package fs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/spf13/viper"

	"gopkg.in/src-d/go-billy.v4/memfs"
)

var (
	configFS       ConfigFS
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
	configFS = ConfigFS{
		FS: memfs.New(),
	}
}

func (s *ConfigFSSuite) TestSetupCreateDir() {
	configFS.Setup()

	_, err := configFS.FS.Stat(testConfigPath)
	s.Nil(err, "Setup should create config dir")
}

func (s *ConfigFSSuite) TestSetTarget() {
	target := "/repo_dir/github.com/TheHipbot/hermes/"

	// set up to create config_dir in memfs
	configFS.Setup()

	configFS.SetTarget(target)
	stat, err := configFS.FS.Stat(fmt.Sprintf("%s%s", testConfigPath, testTargetFile))
	s.Nil(err, "SetTarget should stat a target file")
	bs := make([]byte, stat.Size())
	file, err := configFS.FS.Open(fmt.Sprintf("%s%s", testConfigPath, testTargetFile))
	s.Nil(err, "SetTarget should create a target file")

	_, err = file.Read(bs)
	s.Nil(err, "Target file should be read")
	s.True(strings.Contains(string(bs), target), "Target file content incorrect")
}

func (s *ConfigFSSuite) TestReadCache() {
	cachePath := fmt.Sprintf("%s%s", testConfigPath, testCacheFile)

	file, err := configFS.FS.Create(cachePath)
	s.Nil(err, "Cache file should be created")

	testCache := []byte(`{
	"version": "0.0.1",
	"remotes": {
		"github.com": {
			"name": "github.com",
			"url":  "https://github.com",
			"repos": [
				{
					"name": "github.com/TheHipbot/hermes",
					"repo_path": "/repos/github.com/TheHipbot/hermes",
				},
				{
					"name": "github.com/TheHipbot/dotfiles",
					"repo_path": "/repos/github.com/TheHipbot/dotfiles",
				},
				{
					"name": "github.com/TheHipbot/dockerfiles",
					"repo_path": "/repos/github.com/TheHipbot/dockerfiles",
				},
				{
					"name": "github.com/src-d/go-git",
					"repo_path": "/repos/github.com/src-d/go-git",
				}
			]
		},
		"gitlab.com": {
			"name": "gitlab.com",
			"url":  "https://gitlab.com",
			"repos": [
				{
					"name": "gitlab.com/gitlab-org/gitlab-ce",
					"repo_path": "/repos/gitlab.com/gitlab-org/gitlab-ce"
				},
				{
					"name": "gitlab.com/gnachman/iterm2",
					"repo_path": "/repos/gitlab.com/gnachman/iterm2"
				}
			]
		}
	}
}`)
	file.Write(testCache)

	c, err := configFS.ReadCache()
	s.Nil(err, "Cache file should be read without error")
	s.Equal(string(testCache), string(c), "Cache should be read from cache file in config_path")
}

func (s *ConfigFSSuite) TestReadCacheNoFile() {
	configFS.FS.MkdirAll(viper.GetString("config_path"), 0751)
	_, err := configFS.ReadCache()
	s.NotNil(err, "ReadCache should return an error if no file present")
}

func (s *ConfigFSSuite) TestWriteCache() {
	cachePath := fmt.Sprintf("%s%s", testConfigPath, testCacheFile)
	configFS.FS.MkdirAll(viper.GetString("config_path"), 0751)

	testCache := []byte(`{
	"version": "0.0.1",
	"remotes": {
		"github.com": {
			"name": "github.com",
			"url":  "https://github.com",
			"repos": [
				{
					"name": "github.com/TheHipbot/hermes",
					"repo_path": "/repos/github.com/TheHipbot/hermes",
				},
				{
					"name": "github.com/TheHipbot/dotfiles",
					"repo_path": "/repos/github.com/TheHipbot/dotfiles",
				},
				{
					"name": "github.com/TheHipbot/dockerfiles",
					"repo_path": "/repos/github.com/TheHipbot/dockerfiles",
				},
				{
					"name": "github.com/src-d/go-git",
					"repo_path": "/repos/github.com/src-d/go-git",
				}
			]
		},
		"gitlab.com": {
			"name": "gitlab.com",
			"url":  "https://gitlab.com",
			"repos": [
				{
					"name": "gitlab.com/gitlab-org/gitlab-ce",
					"repo_path": "/repos/gitlab.com/gitlab-org/gitlab-ce"
				},
				{
					"name": "gitlab.com/gnachman/iterm2",
					"repo_path": "/repos/gitlab.com/gnachman/iterm2"
				}
			]
		}
	}
}`)
	err := configFS.WriteCache(testCache)
	s.Nil(err, "WriteCache should run without error")
	stat, err := configFS.FS.Stat(cachePath)
	s.Nil(err, "Cache file should get stat")
	file, err := configFS.FS.Open(cachePath)
	s.Nil(err, "Cache file should exist and be opened")

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	s.Nil(err, "Should read data from cache file")
	s.Equal(string(testCache), string(data), "Cache should be read from cache file in config_path")
}

func (s *ConfigFSSuite) TestWriteCacheOverwrite() {
	cachePath := fmt.Sprintf("%s%s", testConfigPath, testCacheFile)
	configFS.FS.MkdirAll(viper.GetString("config_path"), 0751)

	testCache := []byte(`{
	"version": "0.0.1",
	"remotes": {
		"github.com": {
			"name": "github.com",
			"url":  "https://github.com",
			"repos": [
				{
					"name": "github.com/TheHipbot/hermes",
					"repo_path": "/repos/github.com/TheHipbot/hermes",
				},
				{
					"name": "github.com/TheHipbot/dotfiles",
					"repo_path": "/repos/github.com/TheHipbot/dotfiles",
				},
				{
					"name": "github.com/TheHipbot/dockerfiles",
					"repo_path": "/repos/github.com/TheHipbot/dockerfiles",
				},
				{
					"name": "github.com/src-d/go-git",
					"repo_path": "/repos/github.com/src-d/go-git",
				}
			]
		},
		"gitlab.com": {
			"name": "gitlab.com",
			"url":  "https://gitlab.com",
			"repos": [
				{
					"name": "gitlab.com/gitlab-org/gitlab-ce",
					"repo_path": "/repos/gitlab.com/gitlab-org/gitlab-ce"
				},
				{
					"name": "gitlab.com/gnachman/iterm2",
					"repo_path": "/repos/gitlab.com/gnachman/iterm2"
				}
			]
		}
	}
}`)
	err := configFS.WriteCache(testCache)
	s.Nil(err, "WriteCache should run without error")

	testCacheOverride := []byte(`{
		"version": "0.0.1",
		"remotes": {
			"github.com": {
				"name": "github.com",
				"url":  "https://github.com",
				"repos": [
					{
						"name": "github.com/TheHipbot/dotfiles",
						"repo_path": "/repos/github.com/TheHipbot/dotfiles",
					},
					{
						"name": "github.com/TheHipbot/dockerfiles",
						"repo_path": "/repos/github.com/TheHipbot/dockerfiles",
					},
					{
						"name": "github.com/src-d/go-git",
						"repo_path": "/repos/github.com/src-d/go-git",
					}
				]
			},
			"gitlab.com": {
				"name": "gitlab.com",
				"url":  "https://gitlab.com",
				"repos": [
					{
						"name": "gitlab.com/gnachman/iterm2",
						"repo_path": "/repos/gitlab.com/gnachman/iterm2"
					}
				]
			}
		}
	}`)

	err = configFS.WriteCache(testCacheOverride)
	s.Nil(err, "WriteCache should run without error")
	stat, err := configFS.FS.Stat(cachePath)
	s.Nil(err, "Cache file should get stat")
	file, err := configFS.FS.Open(cachePath)
	s.Nil(err, "Cache file should exist and be opened")

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	s.Nil(err, "Should read data from cache file")
	s.Equal(string(testCacheOverride), string(data), "Cache should be read from cache file in config_path")
}

func TestConfigFSSuite(t *testing.T) {
	suite.Run(t, new(ConfigFSSuite))
}
