package fs

import (
	"errors"
	"fmt"

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

// ReadCache writes the given byte area out to cache.json
func (c *ConfigFS) ReadCache() ([]byte, error) {
	var data []byte
	var file billy.File
	cacheFilePath := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("cache_file"))

	// check for config dir
	if err := c.checkForConfigDir(); err != nil {
		return nil, err
	}

	stat, err := c.FS.Stat(cacheFilePath)
	if err != nil {
		return nil, err
	}
	file, err = c.FS.Open(cacheFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data = make([]byte, stat.Size())

	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// WriteCache writes the given byte area out to cache.json
func (c *ConfigFS) WriteCache(data []byte) error {
	cacheFilePath := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("cache_file"))

	// check for config dir
	if err := c.checkForConfigDir(); err != nil {
		return err
	}

	c.FS.Remove(cacheFilePath)
	file, err := c.FS.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
