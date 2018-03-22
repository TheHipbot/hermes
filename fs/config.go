package fs

import (
	"github.com/spf13/viper"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

var (
	configFs billy.Filesystem
)

func init() {
	configFs = osfs.New(viper.GetString("config_path"))
}

// Setup runs the initial hermes setup
func Setup() error {
	return configFs.MkdirAll(viper.GetString("config_path"), 0751)
}
