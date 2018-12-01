package cmd

import (
	"testing"

	"github.com/TheHipbot/hermes/mock"
	"github.com/TheHipbot/hermes/remote"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type RemoteCmdSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	mockDriver *mock.MockDriver
}

var (
	mockDriver   *remote.Driver
	ctrl         *gomock.Controller
	testRepoPath = "/home/user/test-repos/"
	mockCmd      = &cobra.Command{}
)

func (suite *RemoteCmdSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDriver = mock.NewMockDriver(suite.ctrl)
	drivers = append(drivers, driver{
		Name: "test",
	})
	remote.RegisterDriver("test", func() (remote.Driver, error) {
		return suite.mockDriver, nil
	})
	viper.Set("repo_path", testRepoPath)
}

func (suite *RemoteCmdSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *RemoteCmdSuite) TestWithTokenAuth() {
	ctrl := gomock.NewController(suite.T())
	mockPrompter := mock.NewMockFactory(ctrl)
	prompter = mockPrompter
	mockSelectPrompt := mock.NewMockSelectPrompt(ctrl)
	mockInputPrompt := mock.NewMockInputPrompt(ctrl)
	mockStore := mock.NewMockStorage(ctrl)
	defer ctrl.Finish()

	store = mockStore
	repos := []map[string]string{
		{
			"name": "github.com/thehipbot/hermes",
			"url":  "https://github.com/thehipbot/hermes",
		},
		{
			"name": "github.com/thehipbot/dotfiles",
			"url":  "https://github.com/thehipbot/dotfiles",
		},
		{
			"name": "github.com/carsdotcom/bitcar",
			"url":  "https://github.com/carsdotcom/bitcar",
		},
	}

	gomock.InOrder(
		// create prompt for drivers
		mockPrompter.
			EXPECT().
			CreateSelectPrompt(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockSelectPrompt).
			Times(1),

		// run driver prompt
		mockSelectPrompt.
			EXPECT().
			Run().
			Return(2, "test", nil).
			Times(1),

		suite.mockDriver.
			EXPECT().
			AuthType().
			Return("token").
			Times(1),

		// prompt for token
		mockPrompter.
			EXPECT().
			CreateInputPrompt(gomock.Any()).
			Return(mockInputPrompt).
			Times(1),

		// run token prompt
		mockInputPrompt.
			EXPECT().
			Run().
			Return("1234abcd", nil).
			Times(1),

		suite.mockDriver.
			EXPECT().
			Authenticate(gomock.Eq(remote.Auth{
				Token: "1234abcd",
			})).
			Times(1),

		suite.mockDriver.
			EXPECT().
			GetRepos().
			Return(repos, nil).
			Times(1),

		mockStore.
			EXPECT().
			Open().
			Return().
			Times(1),

		mockStore.
			EXPECT().
			AddRemote("https://github.com", "github.com").
			Return(nil).
			Times(1),
	)

	mockStore.
		EXPECT().
		AddRepository("github.com/thehipbot/hermes", testRepoPath).
		Return(nil).
		Times(1)

	mockStore.
		EXPECT().
		AddRepository("github.com/thehipbot/dotfiles", testRepoPath).
		Return(nil).
		Times(1)

	mockStore.
		EXPECT().
		AddRepository("github.com/carsdotcom/bitcar", testRepoPath).
		Return(nil).
		Times(1)

	gomock.InOrder(
		mockStore.
			EXPECT().
			Save().
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			Close().
			Return(nil).
			Times(1),
	)

	remoteAddHandler(mockCmd, []string{"github.com"})
}

func TestRemoteCmdSuite(t *testing.T) {
	suite.Run(t, new(RemoteCmdSuite))
}
