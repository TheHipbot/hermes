package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/TheHipbot/hermes/pkg/credentials"

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

	// flag vars
	getAllReposFlg bool
	protocolFlg    = ""
	remoteTypeFlg  = ""
	tokenFlg       = ""

	errInput           = errors.New("error retrieving input")
	errInvalidProtocol = fmt.Errorf("invalid protocol, valid values are %s ", protocols)
	errInvalidRemote   = errors.New("invlaid remote url")
	errRetrievingRepos = errors.New("error retrieving repos")
)

func init() {
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteRefreshCmd)
	remoteCmd.Flags().BoolVarP(&getAllReposFlg, "all", "a", false, "get all repos")
	remoteCmd.Flags().StringVarP(&protocolFlg, "protocol", "p", "", "protocol to use for repos of given remote(s)")
	remoteAddCmd.Flags().StringVarP(&remoteTypeFlg, "type", "t", "", "remote type (e.g. github, gitlab, etc.)")
	remoteAddCmd.Flags().StringVar(&tokenFlg, "token", "", "auth token")
}

// remoteCmd represents the base remote command when called without any subcommands
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remotes for hermes repositories",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
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

func promptAndGetAuth(remoteURL *url.URL, driver remote.Driver) (remote.Auth, error) {
	remoteName := remoteURL.Hostname()
	auth := remote.Auth{}
	switch driver.AuthType() {
	case "token":
		if tokenFlg != "" {
			credentialsStorer.Put(remoteName, credentials.Credential{
				Type:  "token",
				Token: tokenFlg,
			})
			auth.Token = tokenFlg
		} else if cred, err := credentialsStorer.Get(remoteName); err != nil {
			fmt.Println(err)
			ip := prompt.CreateTokenInputPrompt(prompter)
			// TODO: handle error here
			// TODO: reprompt on failed auth
			token, err := ip.Run()
			if err != nil {
				return remote.Auth{}, err
			}
			credentialsStorer.Put(remoteName, credentials.Credential{
				Type:  "token",
				Token: token,
			})
			auth.Token = token
		} else {
			auth.Token = cred.Token
		}
	default:
	}
	return auth, nil
}

func getProtocolIndex() (int, error) {
	protocolIndex := -1
	var err error
	if protocolFlg != "" {
		for i, p := range protocols {
			if p == protocolFlg {
				protocolIndex = i
				break
			}
		}
		if protocolIndex < 0 {
			return protocolIndex, errInvalidProtocol
		}
	}

	if protocolIndex < 0 {
		p := prompt.CreateProtocolSelectPrompt(prompter, protocols)
		protocolIndex, _, err = p.Run()
		if err != nil {
			return protocolIndex, errInput
		}
	}

	return protocolIndex, nil
}

func remoteAddHandler(cmd *cobra.Command, args []string) {
	store.Open()
	defer store.Close()
	defer store.Save()
	defer credentialsStorer.Close()
	if err := addReposFromRemote(args[0]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addReposFromRemote(remoteStr string) error {
	remoteURL, err := url.Parse(remoteStr)
	remoteName := remoteURL.Hostname()
	if err != nil {
		return errInvalidRemote
	}
	cachedRemote, remoteCached := store.SearchRemote(remoteName)

	var remoteType string
	if remoteTypeFlg != "" {
		remoteType = remoteTypeFlg
	} else {
		if remoteCached {
			remoteType = cachedRemote.Type
		} else {
			p := prompt.CreateDriverSelectPrompt(prompter, drivers)
			i, _, err := p.Run()
			if err != nil {
				return errInput
			}
			remoteType = drivers[i].Name
		}
	}
	driver, err := remote.NewDriver(remoteType, &remote.DriverOpts{
		AllRepos: getAllReposFlg,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	driver.SetHost(remoteURL.String())

	auth, err := promptAndGetAuth(remoteURL, driver)
	if err != nil {
		return errInput
	}
	driver.Authenticate(auth)
	repos, err := driver.GetRepos()
	for err == remote.ErrAuth {
		fmt.Println("Authentication error received from remote")
		credentialsStorer.Delete(remoteName)
		auth, err = promptAndGetAuth(remoteURL, driver)
		if err == nil {
			driver.Authenticate(auth)
			repos, err = driver.GetRepos()
		}
	}

	if err != nil {
		return errRetrievingRepos
	}

	if !remoteCached {
		protocolIndex, err := getProtocolIndex()
		if err != nil {
			return err
		}

		store.AddRemote(remoteURL.String(), remoteName, remoteType, protocols[protocolIndex])
	} else if protocolFlg != "" {
		protocolIndex, err := getProtocolIndex()
		if err != nil {
			return err
		}
		cachedRemote.Protocol = protocols[protocolIndex]
	}

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
	return nil
}

// remoteRefreshCmd represents the base command when called without any subcommands
var remoteRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh remotes and authentication for hermes repositories",
	Run:   remoteRefreshHandler,
}

func remoteRefreshHandler(cmd *cobra.Command, args []string) {
	store.Open()
	defer store.Close()
	defer store.Save()
	defer credentialsStorer.Close()

	var aggErr error
	for _, r := range store.ListRemotes() {
		fmt.Printf("refreshing %s\n", r.Name)
		if err := addReposFromRemote(r.URL); err != nil {
			fmt.Println(err)
			aggErr = err
		}
	}
	if aggErr != nil {
		os.Exit(1)
	}
}
