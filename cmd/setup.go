package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run initial setup for hermes",
	Long:  `Setup will create the hermes config directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := configFS.Setup(); err != nil {
			fmt.Println("Setup failed")
			os.Exit(1)
		}
	},
}
