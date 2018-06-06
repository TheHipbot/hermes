package fs

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

// ConfigFS configuration filesystem struct
type ConfigFS struct {
	FS billy.Filesystem
}

// NewConfigFS returns a ConfigFS object with the default FS
func NewConfigFS() *ConfigFS {
	return &ConfigFS{
		FS: osfs.New(""),
	}
}

// TODO bubble the error up
func (c *ConfigFS) checkForConfigDir() error {
	stat, err := c.FS.Stat(viper.GetString("config_path"))
	if err != nil || !stat.IsDir() {
		return errors.New("Config directory doesn't exist or can't be opened, please run hermes setup")
	}
	return nil
}

// Setup runs the initial hermes setup
func (c *ConfigFS) Setup() error {
	return c.FS.MkdirAll(viper.GetString("config_path"), 0751)
}

// SetTarget creates a target file with the directory
// to move to
func (c *ConfigFS) SetTarget(target string) error {
	targetFilePath := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("target_file"))

	// check for config dir
	if err := c.checkForConfigDir(); err != nil {
		return err
	}

	file, err := c.FS.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(target))
	if err != nil {
		return err
	}
	return nil
}

// GetCacheFile gets the cache file from the config folder, if file doesn't exists
// it attempts to create it
func (c *ConfigFS) GetCacheFile() (billy.File, error) {
	cacheFilePath := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("cache_file"))
	if _, err := c.FS.Stat(cacheFilePath); err != nil {
		return c.FS.Create(cacheFilePath)
	}
	return c.FS.OpenFile(cacheFilePath, os.O_RDWR, 0666)
}
