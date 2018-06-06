package fs

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/spf13/viper"

	"gopkg.in/src-d/go-billy.v4/memfs"
)

var (
	cfs            ConfigFS
	testConfigPath string
	testTargetFile string
	testCacheFile  string
)

type ConfigFSSuite struct {
	suite.Suite
}

func (s *ConfigFSSuite) SetupTest() {
	testConfigPath = "/test/.hermes/"
	testTargetFile = ".hermes_test"
	testCacheFile = "cache.json"
	viper.Set("config_path", testConfigPath)
	viper.Set("target_file", testTargetFile)
	viper.Set("cache_file", testCacheFile)
	cfs = ConfigFS{
		FS: memfs.New(),
	}
}

func (s *ConfigFSSuite) TestSetupCreateDir() {
	cfs.Setup()

	_, err := cfs.FS.Stat(testConfigPath)
	s.Nil(err, "Setup should create config dir")
}

func (s *ConfigFSSuite) TestSetTarget() {
	target := "/repo_dir/github.com/TheHipbot/hermes/"

	// set up to create config_dir in memfs
	cfs.Setup()

	cfs.SetTarget(target)
	stat, err := cfs.FS.Stat(fmt.Sprintf("%s%s", testConfigPath, testTargetFile))
	s.Nil(err, "SetTarget should stat a target file")
	bs := make([]byte, stat.Size())
	file, err := cfs.FS.Open(fmt.Sprintf("%s%s", testConfigPath, testTargetFile))
	s.Nil(err, "SetTarget should create a target file")

	_, err = file.Read(bs)
	s.Nil(err, "Target file should be read")
	s.True(strings.Contains(string(bs), target), "Target file content incorrect")
}

func (s *ConfigFSSuite) GetCacheFileCreateTest() {
	file, err := cfs.GetCacheFile()
	s.NotNil(err, "GetCacheFile should not error")
	_, err = file.Write([]byte("test"))
	s.Nil(err, "File should be writable")
	file.Close()
	info, err := cfs.FS.Stat(fmt.Sprintf("%s%s", testConfigPath, testCacheFile))
	s.NotNil(err, "Should be able to stat file")
	s.True(info.Mode().IsRegular(), "Cache file should exit")
}

func (s *ConfigFSSuite) GetCacheFileExistingTest() {
	file, err := cfs.GetCacheFile()
	s.NotNil(err, "GetCacheFile should not error")
	file.Write([]byte("test"))
	file.Close()
	file, err = cfs.FS.Open(fmt.Sprintf("%s%s", testConfigPath, testCacheFile))
	s.NotNil(err, "Should be able to stat file")
	content, err := ioutil.ReadAll(file)
	s.NotNil(err, "Should be able to read file")
	s.Equal([]byte("test"), content, "Cache file should exit")
}

func TestConfigFSSuite(t *testing.T) {
	suite.Run(t, new(ConfigFSSuite))
}
