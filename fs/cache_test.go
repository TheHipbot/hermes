package fs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/src-d/go-billy.v4/memfs"
)

var (
	testCache Cache
)

type CacheSuite struct {
	suite.Suite
}

func (s *CacheSuite) SetupTest() {
	viper.Set("config_path", "/test/.hermes/")
	viper.Set("cache_file", "cache.json")

	configFS = &ConfigFS{
		FS: memfs.New(),
	}
	configFS.Setup()
	githubURL, _ := url.Parse("https://github.com")
	gitLabURL, _ := url.Parse("https://gitlab.com")

	testCache = Cache{
		cfs:     configFS,
		Version: cacheFormatVersion,
		Remotes: map[string]*Remote{
			"github.com": &Remote{
				Name: "github.com",
				URL:  githubURL.String(),
				Repos: []Repo{
					Repo{
						Name: "github.com/TheHipbot/hermes",
						Path: "/repos/github.com/TheHipbot/hermes",
					},
					Repo{
						Name: "github.com/TheHipbot/dotfiles",
						Path: "/repos/github.com/TheHipbot/dotfiles",
					},
					Repo{
						Name: "github.com/TheHipbot/dockerfiles",
						Path: "/repos/github.com/TheHipbot/dockerfiles",
					},
					Repo{
						Name: "github.com/src-d/go-git",
						Path: "/repos/github.com/src-d/go-git",
					},
				},
			},
			"gitlab.com": &Remote{
				Name: "gitlab.com",
				URL:  gitLabURL.String(),
				Repos: []Repo{
					Repo{
						Name: "gitlab.com/gitlab-org/gitlab-ce",
						Path: "/repos/gitlab.com/gitlab-org/gitlab-ce",
					},
					Repo{
						Name: "gitlab.com/gnachman/iterm2",
						Path: "/repos/gitlab.com/gnachman/iterm2",
					},
				},
			},
		},
	}
}

func (s *CacheSuite) TestNewClient() {
	cache := NewCache()
	s.Equal(cache.cfs, configFS, "NewClient should return a cache object with the cfs set")
}

func (s *CacheSuite) TestCacheOpenWithFileData() {
	err := configFS.WriteCache([]byte(`{
		"version": "0.0.1",
		"remotes": {
			"github.com": {
				"name": "github.com",
				"url":  "https://github.com",
				"repos": [
					{
						"name": "github.com/TheHipbot/hermes",
						"repo_path": "/repos/github.com/TheHipbot/hermes"
					},
					{
						"name": "github.com/TheHipbot/dotfiles",
						"repo_path": "/repos/github.com/TheHipbot/dotfiles"
					},
					{
						"name": "github.com/TheHipbot/dockerfiles",
						"repo_path": "/repos/github.com/TheHipbot/dockerfiles"
					},
					{
						"name": "github.com/src-d/go-git",
						"repo_path": "/repos/github.com/src-d/go-git"
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
	}`))
	s.Nil(err, "Cache should write successfully")
	cache := &Cache{
		cfs: configFS,
	}
	cache.Open()
	fmt.Println(cache.Remotes)

	s.Equal(cache.Version, "0.0.1", "Cache format version should be 0.0.1")
	s.NotNil(cache.Remotes["github.com"], "There should be repos in the github.com remote")
	s.Equal(len(cache.Remotes["github.com"].Repos), 4, "There should be 4 repos in the github.com remote")
	s.Equal(cache.Remotes["github.com"].Repos[0].Name, "github.com/TheHipbot/hermes", "The first repo in the github.com remote should be hermes")
	s.NotNil(cache.Remotes["gitlab.com"], "There should be repos in the gitlab.com remote")
	s.Equal(len(cache.Remotes["gitlab.com"].Repos), 2, "There should be 4 repos in the gitlab.com remote")
	s.Equal(cache.Remotes["gitlab.com"].Repos[0].Name, "gitlab.com/gitlab-org/gitlab-ce", "The first repo in the gitlab.com remote should be gitlab-ce")
}

func (s *CacheSuite) TestCacheOpenWithReadError() {
	cache := &Cache{
		cfs: configFS,
	}
	cache.Open()

	s.Equal(cache.Version, cacheFormatVersion, "Cache format version should be set")
	s.NotNil(len(cache.Remotes), "There should be no remotes in the Remotes map")
}

func (s *CacheSuite) TestCacheOpenWithInvalidData() {
	configFS.WriteCache([]byte(`{
	"version": "0.0.1",
}`))
	cache := &Cache{
		cfs: configFS,
	}
	cache.Open()

	s.Equal(cache.Version, cacheFormatVersion, "Cache format version should be set")
	s.NotNil(len(cache.Remotes), "There should be no remotes in the Remotes map")
}

func (s *CacheSuite) TestCacheSave() {
	s.Nil(testCache.Save(), "testCache should save successfully")
	cache := &Cache{
		cfs: configFS,
	}
	cache.Open()
	s.Equal(cache.Version, testCache.Version, "Versions between caches should be equal")
	s.Equal(len(cache.Remotes), len(testCache.Remotes), "Caches should have equal number of remotes")
	s.Equal(len(cache.Remotes["github.com"].Repos), len(testCache.Remotes["github.com"].Repos), "Caches should have equal number of github.com repos")
	s.Equal(len(cache.Remotes["gitlab.com"].Repos), len(testCache.Remotes["gitlab.com"].Repos), "Caches should have equal number of gitlab.com repos")
}

func (s *CacheSuite) TestCacheAdd() {
	var results []Repo

	repoCnt := len(testCache.Remotes["github.com"].Repos)
	testCache.Add("github.com/TheHipbot/weather", "/repos/")
	results = testCache.Search("weather")
	s.Len(results, 1, "There should be the new repo")
	s.Equal(repoCnt+1, len(testCache.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")

	testCache.Add("github.com/TheHipbot/docker", "/repos/")
	results = testCache.Search("docker")
	s.Len(results, 2, "There should be the new repo")
	s.Equal(repoCnt+2, len(testCache.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
}

func (s *CacheSuite) TestCacheAddNewRemote() {
	var results []Repo

	remote := testCache.Remotes["gopkg.in"]
	s.Nil(remote, "The remote should not exist")
	testCache.Add("gopkg.in/src-d/go-billy.v4", "/repos/")
	results = testCache.Search("billy")
	s.Len(results, 1, "There should be the new repo")
	remote = testCache.Remotes["gopkg.in"]
	s.NotNil(remote, "The remote should exist")
}

func (s *CacheSuite) TestCacheAddThenSave() {
	var results []Repo

	repoCnt := len(testCache.Remotes["github.com"].Repos)
	testCache.Add("github.com/TheHipbot/weather", "/repos/")
	results = testCache.Search("weather")
	s.Len(results, 1, "There should be the new repo")
	s.Equal(repoCnt+1, len(testCache.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
	testCache.Save()

	var temp Cache
	raw, err := configFS.ReadCache()
	fmt.Println(configFS.ReadCache())
	s.Nil(err, "Should be unmarshallable")
	s.Nil(json.Unmarshal(raw, &temp), "Should be unmarshallable")
	s.Equal(repoCnt+1, len(temp.Remotes["github.com"].Repos), "The new repo should be stored with existing remote in cache")

	testCache.Add("github.com/TheHipbot/docker", "/repos/")
	results = testCache.Search("docker")
	s.Len(results, 2, "There should be the new repo")
	s.Equal(repoCnt+2, len(testCache.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
	testCache.Save()

	raw, err = configFS.ReadCache()
	s.Nil(err, "Should be unmarshallable")
	s.Nil(json.Unmarshal(raw, &temp), "Should be unmarshallable")
	s.Equal(repoCnt+2, len(temp.Remotes["github.com"].Repos), "The new repo should be stored with existing remote in cache")
}

func (s *CacheSuite) TestRemoveCache() {
	repoCnt := len(testCache.Remotes["github.com"].Repos)
	testCache.Remove("github.com/TheHipbot/dotfiles")
	s.Equal(repoCnt-1, len(testCache.Remotes["github.com"].Repos), "github.com remote should have one less repo")
	results := testCache.Search("dotfiles")
	s.Equal(len(results), 0, "There should no longer be a dotfiles repo")
}

func (s *CacheSuite) TestRemoveCacheNoRepo() {
	repoCnt := len(testCache.Remotes["github.com"].Repos)
	err := testCache.Remove("github.com/TheHipbot/docker")
	s.NotNil(err, "There should be an error returned")
	s.Equal(repoCnt, len(testCache.Remotes["github.com"].Repos), "github.com remote should have the same number of repos")
}

func (s *CacheSuite) TestCacheSearchWithResults() {
	var results []Repo

	results = testCache.Search("files")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/dotfiles")
	s.Equal(results[1].Name, "github.com/TheHipbot/dockerfiles")

	results = testCache.Search("gitlab")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "gitlab.com/gitlab-org/gitlab-ce")
	s.Equal(results[1].Name, "gitlab.com/gnachman/iterm2")
}

func (s *CacheSuite) TestCacheSearchCaseInSensitiveWithResults() {
	var results []Repo

	results = testCache.Search("FILES")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/dotfiles")
	s.Equal(results[1].Name, "github.com/TheHipbot/dockerfiles")

	results = testCache.Search("thehipbot")
	s.Len(results, 3, "There are 3 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/hermes")
	s.Equal(results[1].Name, "github.com/TheHipbot/dotfiles")
	s.Equal(results[2].Name, "github.com/TheHipbot/dockerfiles")
}

func (s *CacheSuite) TestCacheSearchWithoutResults() {
	var results []Repo

	results = testCache.Search("test")
	s.Len(results, 0, "There no results")

	testCache.Remotes = map[string]*Remote{}
	results = testCache.Search("files")
	s.Len(results, 0, "There no results")
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
