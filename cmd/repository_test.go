package cmd

import (
	"fmt"
	"testing"

	"github.com/TheHipbot/hermes/pkg/fs"
	"github.com/TheHipbot/hermes/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gopkg.in/src-d/go-billy.v4/memfs"
)

type RepoRmCmdSuite struct {
	suite.Suite
}

func (suite *RepoRmCmdSuite) SetupTest() {
	hardRmFlg = false
	configFS = &fs.ConfigFS{
		FS: memfs.New(),
	}
	configFS.Setup()
	cacheFile, _ = configFS.GetCacheFile()
	appFs = memfs.New()
	store = storage.NewStorage(cacheFile)
	viper.Set("repo_path", "/repos")
}

func (suite *RepoRmCmdSuite) TestRm() {
	store.Open()
	err := store.AddRemote("https://github.com", "github.com", "github", "https")
	suite.Nil(err, "Remote should be added to store without error")
	suite.setupTestRepo(&storage.Repository{
		Name:     "github.com/TheHipbot/hermes",
		Path:     "/repos/github.com/TheHipbot/hermes",
		CloneURL: "https://github.com/TheHipbot/hermes",
	})
	repoRmCommand.Run(&cobra.Command{}, []string{"github.com/TheHipbot/hermes"})
	repos := store.SearchRepositories("github.com/TheHipbot/hermes")
	suite.Empty(repos)
	stat, err := appFs.Stat("/repos/github.com/TheHipbot/hermes")
	suite.Nil(err)
	suite.True(stat.IsDir(), "Repo directory should still be present")
}

func (suite *RepoRmCmdSuite) TestHardRm() {
	store.Open()
	hardRmFlg = true
	err := store.AddRemote("https://github.com", "github.com", "github", "https")
	suite.Nil(err, "Remote should be added to store without error")
	suite.setupTestRepo(&storage.Repository{
		Name:     "github.com/TheHipbot/hermes",
		Path:     "/repos/github.com/TheHipbot/hermes",
		CloneURL: "https://github.com/TheHipbot/hermes",
	})
	suite.setupTestRepo(&storage.Repository{
		Name:     "github.com/spf13/cobra",
		Path:     "/repos/github.com/spf13/cobra",
		CloneURL: "https://github.com/spf13/cobra",
	})
	repoRmCommand.Run(&cobra.Command{}, []string{"github.com/TheHipbot/hermes"})
	repos := store.SearchRepositories("github.com/TheHipbot/hermes")
	suite.Empty(repos)
	repos = store.SearchRepositories("github.com/spf13/cobra")
	suite.NotEmpty(repos)
	stat, err := appFs.Stat("/repos/github.com/TheHipbot/hermes")
	suite.NotNil(err, "Repo directory should not still be present")
	stat, err = appFs.Stat("/repos/github.com/TheHipbot/")
	suite.NotNil(err, "Repo directory should not still be present")
	stat, err = appFs.Stat("/repos/github.com/spf13/cobra")
	suite.Nil(err, "Repo directory should still be present")
	suite.True(stat.IsDir(), "Repo directory should still be present")
}

func (suite *RepoRmCmdSuite) TestHardRmNoOtherRepos() {
	store.Open()
	hardRmFlg = true
	err := store.AddRemote("https://github.com", "github.com", "github", "https")
	suite.Nil(err, "Remote should be added to store without error")
	suite.setupTestRepo(&storage.Repository{
		Name:     "github.com/TheHipbot/hermes",
		Path:     "/repos/github.com/TheHipbot/hermes",
		CloneURL: "https://github.com/TheHipbot/hermes",
	})
	repoRmCommand.Run(&cobra.Command{}, []string{"github.com/TheHipbot/hermes"})
	repos := store.SearchRepositories("github.com/TheHipbot/hermes")
	suite.Empty(repos)
	stat, err := appFs.Stat("/repos/github.com/TheHipbot/hermes")
	suite.NotNil(err, "Repo directory should not still be present")
	stat, err = appFs.Stat("/repos/github.com/TheHipbot/")
	suite.NotNil(err, "Repo directory should not still be present")
	stat, err = appFs.Stat("/repos/")
	suite.Nil(err, "Repo directory should still be present")
	suite.True(stat.IsDir(), "Repo directory should still be present")
}

func (suite *RepoRmCmdSuite) setupTestRepo(repo *storage.Repository) {
	err := store.AddRepository(repo)
	suite.Nil(err, "Repository should be added to store without error")
	suite.Nil(store.Save(), "store should be saved")
	err = appFs.MkdirAll(repo.Path, 777)
	suite.Nil(err, "Repository path should be created")
	_, err = appFs.Create(fmt.Sprintf("%s/test.txt", repo.Path))
	suite.Nil(err, "Repository file should be created")
}

func TestRepoRmSuite(t *testing.T) {
	suite.Run(t, new(RepoRmCmdSuite))
}
