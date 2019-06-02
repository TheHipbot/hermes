// Copyright © 2018 Jeremy Chambers <jeromext@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/TheHipbot/hermes/pkg/credentials"

	fscred "github.com/TheHipbot/hermes/pkg/credentials/osfs"
	"github.com/TheHipbot/hermes/pkg/fs"
	"github.com/TheHipbot/hermes/pkg/prompt"
	"github.com/TheHipbot/hermes/pkg/repo"
	"github.com/TheHipbot/hermes/pkg/storage"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

var (
	cfgFile   string
	aliasFlg  bool
	appFs     billy.Filesystem
	configFS  *fs.ConfigFS
	store     storage.Storage
	prompter  prompt.Factory
	protocols = []string{
		"https",
		"ssh",
		"http",
	}
	credentialsStorer credentials.Storer
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hermes",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	Run: getHandler,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a new repo or go to an existing",
	Run:   getHandler,
}

// TODO: add repository only on successful clone
// TODO: update ssh to not assume username
func getHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Requires repo as an argument")
		os.Exit(ExitInvalidArguments)
	}
	repoName := args[0]
	pathToRepo := fmt.Sprintf("%s%s/", viper.GetString("repo_path"), repoName)
	store.Open()
	defer store.Close()

	var selectedRepo storage.Repository
	var remote *storage.Remote
	cachedRepos := store.SearchRepositories(repoName)
	if len(cachedRepos) == 1 {
		selectedRepo = cachedRepos[0]
		remote, _ = store.SearchRemote(strings.Split(selectedRepo.Name, "/")[0])
	} else if len(cachedRepos) == 0 {
		parts := strings.Split(repoName, "/")
		if len(parts) < 3 {
			fmt.Printf(`No repo found, a new repo must be in the form
<remote hostname>/<user or group>/<repo name>
`)
			os.Exit(ExitInvalidArguments)
		}
		remoteName := parts[0]
		remote, ok := store.SearchRemote(remoteName)
		selectedRepo = storage.Repository{
			Name: repoName,
			Path: pathToRepo,
		}
		if err := store.AddRepository(&storage.Repository{
			Name: repoName,
			Path: viper.GetString("repo_path"),
		}); err != nil {
			fmt.Printf("Error adding repo to cache %s\n%s\n", pathToRepo, err)
		}
		if !ok {
			// prompt user for protocol
			p := prompt.CreateProtoclSelectPrompt(prompter, protocols)
			i, _, err := p.Run()
			if err != nil {
				fmt.Printf("Error retrieving input\n")
				os.Exit(1)
			}
			remote, _ = store.SearchRemote(remoteName)
			remote.Protocol = protocols[i]
		}
		store.Save()
	} else {
		p := prompt.CreateRepoSelectPrompt(prompter, cachedRepos)
		i, _, err := p.Run()
		selectedRepo = cachedRepos[i]
		if err != nil {
			fmt.Printf("Error selecting repo\n%s\n", err)
			os.Exit(1)
		}
		remote, _ = store.SearchRemote(strings.Split(selectedRepo.Name, "/")[0])
	}

	targetRepo := repo.NewGitRepository(selectedRepo.Name, "")
	targetRepo.Fs = appFs
	cloner, _ := repo.NewCloner("git")
	targetRepo.Cloner = cloner

	switch remote.Protocol {
	case "ssh":
		if selectedRepo.SSHURL != "" {
			targetRepo.URL = selectedRepo.SSHURL
			targetRepo.Protocol = "ssh"
		} else {
			targetRepo.URL = fmt.Sprintf("ssh://git@%s", selectedRepo.Name)
			targetRepo.Protocol = "ssh"
		}
	case "http":
		if selectedRepo.CloneURL != "" {
			targetRepo.URL = selectedRepo.CloneURL
		} else {
			targetRepo.URL = fmt.Sprintf("http://%s", selectedRepo.Name)
		}
	default:
		if selectedRepo.CloneURL != "" {
			targetRepo.URL = selectedRepo.CloneURL
		} else {
			targetRepo.URL = fmt.Sprintf("https://%s", selectedRepo.Name)
		}
	}

	if err := targetRepo.Clone(selectedRepo.Path); err != nil && err != repo.ErrRepoAlreadyExists {
		fmt.Printf("Error cloning repo %s\n%s\n", selectedRepo.Path, err)
		os.Exit(1)
	}

	if err := configFS.SetTarget(selectedRepo.Path); err != nil {
		fmt.Printf("Error creating target file\n%s\n", err)
		os.Exit(1)
	}
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(ExitCannotExecute)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hermes.yaml)")

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetDefault("repo_path", fmt.Sprintf("%s/hermes-repos/", home))
	viper.SetDefault("config_path", fmt.Sprintf("%s/.hermes/", home))
	viper.SetDefault("target_file", ".hermes_target")
	viper.SetDefault("cache_file", "cache.json")
	viper.SetDefault("alias_name", "hermes")
	viper.SetDefault("remotes_file", "remotes.json")
	viper.SetDefault("credentials_type", "none")
	viper.SetDefault("credentials_file", "credentials.yml")

	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(remoteCmd)
	rootCmd.AddCommand(versionCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hermes" (without extension).
		viper.SetConfigName(".hermes")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/hermes/")
	}

	viper.ReadInConfig()
	viper.AutomaticEnv() // read in environment variables that match
	configFS = fs.NewConfigFS()
	appFs = osfs.New("")

	prompter = &prompt.Prompter{}

	cacheFile, err := configFS.GetCacheFile()
	if err != nil {
		fmt.Println("Cache file could not be opened or created")
	}
	store = storage.NewStorage(cacheFile)

	switch viper.GetString("credentials_type") {
	case "file":
		credentialFileName := fmt.Sprintf("%s/%s", viper.GetString("config_path"), viper.GetString("credentials_file"))
		if file, err := appFs.OpenFile(credentialFileName, os.O_RDWR, 0666); err == nil {
			credentialsStorer = fscred.NewFSStorer(file)
		} else if file, err := appFs.Create(credentialFileName); err == nil {
			credentialsStorer = fscred.NewFSStorer(file)
		} else {
			fmt.Println("Credentials file could not be opened or created, no credentials will be persisted")
		}
	}

	if credentialsStorer == nil {
		credentialsStorer = credentials.NewMemStorer()
	}
}
