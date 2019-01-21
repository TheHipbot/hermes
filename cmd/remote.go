package cmd

import (
	"fmt"
	"os"

	"github.com/TheHipbot/hermes/pkg/prompt"
	"github.com/TheHipbot/hermes/pkg/remote"
	"github.com/TheHipbot/hermes/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type driver struct {
	Name string
}

var (
	drivers = []driver{
		driver{
			Name: "github",
		},
		driver{
			Name: "gitlab",
		},
	}
	getAllRepos bool
)

func init() {
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.Flags().BoolVarP(&getAllRepos, "all", "a", false, "get all repos")
}

// remoteCmd represents the base remote command when called without any subcommands
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remotes for hermes repositories",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

// remoteAddCmd represents the base command when called without any subcommands
var remoteAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add remotes and authentication for hermes repositories",
	Args:  cobra.MinimumNArgs(1),
	Run:   remoteAddHandler,
}

func remoteAddHandler(cmd *cobra.Command, args []string) {
	remoteName := args[0]
	p := prompt.CreateDriverSelectPrompt(prompter, drivers)
	i, _, err := p.Run()
	if err != nil {
		fmt.Printf("Error retrieving input")
		os.Exit(1)
	}

	driver, _ := remote.NewDriver(drivers[i].Name, &remote.DriverOpts{
		AllRepos: getAllRepos,
	})
	driver.SetHost(remoteName)
	auth := remote.Auth{}
	switch driver.AuthType() {
	case "token":
		ip := prompt.CreateTokenInputPrompt(prompter)
		// TODO: handle error here
		token, _ := ip.Run()
		auth.Token = token
	default:
	}

	driver.Authenticate(auth)
	repos, err := driver.GetRepos()
	if err != nil {
		fmt.Println("Error retrieving repos")
		os.Exit(1)
	}

	store.Open()
	defer store.Close()
	defer store.Save()

	p = prompt.CreateProtoclSelectPrompt(prompter, protocols)
	i, _, err = p.Run()
	if err != nil {
		fmt.Printf("Error retrieving input")
		os.Exit(1)
	}

	// TODO check if remote already present
	store.AddRemote(fmt.Sprintf("https://%s", remoteName), remoteName, protocols[i])

	// add repos to cache
	for _, r := range repos {
		repoToAdd := &storage.Repository{
			Name:     r["name"],
			Path:     fmt.Sprintf("%s%s", viper.GetString("repo_path"), r["name"]),
			CloneURL: r["clone_url"],
			SSHURL:   r["ssh_url"],
		}
		store.AddRepository(repoToAdd)
	}
}
