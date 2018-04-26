package repo

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-billy.v4/memfs"
)

var (
	testReposPath = "/home/test-repos/"
	testRepoName  = "github.com/TheHipbot/hermes"
)

func init() {
	appFs = memfs.New()
}

func TestCloneRepo(t *testing.T) {
	assert := assert.New(t)
	pathToClone := fmt.Sprintf("%s%s", testReposPath, testRepoName)
	repoURL, err := url.Parse("https://github.com/TheHipbot/hermes")
	if err != nil {
		t.Fatalf("Test URL could not be parsed")
	}

	repo := GitRepository{
		Name: testRepoName,
		URL:  repoURL,
	}

	assert.Nil(repo.Clone(pathToClone), "Error cloning repo")

	// is there a directory in the memfs for the cloned repo
	fileInfo, err := appFs.Stat(pathToClone)
	assert.Nil(err, fmt.Sprintf("Error getting directory %s stat", pathToClone))
	assert.True(fileInfo.IsDir(), "Repo path should be a directory")

	// is .git a directory
	gitPath := fmt.Sprintf("%s/.git", pathToClone)
	fileInfo, err = appFs.Stat(gitPath)
	assert.Nil(err, fmt.Sprintf("Error getting directory %s stat", gitPath))
	assert.True(fileInfo.IsDir(), fmt.Sprintf("%s path should be a directory", gitPath))

	// is there a README and main.go
	fileInfo, err = appFs.Stat(fmt.Sprintf("%s/README.md", pathToClone))
	assert.Nil(err, "Error getting README stat")
	assert.True(fileInfo.Mode().IsRegular(), "README is missing")

	fileInfo, err = appFs.Stat(fmt.Sprintf("%s/main.go", pathToClone))
	assert.Nil(err, "Error getting main stat")
	assert.True(fileInfo.Mode().IsRegular(), "main.go is missing")
}

func TestCloneExistingRepo(t *testing.T) {
	assert := assert.New(t)

	pathToClone := fmt.Sprintf("%s%s", testReposPath, testRepoName)
	repoURL, err := url.Parse("https://github.com/TheHipbot/hermes")
	if err != nil {
		t.Fatalf("Test URL could not be parsed")
	}

	repo := GitRepository{
		Name: testRepoName,
		URL:  repoURL,
	}

	assert.Nil(repo.Clone(pathToClone), "Error cloning repo")
	assert.Nil(repo.Clone(pathToClone), "Error cloning repo second time")
}
