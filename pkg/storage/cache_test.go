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
				Repos: map[string]*Repository{
					"github.com/TheHipbot/hermes": &Repository{
						Name: "github.com/TheHipbot/hermes",
						Path: "/repos/github.com/TheHipbot/hermes",
					},
					"github.com/TheHipbot/dotfiles": &Repository{
						Name: "github.com/TheHipbot/dotfiles",
						Path: "/repos/github.com/TheHipbot/dotfiles",
					},
					"github.com/TheHipbot/dockerfiles": &Repository{
						Name: "github.com/TheHipbot/dockerfiles",
						Path: "/repos/github.com/TheHipbot/dockerfiles",
					},
					"github.com/src-d/go-git": &Repository{
						Name: "github.com/src-d/go-git",
						Path: "/repos/github.com/src-d/go-git",
					},
				},
			},
			"gitlab.com": &Remote{
				Name: "gitlab.com",
				URL:  gitLabURL.String(),
				Repos: map[string]*Repository{
					"gitlab.com/gitlab-org/gitlab-ce": &Repository{
						Name: "gitlab.com/gitlab-org/gitlab-ce",
						Path: "/repos/gitlab.com/gitlab-org/gitlab-ce",
					},
					"gitlab.com/gnachman/iterm2": &Repository{
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
				"repos": {
					"github.com/TheHipbot/hermes": {
						"name": "github.com/TheHipbot/hermes",
						"repo_path": "/repos/github.com/TheHipbot/hermes"
					},
					"github.com/TheHipbot/dotfiles": {
						"name": "github.com/TheHipbot/dotfiles",
						"repo_path": "/repos/github.com/TheHipbot/dotfiles"
					},
					"github.com/TheHipbot/dockerfiles": {
						"name": "github.com/TheHipbot/dockerfiles",
						"repo_path": "/repos/github.com/TheHipbot/dockerfiles"
					},
					"github.com/src-d/go-git": {
						"name": "github.com/src-d/go-git",
						"repo_path": "/repos/github.com/src-d/go-git"
					}
				}
			},
			"gitlab.com": {
				"name": "gitlab.com",
				"url":  "https://gitlab.com",
				"repos": {
					"gitlab.com/gitlab-org/gitlab-ce": {
						"name": "gitlab.com/gitlab-org/gitlab-ce",
						"repo_path": "/repos/gitlab.com/gitlab-org/gitlab-ce"
					},
					"gitlab.com/gnachman/iterm2": {
						"name": "gitlab.com/gnachman/iterm2",
						"repo_path": "/repos/gitlab.com/gnachman/iterm2"
					}
				}
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
	s.Equal(storage.Remotes["github.com"].Repos["github.com/TheHipbot/hermes"].Name, "github.com/TheHipbot/hermes", "github.com/TheHipbot/herme should be a repo in the github.com remote should be hermes")
	s.NotNil(storage.Remotes["gitlab.com"], "There should be repos in the gitlab.com remote")
	s.Equal(len(storage.Remotes["gitlab.com"].Repos), 2, "There should be 4 repos in the gitlab.com remote")
	s.Equal(storage.Remotes["gitlab.com"].Repos["gitlab.com/gitlab-org/gitlab-ce"].Name, "gitlab.com/gitlab-org/gitlab-ce", "gitlab.com/gitlab-org/gitlab-ce should be a repo in the gitlab.com remote should be gitlab-ce")
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
	testStorage.AddRepository(&Repository{
		Name: "github.com/TheHipbot/weather",
		Path: "/repos/",
	})
	results = testStorage.SearchRepositories("weather")
	s.Len(results, 1, "There should be the new repo")
	s.Equal(repoCnt+1, len(testStorage.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")

	testStorage.AddRepository(&Repository{
		Name: "github.com/TheHipbot/docker",
		Path: "/repos/",
	})
	results = testStorage.SearchRepositories("docker")
	s.Len(results, 2, "There should be the new repo")
	s.Equal(repoCnt+2, len(testStorage.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
}

func (s *StorageSuite) TestCacheAddSameRepoMultipleTimes() {
	var results []Repository

	repoCnt := len(testStorage.Remotes["github.com"].Repos)
	testStorage.AddRepository(&Repository{
		Name: "github.com/TheHipbot/weather",
		Path: "/repos/",
	})
	results = testStorage.SearchRepositories("weather")
	s.Len(results, 1, "There should be the new repo")
	s.Equal(repoCnt+1, len(testStorage.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")

	testStorage.AddRepository(&Repository{
		Name: "github.com/TheHipbot/weather",
		Path: "/repos/",
	})
	results = testStorage.SearchRepositories("weather")
	s.Len(results, 1, "There should still only be one entry")
}

func (s *StorageSuite) TestCacheAddNewRemote() {
	var results []Repository

	remote := testStorage.Remotes["gopkg.in"]
	s.Nil(remote, "The remote should not exist")
	testStorage.AddRepository(&Repository{
		Name: "gopkg.in/src-d/go-billy.v4",
		Path: "/repos/",
	})
	results = testStorage.SearchRepositories("billy")
	s.Len(results, 1, "There should be the new repo")
	remote = testStorage.Remotes["gopkg.in"]
	s.NotNil(remote, "The remote should exist")
}

func (s *StorageSuite) TestCacheAddThenSave() {
	var results []Repository

	repoCnt := len(testStorage.Remotes["github.com"].Repos)
	testStorage.AddRepository(&Repository{
		Name: "github.com/TheHipbot/weather",
		Path: "/repos/",
	})
	results = testStorage.SearchRepositories("weather")
	s.Len(results, 1, "There should be the new repo")
	s.Equal(repoCnt+1, len(testStorage.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
	err := testStorage.Save()
	s.Nil(err, "Should save")
	var temp storage
	raw, err := ioutil.ReadAll(testStorer)
	s.Nil(err, "Should be unmarshallable")
	s.Nil(json.Unmarshal(raw, &temp), "Should be unmarshallable")
	s.Equal(repoCnt+1, len(temp.Remotes["github.com"].Repos), "The new repo should be stored with existing remote in cache")

	testStorage.AddRepository(&Repository{
		Name: "github.com/TheHipbot/docker",
		Path: "/repos/",
	})
	results = testStorage.SearchRepositories("docker")
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
	results := testStorage.SearchRepositories("dotfiles")
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
	results := cache.SearchRepositories("dotfiles")
	s.Equal(len(results), 0, "There should no longer be a dotfiles repo")
}

func (s *StorageSuite) TestStorageSearchWithResults() {
	var results []Repository

	// Results should be in alphabetical order by repo name
	results = testStorage.SearchRepositories("files")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/dockerfiles")
	s.Equal(results[1].Name, "github.com/TheHipbot/dotfiles")

	results = testStorage.SearchRepositories("gitlab")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "gitlab.com/gitlab-org/gitlab-ce")
	s.Equal(results[1].Name, "gitlab.com/gnachman/iterm2")
}

func (s *StorageSuite) TestStorageSearchCaseInSensitiveWithResults() {
	var results []Repository

	// Results should be in alphabetical order by repo name
	results = testStorage.SearchRepositories("FILES")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/dockerfiles")
	s.Equal(results[1].Name, "github.com/TheHipbot/dotfiles")

	results = testStorage.SearchRepositories("thehipbot")
	s.Len(results, 3, "There are 3 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/dockerfiles")
	s.Equal(results[1].Name, "github.com/TheHipbot/dotfiles")
	s.Equal(results[2].Name, "github.com/TheHipbot/hermes")
}

func (s *StorageSuite) TestStorageSearchWithoutResults() {
	var results []Repository

	results = testStorage.SearchRepositories("test")
	s.Len(results, 0, "There no results")

	testStorage.Remotes = map[string]*Remote{}
	results = testStorage.SearchRepositories("files")
	s.Len(results, 0, "There no results")
}

func (s *StorageSuite) TestStorageSearchRemote() {
	res, ok := testStorage.SearchRemote("github.com")
	s.True(ok)
	s.Equal("github.com", res.Name)
	s.Equal("https://github.com", res.URL)
}

func (s *StorageSuite) TestStorageSearchRemoteWithoutResult() {
	res, ok := testStorage.SearchRemote("dne.com")
	s.False(ok)
	s.Equal(&Remote{}, res)
}

func (s *StorageSuite) TestListRemotes() {
	remotes := testStorage.ListRemotes()

	s.Len(remotes, 2, "There should be 2 remotes")
	for _, r := range remotes {
		s.Contains([]string{"github.com", "gitlab.com"}, r.Name, "It should be one of the two remotes in test cache")
	}
}

func (s *StorageSuite) TestListRemotesAfterAdd() {
	remotes := testStorage.ListRemotes()
	s.Len(remotes, 2, "There should be 2 remotes")

	err := testStorage.AddRemote("test.com", "test", "test", "https")
	s.Nil(err, "Add remote should not error")
	remotes = testStorage.ListRemotes()
	s.Len(remotes, 3, "There should be 3 remotes after adding one")
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}
