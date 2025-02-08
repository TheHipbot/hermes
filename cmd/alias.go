package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/alecthomas/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aliasCmd represents the base command when called without any subcommands
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Outputs shell function for hermes alias",
	Run: func(cmd *cobra.Command, args []string) {
		alias, err := generateAlias()
		if err != nil {
			fmt.Printf("Error generating alias\n%s", err)
			os.Exit(1)
		}
		fmt.Print(alias)
		os.Exit(0)
	},
}

// AliasData is a struct containing the data
// for the alias function template
type AliasData struct {
	ConfigDir      string
	TargetFileName string
	AliasName      string
}

func generateAlias() (string, error) {
	var resolved bytes.Buffer
	data := AliasData{
		ConfigDir:      viper.GetString("config_path"),
		TargetFileName: viper.GetString("target_file"),
		AliasName:      viper.GetString("alias_name"),
	}

	aliasTemplate := `function {{ .AliasName }}() {
	local HERMES_BIN=$(whereis hermes | awk -F: '{gsub(/ /, "", $2); print $2}' || which hermes)
	$HERMES_BIN $@
	local EXIT_STATUS=$?
	if [ -f {{ .ConfigDir }}{{ .TargetFileName }} ]; then
		cd $(cat {{ .ConfigDir }}{{ .TargetFileName }}) && rm {{ .ConfigDir }}{{ .TargetFileName }}
	fi
	return $EXIT_STATUS
}`

	t, err := template.New("alias").Parse(aliasTemplate)
	if err != nil {
		return "", err
	}

	if err = t.Execute(&resolved, data); err != nil {
		return "", err
	}

	return resolved.String(), nil
}
