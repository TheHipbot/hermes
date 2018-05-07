package fs

import (
	"errors"
	"fmt"

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

// TODO bubble the error up
func checkForConfigDir() error {
	stat, err := configFs.Lstat(viper.GetString("config_path"))
	if err != nil || !stat.IsDir() {
		return errors.New("Config directory doesn't exist or can't be opened, please run hermes setup")
	}
	return nil
}

// Setup runs the initial hermes setup
func Setup() error {
	return configFs.MkdirAll(viper.GetString("config_path"), 0751)
}

// SetTarget creates a target file with the directory
// to move to
func SetTarget(target string) error {
	targetFilePath := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("target_file"))

	// check for config dir
	if err := checkForConfigDir(); err != nil {
		return err
	}

	file, err := configFs.Create(targetFilePath)
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
func ReadCache() ([]byte, error) {
	var data []byte
	var file billy.File
	cacheFilePath := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("cache_file"))

	// check for config dir
	if err := checkForConfigDir(); err != nil {
		return nil, err
	}

	stat, err := configFs.Stat(cacheFilePath)
	if err != nil {
		return nil, err
	}
	file, err = configFs.Open(cacheFilePath)
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
func WriteCache(data []byte) error {
	cacheFilePath := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("cache_file"))

	// check for config dir
	if err := checkForConfigDir(); err != nil {
		return err
	}

	configFs.Remove(cacheFilePath)
	file, err := configFs.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
