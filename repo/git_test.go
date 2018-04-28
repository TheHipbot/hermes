package repo

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
)

var (
	testReposPath = "/home/test-repos/"
	testRepoName  = "github.com/TheHipbot/hermes"
)

type GitRepositorySuite struct {
	suite.Suite
}

func (suite *GitRepositorySuite) SetupTest() {
	appFs = memfs.New()
}

func (suite *GitRepositorySuite) TestCloneRepo() {
	pathToClone := fmt.Sprintf("%s%s", testReposPath, testRepoName)
	repoURL, err := url.Parse("https://github.com/TheHipbot/hermes")
	suite.Nil(err, "Test URL could not be parsed")

	repo := GitRepository{
		Name: testRepoName,
		URL:  repoURL,
	}

	suite.Nil(repo.Clone(pathToClone), "Error cloning repo")

	// is there a directory in the memfs for the cloned repo
	fileInfo, err := appFs.Stat(pathToClone)
	suite.Nil(err, fmt.Sprintf("Error getting directory %s stat", pathToClone))
	suite.True(fileInfo.IsDir(), "Repo path should be a directory")

	// is .git a directory
	gitPath := fmt.Sprintf("%s/.git", pathToClone)
	fileInfo, err = appFs.Stat(gitPath)
	suite.Nil(err, fmt.Sprintf("Error getting directory %s stat", gitPath))
	suite.True(fileInfo.IsDir(), fmt.Sprintf("%s path should be a directory", gitPath))

	// is there a README and main.go
	fileInfo, err = appFs.Stat(fmt.Sprintf("%s/README.md", pathToClone))
	suite.Nil(err, "Error getting README stat")
	suite.True(fileInfo.Mode().IsRegular(), "README is missing")

	fileInfo, err = appFs.Stat(fmt.Sprintf("%s/main.go", pathToClone))
	suite.Nil(err, "Error getting main stat")
	suite.True(fileInfo.Mode().IsRegular(), "main.go is missing")
}

func (suite *GitRepositorySuite) TestCloneExistingRepo() {
	pathToClone := fmt.Sprintf("%s%s", testReposPath, testRepoName)
	repoURL, err := url.Parse("https://github.com/TheHipbot/hermes")
	suite.Nil(err, "Test URL could not be parsed")

	repo := GitRepository{
		Name: testRepoName,
		URL:  repoURL,
	}

	suite.Nil(repo.Clone(pathToClone), "Error cloning repo")
	suite.Equal(repo.Clone(pathToClone), git.ErrRepositoryAlreadyExists, "Should throw ErrRepositoryAlreadyExists error")
}

func TestGitRepositorySuite(t *testing.T) {
	suite.Run(t, new(GitRepositorySuite))
}
