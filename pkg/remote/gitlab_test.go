package remote

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	gitlab "github.com/xanzy/go-gitlab"
)

type GitLabRemoteSuite struct {
	suite.Suite
}

func (s *GitLabRemoteSuite) TestGitlabCreator() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := gitlabCreator(opts)
	s.Nil(err, "Should return without error")
	s.IsType(d, &GitLab{}, "Driver should be Gitlab type")
	gl := d.(*GitLab)
	s.Equal(defaultGitlabAPIHost, gl.Host, "Gitlab Driver should be created with default host")
}

func (s *GitLabRemoteSuite) TestGitlabSetHost() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := gitlabCreator(opts)
	s.Nil(err, "Creator should not return error")
	gl := d.(*GitLab)
	d.SetHost("http://test.gitlab.com")
	s.Equal("http://test.gitlab.com", gl.Host, "Host should be set")
}

func (s *GitLabRemoteSuite) TestGitlabSetAuth() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := gitlabCreator(opts)
	s.Nil(err, "Creator should not return error")
	testAuth := Auth{
		Token: "1234abc",
	}
	gl := d.(*GitLab)
	d.Authenticate(testAuth)
	s.Equal(testAuth, gl.Auth, "Auth should be set")
}

func (s *GitLabRemoteSuite) TestGitlabAuthType() {
	opts := &DriverOpts{
		Auth: &Auth{
			Token: "abcd123",
			Type:  "token",
		},
	}
	d, err := gitlabCreator(opts)
	s.Nil(err, "Creator should not return error")
	s.Equal(authToken, d.AuthType(), "AuthType should be authToken")
}

func (s *GitLabRemoteSuite) TestGitLabMapper() {
	res := []map[string]string{}

	htmlURL1 := "https://github.com/carsdotcom/beacon"
	cloneURL1 := "https://github.com/carsdotcom/beacon.git"
	sshURL1 := "git@github.com:carsdotcom/beacon.git"
	htmlURL2 := "https://github.com/carsdotcom/bitcar"
	cloneURL2 := "https://github.com/carsdotcom/bitcar.git"
	sshURL2 := "git@github.com:carsdotcom/bitcar.git"
	testRepos := []*gitlab.Project{
		&gitlab.Project{
			WebURL:        htmlURL1,
			HTTPURLToRepo: cloneURL1,
			SSHURLToRepo:  sshURL1,
		},
		&gitlab.Project{
			WebURL:        htmlURL2,
			HTTPURLToRepo: cloneURL2,
			SSHURLToRepo:  sshURL2,
		},
	}
	res, err := mapGitLabProjects(res, testRepos)
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

func TestGitLabRemoteSuite(t *testing.T) {
	suite.Run(t, new(GitLabRemoteSuite))
}
