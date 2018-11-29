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

	"github.com/TheHipbot/hermes/cache"
	"github.com/TheHipbot/hermes/fs"
	"github.com/TheHipbot/hermes/pkg/prompt"
	"github.com/TheHipbot/hermes/pkg/storage"
	"github.com/TheHipbot/hermes/repo"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"
)

var (
	cfgFile  string
	aliasFlg bool
	configFS *fs.ConfigFS
	fsCache  cache.Cache
	store *storage.Storage
	prompter prompt.Factory
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

func getHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Requires repo as an argument")
		os.Exit(1)
	}
	repoName := args[0]
	pathToRepo := fmt.Sprintf("%s%s/", viper.GetString("repo_path"), repoName)
	// repoURL, err := url.Parse(fmt.Sprintf("https://%s", repoName))
	store.Open()
	defer store.Close()

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	var selectedRepo storage.Repository
	cachedRepos := store.Search(repoName)
	if len(cachedRepos) == 1 {
		selectedRepo = cachedRepos[0]
	} else if len(cachedRepos) == 0 {
		selectedRepo = storage.Repository{
			Name: repoName,
			Path: pathToRepo,
		}
		if err := store.AddRepo(repoName, viper.GetString("repo_path")); err != nil {
			fmt.Printf("Error adding repo to cache %s\n%s\n", pathToRepo, err)
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
	}

	repo := repo.GitRepository{
		Name: selectedRepo.Name,
		URL:  fmt.Sprintf("https://%s", selectedRepo.Name),
	}

	if err := repo.Clone(selectedRepo.Path); err != nil && err != git.ErrRepositoryAlreadyExists {
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
		os.Exit(1)
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
	viper.SetDefault("remotes_file", "remotes.json")

	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(remoteCmd)
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
	prompter = &prompt.Prompter{}

	cacheFile, err := configFS.GetCacheFile()
	if err != nil {
		fmt.Println("Cache file could not be opened or created")
	}
	store = storage.NewStorage(cacheFile)
}
