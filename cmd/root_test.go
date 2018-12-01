package cmd

import (
	"fmt"
	"testing"

	"gopkg.in/src-d/go-billy.v4"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/golang/mock/gomock"

	"github.com/TheHipbot/hermes/fs"
	mock_prompt "github.com/TheHipbot/hermes/mock"
	"github.com/TheHipbot/hermes/pkg/storage"
	"github.com/stretchr/testify/suite"
)

var (
	cmd       = &cobra.Command{}
	cacheFile billy.File
)

type RootCmdSuite struct {
	suite.Suite
}

func (s *RootCmdSuite) SetupTest() {
	configFS = &fs.ConfigFS{
		FS: memfs.New(),
	}
	configFS.Setup()
	cacheFile, _ = configFS.GetCacheFile()
	appFs = memfs.New()
	store = storage.NewStorage(cacheFile)
	viper.Set("repo_path", "/repos") 
}

func (s *RootCmdSuite) TearDownSuite() {
	store.Close()
}

func (s *RootCmdSuite) TestGetHandlerSingleCachedRepo() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	cacheFile.Seek(0, 0)
	p, _ := cacheFile.Write([]byte(`{
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
	cacheFile.Truncate(int64(p))
	store.Open()

	mockPrompter := mock_prompt.NewMockFactory(ctrl)
	mockPrompter.
		EXPECT().
		CreateSelectPrompt(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)
	prompter = mockPrompter

	getHandler(cmd, []string{"github.com/TheHipbot/hermes"})
	target := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("target_file"))
	stat, _ := configFS.FS.Stat(target)
	gitFileStat, err := appFs.Stat(fmt.Sprintf("%s/%s", viper.GetString("repo_path"), "github.com/TheHipbot/hermes"))
	
	s.Nil(err)
	targetFile, err := configFS.FS.Open(target)
	defer targetFile.Close()
	content := make([]byte, stat.Size())
	targetFile.Read(content)
	s.Nil(err, "Target file should exist")
	s.Equal(string(content), "/repos/github.com/TheHipbot/hermes", "Get should find one repo and set target path")
	s.True(gitFileStat.IsDir(), ".git folder in repo should exist")
}

func (s *RootCmdSuite) TestGetHandlerMultipleCachedRepos() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	cacheFile.Seek(0, 0)
	p, _ := cacheFile.Write([]byte(`{
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
	cacheFile.Truncate(int64(p))
	store.Open()
	repos := []storage.Repository{
		storage.Repository{
			Name: "github.com/TheHipbot/hermes",
			Path: "/repos/github.com/TheHipbot/hermes",
		},
		storage.Repository{
			Name: "github.com/TheHipbot/dotfiles",
			Path: "/repos/github.com/TheHipbot/dotfiles",
		},
		storage.Repository{
			Name: "github.com/TheHipbot/dockerfiles",
			Path: "/repos/github.com/TheHipbot/dockerfiles",
		},
	}

	mockPrompter := mock_prompt.NewMockFactory(ctrl)
	mockPrompt := mock_prompt.NewMockSelectPrompt(ctrl)
	mockPrompt.
		EXPECT().
		Run().
		Return(1, "blah", nil).
		Times(1)

	mockPrompter.
		EXPECT().
		CreateSelectPrompt(gomock.Any(), gomock.Eq(repos), gomock.Any()).
		Return(mockPrompt).
		Times(1)

	prompter = mockPrompter

	getHandler(cmd, []string{"hipbot"})
	target := fmt.Sprintf("%s%s", viper.GetString("config_path"), viper.GetString("target_file"))
	stat, _ := configFS.FS.Stat(target)
	targetFile, err := configFS.FS.Open(target)
	defer targetFile.Close()
	content := make([]byte, stat.Size())
	targetFile.Read(content)
	s.Nil(err, "Target file should exist")
	s.Equal(string(content), "/repos/github.com/TheHipbot/dotfiles", "Get should find one repo and set target path")
}

func TestRootCmdSuite(t *testing.T) {
	suite.Run(t, new(RootCmdSuite))
}
