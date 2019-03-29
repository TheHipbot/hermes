package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/template"
	"github.com/spf13/cobra"
)

var (
	version         = "????"
	commit          = "????"
	timestamp       = "????"
	versionTemplate = `Version:    {{ .Version }}
Git Commit: {{ .GitCommit }}
Built:      {{ .BuildTime }}`

	errParseTemplate   = errors.New("Error parsing version template")
	errExecuteTemplate = errors.New("Error executing version template")
)

// VersionData contains data on the version
// and build information of the hermes binary
type VersionData struct {
	Version   string
	GitCommit string
	BuildTime string
}

// versionCmd will output the version of the hermes binary in use
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "View version of hermes binary currently in use",
	Run:   versionHandler,
}

func versionHandler(cmd *cobra.Command, args []string) {
	data := &VersionData{
		Version:   version,
		GitCommit: commit,
		BuildTime: timestamp,
	}
	result, err := generateVersionOutput(versionTemplate, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func generateVersionOutput(tmpl string, data *VersionData) (string, error) {
	t, err := template.New("version").Parse(tmpl)
	if err != nil {
		return "", errParseTemplate
	}
	var resolved bytes.Buffer
	if err := t.Execute(&resolved, data); err != nil {
		return "", errExecuteTemplate
	}
	return resolved.String(), nil
}
