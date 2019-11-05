package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	hardRmFlg bool
)

func init() {
	repoCmd.AddCommand(repoRmCommand)

	repoRmCommand.Flags().BoolVar(&hardRmFlg, "hard", false, "remove repo from disk")
}

// repoCmd represents the base remote command when called without any subcommands
var repoCmd = &cobra.Command{
	Use:     "repo [subcommand]",
	Aliases: []string{"repository"},
	Short:   "Manage repositories tracked in hermes",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var repoRmCommand = &cobra.Command{
	Use:     "rm [repo name]",
	Aliases: []string{"remove"},
	Short:   "Remove a repo from cache and optionally from disk",
	Run:     repoRmHandler,
}

func repoRmHandler(cmd *cobra.Command, args []string) {
	store.Open()
	defer store.Close()
	repos := store.SearchRepositories(args[0])
	switch len(repos) {
	case 0:
		fmt.Printf("no repo %s found\n", args[0])
		os.Exit(1)
	case 1:
		if err := store.RemoveRepository(repos[0].Name); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("repo %s removed from cache\n", repos[0].Name)
	default:
		fmt.Println("many repos match your entry, please choose one")
		for _, r := range repos {
			fmt.Printf("  %s\n", r.Name)
		}
		os.Exit(1)
	}

	if hardRmFlg && len(repos) == 1 {
		fmt.Println("here")
		if stat, err := appFs.Stat(repos[0].Path); err != nil {
			fmt.Println("error for repo directory stat, the directory cannot be removed")
			os.Exit(1)
		} else if stat.IsDir() {
			removeDirRecursive(repos[0].Path)
			removeEmptyDirs(repos[0].Path[:strings.LastIndex(repos[0].Path[:len(repos[0].Path)-1], "/")+1], viper.GetString("repo_path"))
		}
	}
}

func removeEmptyDirs(path, base string) error {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	if !strings.HasSuffix(base, "/") {
		base = base + "/"
	}

	if path == base {
		return nil
	}

	if items, err := appFs.ReadDir(path); err != nil {
		return err
	} else if len(items) == 0 {
		if err := appFs.Remove(path); err != nil {
			return nil
		}
		return removeEmptyDirs(path[:strings.LastIndex(path[:len(path)-1], "/")+1], base)
	}
	return nil
}

func removeDirRecursive(path string) error {
	if stat, err := appFs.Stat(path); err != nil {
		fmt.Println("error for repo directory stat, the directory cannot be removed")
		os.Exit(1)
	} else if stat.IsDir() {
		items, _ := appFs.ReadDir(path)
		for _, item := range items {
			if !item.IsDir() {
				if err := appFs.Remove(fmt.Sprintf("%s/%s", path, item.Name())); err != nil {
					return err
				}
			} else {
				if err := removeDirRecursive(fmt.Sprintf("%s%s", path, item.Name())); err != nil {
					return err
				}
			}
		}
		return appFs.Remove(path)
	}
	return nil
}
