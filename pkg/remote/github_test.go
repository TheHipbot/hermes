package remote

import (
	"strings"
	"testing"

	"github.com/google/go-github/v29/github"
	"github.com/stretchr/testify/suite"
)

type GitHubRemoteSuite struct {
	suite.Suite
}

func (s *GitHubRemoteSuite) TestGitHubCreator() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := githubCreator(opts)
	s.Nil(err, "Should return without error")
	s.IsType(d, &GitHub{}, "Driver should be GitHub type")
	gh := d.(*GitHub)
	s.Equal(defaultGitHubAPIHost, gh.Host, "GitHub Driver should be created with default host")
	s.Equal(opts, gh.Opts, "Options should be passed into Github struct")
}

func (s *GitHubRemoteSuite) TestSetHost() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := githubCreator(opts)
	s.Nil(err, "Creator should not return error")
	gh := d.(*GitHub)
	d.SetHost("http://test.github.com")
	s.Equal("http://test.github.com", gh.Host, "Host should be set")
}

func (s *GitHubRemoteSuite) TestSetHostToGithub() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := githubCreator(opts)
	s.Nil(err, "Creator should not return error")
	gh := d.(*GitHub)
	d.SetHost("https://github.com")
	s.Equal("https://api.github.com", gh.Host, "Host should be set to default")
	d.SetHost("github.com")
	s.Equal("https://api.github.com", gh.Host, "Host should be set to default")
}

func (s *GitHubRemoteSuite) TestGitHubSetAuth() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := githubCreator(opts)
	s.Nil(err, "Creator should not return error")
	testAuth := Auth{
		Token: "1234abc",
	}
	gh := d.(*GitHub)
	d.Authenticate(testAuth)
	s.Equal(testAuth, gh.Auth, "Auth should be set")
}

func (s *GitHubRemoteSuite) TestGithubMapper() {
	res := []map[string]string{}

	htmlURL1 := "https://github.com/carsdotcom/beacon"
	cloneURL1 := "https://github.com/carsdotcom/beacon.git"
	sshURL1 := "git@github.com:carsdotcom/beacon.git"
	htmlURL2 := "https://github.com/carsdotcom/bitcar"
	cloneURL2 := "https://github.com/carsdotcom/bitcar.git"
	sshURL2 := "git@github.com:carsdotcom/bitcar.git"
	testRepos := []*github.Repository{
		&github.Repository{
			HTMLURL:  &htmlURL1,
			CloneURL: &cloneURL1,
			SSHURL:   &sshURL1,
		},
		&github.Repository{
			HTMLURL:  &htmlURL2,
			CloneURL: &cloneURL2,
			SSHURL:   &sshURL2,
		},
	}
	res, err := mapGitHubRepos(res, testRepos)
	s.Nil(err)
	s.Equal(res[0]["url"], htmlURL1)
	s.Equal(res[0]["name"], strings.Split(htmlURL1, "://")[1])
	s.Equal(res[0]["ssh_url"], sshURL1)
	s.Equal(res[0]["clone_url"], cloneURL1)
	s.Equal(res[1]["url"], htmlURL2)
	s.Equal(res[1]["name"], strings.Split(htmlURL2, "://")[1])
	s.Equal(res[1]["ssh_url"], sshURL2)
	s.Equal(res[1]["clone_url"], cloneURL2)
}

func (s *GitHubRemoteSuite) TestGitHubAuthType() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := githubCreator(opts)
	s.Nil(err, "Creator should not return error")
	s.Equal(authToken, d.AuthType(), "AuthType should be authToken")
}

func TestGitHubRemoteSuite(t *testing.T) {
	suite.Run(t, new(GitHubRemoteSuite))
}
