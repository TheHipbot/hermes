package fs

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

var (
	configFs billy.Filesystem
)

const (
	targetFileName = ".hermes_target"
)

func init() {
	configFs = osfs.New(viper.GetString("config_path"))
}

// TODO bubble the error up
func checkForConfigDir() {
	stat, err := configFs.Lstat(viper.GetString("config_path"))
	if err != nil || !stat.IsDir() {
		fmt.Printf("Config directory doesn't exist or can't be opened, please run hermes setup.")
		os.Exit(1)
	}
}

// Setup runs the initial hermes setup
func Setup() error {
	return configFs.MkdirAll(viper.GetString("config_path"), 0751)
}

// SetTarget creates a target file with the directory
// to move to
func SetTarget(target string) error {
	targetFilePath := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("target_name"))

	// check for config dir
	checkForConfigDir()

	file, err := configFs.Create(targetFilePath)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(target))
	if err != nil {
		return err
	}
	return nil
}
