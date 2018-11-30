package storage

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"testing"

	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/stretchr/testify/suite"
)

var (
	testStorage storage
	testStorer  storer
)

type StorageSuite struct {
	suite.Suite
}

func (s *StorageSuite) SetupTest() {
	var err error
	testFs := memfs.New()
	testStorer, err = testFs.Create("cache.json")
	s.Nil(err, "Setup should be able to create test cache file")

	githubURL, _ := url.Parse("https://github.com")
	gitLabURL, _ := url.Parse("https://gitlab.com")

	testStorage = storage{
		storer:  testStorer,
		Version: cacheFormatVersion,
		Remotes: map[string]*Remote{
			"github.com": &Remote{
				Name: "github.com",
				URL:  githubURL.String(),
				Repos: []Repository{
					Repository{
						Name: "github.com/TheHipbot/hermes",
						Path: "/repos/github.com/TheHipbot/hermes",
					},
					Repository{
						Name: "github.com/TheHipbot/dotfiles",
						Path: "/repos/github.com/TheHipbot/dotfiles",
					},
					Repository{
						Name: "github.com/TheHipbot/dockerfiles",
						Path: "/repos/github.com/TheHipbot/dockerfiles",
					},
					Repository{
						Name: "github.com/src-d/go-git",
						Path: "/repos/github.com/src-d/go-git",
					},
				},
			},
			"gitlab.com": &Remote{
				Name: "gitlab.com",
				URL:  gitLabURL.String(),
				Repos: []Repository{
					Repository{
						Name: "gitlab.com/gitlab-org/gitlab-ce",
						Path: "/repos/gitlab.com/gitlab-org/gitlab-ce",
					},
					Repository{
						Name: "gitlab.com/gnachman/iterm2",
						Path: "/repos/gitlab.com/gnachman/iterm2",
					},
				},
			},
		},
	}
}

func (s *StorageSuite) TestNewClient() {
	store := NewStorage(testStorer)
	impl, ok := store.(*storage)

	s.True(ok)
	s.Equal(impl.storer, testStorer, "NewClient should return a cache object with the storer set")
}

func (s *StorageSuite) TestStorageOpenWithFileData() {
	_, err := testStorer.Write([]byte(`{
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
	s.Nil(err, "Storage should write successfully")
	storage := &storage{
		storer: testStorer,
	}
	storage.Open()

	s.Equal(storage.Version, "0.0.1", "storage format version should be 0.0.1")
	s.NotNil(storage.Remotes["github.com"], "There should be repos in the github.com remote")
	s.Equal(len(storage.Remotes["github.com"].Repos), 4, "There should be 4 repos in the github.com remote")
	s.Equal(storage.Remotes["github.com"].Repos[0].Name, "github.com/TheHipbot/hermes", "The first repo in the github.com remote should be hermes")
	s.NotNil(storage.Remotes["gitlab.com"], "There should be repos in the gitlab.com remote")
	s.Equal(len(storage.Remotes["gitlab.com"].Repos), 2, "There should be 4 repos in the gitlab.com remote")
	s.Equal(storage.Remotes["gitlab.com"].Repos[0].Name, "gitlab.com/gitlab-org/gitlab-ce", "The first repo in the gitlab.com remote should be gitlab-ce")
}

func (s *StorageSuite) TestStorageOpenWithReadError() {
	cache := &storage{
		storer: testStorer,
	}
	cache.Open()

	s.Equal(cache.Version, cacheFormatVersion, "Cache format version should be set")
	s.NotNil(len(cache.Remotes), "There should be no remotes in the Remotes map")
}

func (s *StorageSuite) TestStorageOpenWithInvalidData() {
	testStorer.Write([]byte(`{
	"version": "0.0.1",
}`))
	cache := &storage{
		storer: testStorer,
	}
	cache.Open()

	s.Equal(cache.Version, cacheFormatVersion, "Cache format version should be set")
	s.NotNil(len(cache.Remotes), "There should be no remotes in the Remotes map")
}

func (s *StorageSuite) TestCacheSave() {
	s.Nil(testStorage.Save(), "testCache should save successfully")
	cache := &storage{
		storer: testStorer,
	}
	cache.Open()
	s.Equal(cache.Version, testStorage.Version, "Versions between caches should be equal")
	s.Equal(len(cache.Remotes), len(testStorage.Remotes), "Caches should have equal number of remotes")
	s.Equal(len(cache.Remotes["github.com"].Repos), len(testStorage.Remotes["github.com"].Repos), "Caches should have equal number of github.com repos")
	s.Equal(len(cache.Remotes["gitlab.com"].Repos), len(testStorage.Remotes["gitlab.com"].Repos), "Caches should have equal number of gitlab.com repos")
}

func (s *StorageSuite) TestCacheAdd() {
	var results []Repository

	repoCnt := len(testStorage.Remotes["github.com"].Repos)
	testStorage.AddRepository("github.com/TheHipbot/weather", "/repos/")
	results = testStorage.Search("weather")
	s.Len(results, 1, "There should be the new repo")
	s.Equal(repoCnt+1, len(testStorage.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")

	testStorage.AddRepository("github.com/TheHipbot/docker", "/repos/")
	results = testStorage.Search("docker")
	s.Len(results, 2, "There should be the new repo")
	s.Equal(repoCnt+2, len(testStorage.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
}

func (s *StorageSuite) TestCacheAddNewRemote() {
	var results []Repository

	remote := testStorage.Remotes["gopkg.in"]
	s.Nil(remote, "The remote should not exist")
	testStorage.AddRepository("gopkg.in/src-d/go-billy.v4", "/repos/")
	results = testStorage.Search("billy")
	s.Len(results, 1, "There should be the new repo")
	remote = testStorage.Remotes["gopkg.in"]
	s.NotNil(remote, "The remote should exist")
}

func (s *StorageSuite) TestCacheAddThenSave() {
	var results []Repository

	repoCnt := len(testStorage.Remotes["github.com"].Repos)
	testStorage.AddRepository("github.com/TheHipbot/weather", "/repos/")
	results = testStorage.Search("weather")
	s.Len(results, 1, "There should be the new repo")
	s.Equal(repoCnt+1, len(testStorage.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
	err := testStorage.Save()
	s.Nil(err, "Should save")
	var temp storage
	raw, err := ioutil.ReadAll(testStorer)
	s.Nil(err, "Should be unmarshallable")
	s.Nil(json.Unmarshal(raw, &temp), "Should be unmarshallable")
	s.Equal(repoCnt+1, len(temp.Remotes["github.com"].Repos), "The new repo should be stored with existing remote in cache")

	testStorage.AddRepository("github.com/TheHipbot/docker", "/repos/")
	results = testStorage.Search("docker")
	s.Len(results, 2, "There should be the new repo")
	s.Equal(repoCnt+2, len(testStorage.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
	testStorage.Save()

	raw, err = ioutil.ReadAll(testStorer)
	s.Nil(err, "Should be unmarshallable")
	s.Nil(json.Unmarshal(raw, &temp), "Should be unmarshallable")
	s.Equal(repoCnt+2, len(temp.Remotes["github.com"].Repos), "The new repo should be stored with existing remote in cache")
}

func (s *StorageSuite) TestRemoveRepository() {
	repoCnt := len(testStorage.Remotes["github.com"].Repos)
	testStorage.RemoveRepository("github.com/TheHipbot/dotfiles")
	s.Equal(repoCnt-1, len(testStorage.Remotes["github.com"].Repos), "github.com remote should have one less repo")
	results := testStorage.Search("dotfiles")
	s.Equal(len(results), 0, "There should no longer be a dotfiles repo")
}

func (s *StorageSuite) TestRemoveRepoNoRepo() {
	repoCnt := len(testStorage.Remotes["github.com"].Repos)
	err := testStorage.RemoveRepository("github.com/TheHipbot/docker")
	s.NotNil(err, "There should be an error returned")
	s.Equal(repoCnt, len(testStorage.Remotes["github.com"].Repos), "github.com remote should have the same number of repos")
}

func (s *StorageSuite) TestRemoveRepoAndSave() {
	repoCnt := len(testStorage.Remotes["github.com"].Repos)
	testStorage.Save()
	testStorage.RemoveRepository("github.com/TheHipbot/dotfiles")
	testStorage.Save()
	cache := storage{
		storer: testStorer,
	}
	cache.Open()
	s.Equal(repoCnt-1, len(cache.Remotes["github.com"].Repos), "github.com remote should have one less repo")
	results := cache.Search("dotfiles")
	s.Equal(len(results), 0, "There should no longer be a dotfiles repo")
}

func (s *StorageSuite) TestStorageSearchWithResults() {
	var results []Repository

	results = testStorage.Search("files")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/dotfiles")
	s.Equal(results[1].Name, "github.com/TheHipbot/dockerfiles")

	results = testStorage.Search("gitlab")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "gitlab.com/gitlab-org/gitlab-ce")
	s.Equal(results[1].Name, "gitlab.com/gnachman/iterm2")
}

func (s *StorageSuite) TestStorageSearchCaseInSensitiveWithResults() {
	var results []Repository

	results = testStorage.Search("FILES")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/dotfiles")
	s.Equal(results[1].Name, "github.com/TheHipbot/dockerfiles")

	results = testStorage.Search("thehipbot")
	s.Len(results, 3, "There are 3 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/hermes")
	s.Equal(results[1].Name, "github.com/TheHipbot/dotfiles")
	s.Equal(results[2].Name, "github.com/TheHipbot/dockerfiles")
}

func (s *StorageSuite) TestStorageSearchWithoutResults() {
	var results []Repository

	results = testStorage.Search("test")
	s.Len(results, 0, "There no results")

	testStorage.Remotes = map[string]*Remote{}
	results = testStorage.Search("files")
	s.Len(results, 0, "There no results")
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}
