package cache

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CacheSuite struct {
	suite.Suite
}

func (s *CacheSuite) SetupTest() {
	githubURL, _ := url.Parse("https://github.com")
	gitLabURL, _ := url.Parse("https://gitlab.com")
	cache.Remotes = map[string]*Remote{
		"github.com": &Remote{
			Name: "github.com",
			URL:  githubURL,
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
			URL:  gitLabURL,
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
	}
}

func (s *CacheSuite) TestCacheSearchWithResults() {
	var results []Repo

	results = Search("files")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "github.com/TheHipbot/dotfiles")
	s.Equal(results[1].Name, "github.com/TheHipbot/dockerfiles")

	results = Search("gitlab")
	s.Len(results, 2, "There are 2 repos with files in the name")
	s.Equal(results[0].Name, "gitlab.com/gitlab-org/gitlab-ce")
	s.Equal(results[1].Name, "gitlab.com/gnachman/iterm2")
}

func (s *CacheSuite) TestCacheSearchWithoutResults() {
	var results []Repo

	results = Search("test")
	s.Len(results, 0, "There no results")

	cache.Remotes = map[string]*Remote{}
	results = Search("files")
	s.Len(results, 0, "There no results")
}

func (s *CacheSuite) TestCacheAdd() {
	var results []Repo

	repoCnt := len(cache.Remotes["github.com"].Repos)
	Add("github.com/TheHipbot/weather", "/repos/")
	results = Search("weather")
	s.Len(results, 1, "There should be the new repo")
	s.Equal(repoCnt+1, len(cache.Remotes["github.com"].Repos), "The new repo should be stored with existing remote")
}

func (s *CacheSuite) TestCacheAddNewRemote() {
	var results []Repo

	remote := cache.Remotes["gopkg.in"]
	s.Nil(remote, "The remote should not exist")
	Add("gopkg.in/src-d/go-billy.v4", "/repos/")
	results = Search("billy")
	s.Len(results, 1, "There should be the new repo")
	remote = cache.Remotes["gopkg.in"]
	s.NotNil(remote, "The remote should exist")
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
