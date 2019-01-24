package repo

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/suite"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

var (
	testReposPath = "/home/test-repos/"
	testRepoName  = "github.com/TheHipbot/hermes"

	testPEMKey = `-----BEGIN DSA PRIVATE KEY-----
MIIBugIBAAKBgQCaqh8rALlKS8TwamPK5jkjkuvPeg0fJhUMM85U9AsfTSRk7SS3
VsCszB8StadHMSTTqzY/vze4Jx69RV6M7takDbCuSB9zQse2O3FWP/x1CDdj7QLL
9sRbceYZ8uwwjUkep/342Acb7qFkCPhvqFI+8sbcV9fY1vb/tSATHJGC3QIVAKH6
v71V3XJe8IcuBct8BRriCBdrAoGAerahuBRXqvPGvMhxA0Yfo1t6SFepalth2Suk
QiCk4uUL6GpyoATYTVBfAZWQlG/3G0XnmGRegp2adKMJMrVr3tAlcRkC4FBZ3OlF
W5tcrkARQ6+bkKGxPPSgTJzp4KoaJQt5diAsHUEYioE4eJq10v0JMYqgLqW6pkJW
oVBCtbwCgYB0bRiyQpx1wnAN1kZY57O9NZmvGUwbaRy9OKhdl+xwV1reh4454ncx
tbdlwS6XNMSWqQ2nAW1unrvahmSiugMALYNRT7Oz9VrMJKg2uymYyLLDU3K9BsJI
IA+rMecwVEVnwVpLe9sFv2Ax/PXYUKmEMntpUZe59iWilzK6Is3vGQIUfJtpNuzI
cAbBwCiZESVy2GNv6V0=
-----END DSA PRIVATE KEY-----
`

	sshConfigFmtStr = `Host %s
	user git
	IdentityFile %s`
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

	repo := NewGitRepository(testRepoName, repoURL.String())
	repo.Fs = appFs

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

	repo := NewGitRepository(testRepoName, repoURL.String())
	repo.Fs = appFs

	suite.Nil(repo.Clone(pathToClone), "Error cloning repo")
	suite.Equal(repo.Clone(pathToClone), git.ErrRepositoryAlreadyExists, "Should throw ErrRepositoryAlreadyExists error")
}

type testCloner struct {
	suite *GitRepositorySuite
}

func (t *testCloner) clone(storer *filesystem.Storage, tree billy.Filesystem, opts *git.CloneOptions) error {
	t.suite.NotNil(opts.Auth)
	return nil
}

func (suite *GitRepositorySuite) TestCloneSSHConfig() {
	pathToClone := fmt.Sprintf("%s%s", testReposPath, testRepoName)
	repoURL := "git@github.com:/TheHipbot/hermes"

	repo := NewGitRepository(testRepoName, repoURL)
	repo.Fs = appFs
	repo.Protocol = "ssh"
	repo.cloner = &testCloner{
		suite: suite,
	}

	fullPath, _ := homedir.Dir()

	appFs.MkdirAll(fmt.Sprintf("%s/.ssh/", fullPath), os.ModeDir)
	writePEMFile("~/.ssh/id_rsa")
	writeSSHConfigFile("github.com", "~/.ssh/id_rsa")
	suite.Nil(repo.Clone(pathToClone), "Error cloning repo")
}

func (suite *GitRepositorySuite) TestCloneSSHConfigWithPort() {
	pathToClone := fmt.Sprintf("%s%s", testReposPath, "gitlab.hipbot.com/TheHipbot/hermes")
	repoURL, err := url.Parse("ssh://git@gitlab.hipbot.com:8893/TheHipbot/hermes")
	suite.Nil(err, "Test URL could not be parsed")

	repo := NewGitRepository(testRepoName, repoURL.String())
	repo.Fs = appFs
	repo.Protocol = "ssh"
	repo.cloner = &testCloner{
		suite: suite,
	}

	fullPath, _ := homedir.Dir()

	appFs.MkdirAll(fmt.Sprintf("%s/.ssh/", fullPath), os.ModeDir)
	writePEMFile("~/.ssh/id_rsa")
	writeSSHConfigFile("gitlab.hipbot.com", "~/.ssh/id_rsa")
	suite.Nil(repo.Clone(pathToClone), "Error cloning repo")
}

func TestGitRepositorySuite(t *testing.T) {
	suite.Run(t, new(GitRepositorySuite))
}

func writePEMFile(path string) error {
	fullPath, err := homedir.Expand(path)
	if err != nil {
		return err
	}
	file, err := appFs.Create(fullPath)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(testPEMKey))
	return err
}

func writeSSHConfigFile(host, path string) error {
	fullPath, err := homedir.Dir()
	if err != nil {
		return err
	}
	file, err := appFs.Create(fmt.Sprintf("%s/.ssh/config", fullPath))
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(fmt.Sprintf(sshConfigFmtStr, host, path)))
	return err
}
