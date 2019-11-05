package cmd

import (
	"fmt"
	"testing"

	"github.com/TheHipbot/hermes/pkg/credentials"

	"github.com/TheHipbot/hermes/mock"
	"github.com/TheHipbot/hermes/pkg/remote"
	"github.com/TheHipbot/hermes/pkg/storage"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type RemoteCmdSuite struct {
	suite.Suite
	ctrl                 *gomock.Controller
	mockDriver           *mock.MockDriver
	mockCredentialStorer *mock.MockCredentialsStorer
}

var (
	mockDriver   *remote.Driver
	ctrl         *gomock.Controller
	testRepoPath = "/home/user/test-repos/"
	mockCmd      = &cobra.Command{}
	optsHarness  *remote.DriverOpts
)

func (suite *RemoteCmdSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDriver = mock.NewMockDriver(suite.ctrl)
	drivers = append(drivers, driver{
		Name: "test",
	})
	remote.RegisterDriver("test", func(opts *remote.DriverOpts) (remote.Driver, error) {
		optsHarness = opts
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
	suite.mockCredentialStorer = mock.NewMockCredentialsStorer(ctrl)
	credentialsStorer = suite.mockCredentialStorer
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

	mockStore.
		EXPECT().
		SearchRemote("github.com").
		Return(nil, false).
		Times(1)

	promptForDriverThenSetHost(mockPrompter, mockSelectPrompt, suite.mockDriver, 2, "https://github.com")

	suite.mockDriver.
		EXPECT().
		AuthType().
		Return("token").
		Times(1)

	suite.mockCredentialStorer.
		EXPECT().
		Get("github.com").
		Return(credentials.Credential{}, credentials.ErrCredentialNotFound).
		Times(1)

	suite.mockCredentialStorer.
		EXPECT().
		Close().
		Return(nil).
		Times(1)

	promptForAuthThenStore(mockPrompter, mockInputPrompt, suite.mockCredentialStorer, suite.mockDriver, "github.com")

	promptForProtocol(mockPrompter, mockSelectPrompt, 0, "https")

	suite.mockDriver.
		EXPECT().
		GetRepos().
		Return(repos, nil).
		Times(1)

	// Open store, add remote, add repos, save, close
	gomock.InOrder(
		mockStore.
			EXPECT().
			Open().
			Return().
			Times(1),

		mockStore.
			EXPECT().
			AddRemote("https://github.com", "github.com", "test", "https").
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/thehipbot/hermes",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/hermes"),
			}).
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/thehipbot/dotfiles",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/dotfiles"),
			}).
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/carsdotcom/bitcar",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/carsdotcom/bitcar"),
			}).
			Return(nil).
			Times(1),
	)

	saveAndCloseStorage(mockStore)

	remoteAddHandler(mockCmd, []string{"https://github.com"})
}

func (suite *RemoteCmdSuite) TestAuthError() {
	ctrl := gomock.NewController(suite.T())
	mockPrompter := mock.NewMockFactory(ctrl)
	prompter = mockPrompter
	mockSelectPrompt := mock.NewMockSelectPrompt(ctrl)
	mockInputPrompt := mock.NewMockInputPrompt(ctrl)
	mockStore := mock.NewMockStorage(ctrl)
	suite.mockCredentialStorer = mock.NewMockCredentialsStorer(ctrl)
	credentialsStorer = suite.mockCredentialStorer
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

	mockStore.
		EXPECT().
		SearchRemote("github.com").
		Return(nil, false).
		Times(1)

	suite.mockDriver.
		EXPECT().
		AuthType().
		Return("token").
		Times(2)

	promptForDriverThenSetHost(mockPrompter, mockSelectPrompt, suite.mockDriver, 2, "https://github.com")

	promptForAuthThenStore(mockPrompter, mockInputPrompt, suite.mockCredentialStorer, suite.mockDriver, "github.com")

	// Open store, add remote, add repos, save, close
	gomock.InOrder(
		mockStore.
			EXPECT().
			Open().
			Return().
			Times(1),

		mockStore.
			EXPECT().
			AddRemote("https://github.com", "github.com", "test", "https").
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/thehipbot/hermes",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/hermes"),
			}).
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/thehipbot/dotfiles",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/dotfiles"),
			}).
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/carsdotcom/bitcar",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/carsdotcom/bitcar"),
			}).
			Return(nil).
			Times(1),
	)

	saveAndCloseStorage(mockStore)

	suite.mockCredentialStorer.
		EXPECT().
		Get("github.com").
		Return(credentials.Credential{}, credentials.ErrCredentialNotFound).
		Times(2)

	suite.mockCredentialStorer.
		EXPECT().
		Close().
		Return(nil).
		Times(1)

	// return auth error from GetRepos
	suite.mockDriver.
		EXPECT().
		GetRepos().
		Return(nil, remote.ErrAuth).
		Times(1)

	gomock.InOrder(
		// delete credential from storer
		suite.mockCredentialStorer.
			EXPECT().
			Delete("github.com").
			Return(nil).
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

		// store auth from input
		suite.mockCredentialStorer.
			EXPECT().
			Put("github.com", credentials.Credential{
				Type:  "token",
				Token: "1234abcd",
			}).
			Return(nil).
			Times(1),

		suite.mockDriver.
			EXPECT().
			Authenticate(gomock.Eq(remote.Auth{
				Token: "1234abcd",
			})).
			Times(1),
	)

	suite.mockDriver.
		EXPECT().
		GetRepos().
		Return(repos, nil).
		Times(1)

	promptForProtocol(mockPrompter, mockSelectPrompt, 0, "https")

	remoteAddHandler(mockCmd, []string{"https://github.com"})
}

func (suite *RemoteCmdSuite) TestWithStoredTokenAuth() {
	ctrl := gomock.NewController(suite.T())
	mockPrompter := mock.NewMockFactory(ctrl)
	prompter = mockPrompter
	mockSelectPrompt := mock.NewMockSelectPrompt(ctrl)
	mockStore := mock.NewMockStorage(ctrl)
	suite.mockCredentialStorer = mock.NewMockCredentialsStorer(ctrl)
	credentialsStorer = suite.mockCredentialStorer
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

	mockStore.
		EXPECT().
		SearchRemote("github.com").
		Return(nil, false).
		Times(1)

	promptForDriverThenSetHost(mockPrompter, mockSelectPrompt, suite.mockDriver, 2, "https://github.com")

	suite.mockDriver.
		EXPECT().
		AuthType().
		Return("token").
		Times(1)

	suite.mockCredentialStorer.
		EXPECT().
		Get("github.com").
		Return(credentials.Credential{
			Type:  "token",
			Token: "1234abcd",
		}, nil).
		Times(1)

	suite.mockCredentialStorer.
		EXPECT().
		Close().
		Return(nil).
		Times(1)

	mockPrompter.
		EXPECT().
		CreateInputPrompt(gomock.Any()).
		Times(0)

	suite.mockDriver.
		EXPECT().
		Authenticate(gomock.Eq(remote.Auth{
			Token: "1234abcd",
		})).
		Times(1)

	promptForProtocol(mockPrompter, mockSelectPrompt, 0, "https")

	suite.mockDriver.
		EXPECT().
		GetRepos().
		Return(repos, nil).
		Times(1)

	// Open store, add remote, add repos, save, close
	gomock.InOrder(
		mockStore.
			EXPECT().
			Open().
			Return().
			Times(1),

		mockStore.
			EXPECT().
			AddRemote("https://github.com", "github.com", "test", "https").
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/thehipbot/hermes",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/hermes"),
			}).
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/thehipbot/dotfiles",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/dotfiles"),
			}).
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/carsdotcom/bitcar",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/carsdotcom/bitcar"),
			}).
			Return(nil).
			Times(1),
	)

	saveAndCloseStorage(mockStore)

	remoteAddHandler(mockCmd, []string{"https://github.com"})
}

func (suite *RemoteCmdSuite) TestRemoteAddWithAll() {
	ctrl := gomock.NewController(suite.T())
	mockPrompter := mock.NewMockFactory(ctrl)
	prompter = mockPrompter
	mockSelectPrompt := mock.NewMockSelectPrompt(ctrl)
	mockInputPrompt := mock.NewMockInputPrompt(ctrl)
	mockStore := mock.NewMockStorage(ctrl)
	suite.mockCredentialStorer = mock.NewMockCredentialsStorer(ctrl)
	credentialsStorer = suite.mockCredentialStorer
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

	mockStore.
		EXPECT().
		Open().
		Return().
		Times(1)

	mockStore.
		EXPECT().
		SearchRemote("github.com").
		Return(nil, false).
		Times(1)

	promptForDriverThenSetHost(mockPrompter, mockSelectPrompt, suite.mockDriver, 2, "https://github.com")

	gomock.InOrder(
		suite.mockDriver.
			EXPECT().
			AuthType().
			Return("token").
			Times(1),

		suite.mockCredentialStorer.
			EXPECT().
			Get("github.com").
			Return(credentials.Credential{}, credentials.ErrCredentialNotFound).
			Times(1),
	)

	promptForAuthThenStore(mockPrompter, mockInputPrompt, suite.mockCredentialStorer, suite.mockDriver, "github.com")

	gomock.InOrder(
		suite.mockDriver.
			EXPECT().
			GetRepos().
			Return(repos, nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRemote("https://github.com", "github.com", "test", "https").
			Return(nil).
			Times(1),

		suite.mockCredentialStorer.
			EXPECT().
			Close().
			Return(nil).
			Times(1),
	)

	promptForProtocol(mockPrompter, mockSelectPrompt, 0, "https")

	mockStore.
		EXPECT().
		AddRepository(&storage.Repository{
			Name: "github.com/thehipbot/hermes",
			Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/hermes"),
		}).
		Return(nil).
		Times(1)

	mockStore.
		EXPECT().
		AddRepository(&storage.Repository{
			Name: "github.com/thehipbot/dotfiles",
			Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/dotfiles"),
		}).
		Return(nil).
		Times(1)

	mockStore.
		EXPECT().
		AddRepository(&storage.Repository{
			Name: "github.com/carsdotcom/bitcar",
			Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/carsdotcom/bitcar"),
		}).
		Return(nil).
		Times(1)

	saveAndCloseStorage(mockStore)

	getAllReposFlg = true
	remoteAddHandler(mockCmd, []string{"https://github.com"})
	suite.True(optsHarness.AllRepos)
}

func (suite *RemoteCmdSuite) TestRemoteAddSSH() {
	ctrl := gomock.NewController(suite.T())
	mockPrompter := mock.NewMockFactory(ctrl)
	prompter = mockPrompter
	mockSelectPrompt := mock.NewMockSelectPrompt(ctrl)
	mockInputPrompt := mock.NewMockInputPrompt(ctrl)
	mockStore := mock.NewMockStorage(ctrl)
	suite.mockCredentialStorer = mock.NewMockCredentialsStorer(ctrl)
	credentialsStorer = suite.mockCredentialStorer
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

	mockStore.
		EXPECT().
		SearchRemote("github.com").
		Return(nil, false).
		Times(1)

	promptForDriverThenSetHost(mockPrompter, mockSelectPrompt, suite.mockDriver, 2, "https://github.com")

	suite.mockDriver.
		EXPECT().
		AuthType().
		Return("token").
		Times(1)

	suite.mockCredentialStorer.
		EXPECT().
		Get("github.com").
		Return(credentials.Credential{}, credentials.ErrCredentialNotFound).
		Times(1)

	suite.mockCredentialStorer.
		EXPECT().
		Close().
		Return(nil).
		Times(1)

	promptForAuthThenStore(mockPrompter, mockInputPrompt, suite.mockCredentialStorer, suite.mockDriver, "github.com")

	promptForProtocol(mockPrompter, mockSelectPrompt, 1, "ssh")

	suite.mockDriver.
		EXPECT().
		GetRepos().
		Return(repos, nil).
		Times(1)

	// Open store, add remote, add repos, save, close
	gomock.InOrder(
		mockStore.
			EXPECT().
			Open().
			Return().
			Times(1),

		mockStore.
			EXPECT().
			AddRemote("https://github.com", "github.com", "test", "ssh").
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/thehipbot/hermes",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/hermes"),
			}).
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/thehipbot/dotfiles",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/thehipbot/dotfiles"),
			}).
			Return(nil).
			Times(1),

		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: "github.com/carsdotcom/bitcar",
				Path: fmt.Sprintf("%s%s", testRepoPath, "github.com/carsdotcom/bitcar"),
			}).
			Return(nil).
			Times(1),
	)

	saveAndCloseStorage(mockStore)

	remoteAddHandler(mockCmd, []string{"https://github.com"})
}

func (suite *RemoteCmdSuite) TestRemoteRefresh() {
	ctrl := gomock.NewController(suite.T())
	mockStore := mock.NewMockStorage(ctrl)
	suite.mockCredentialStorer = mock.NewMockCredentialsStorer(ctrl)
	credentialsStorer = suite.mockCredentialStorer

	defer ctrl.Finish()

	store = mockStore

	githubRepos := []map[string]string{
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
		{
			"name": "github.com/thehipbot/harp",
			"url":  "https://github.com/thehipbot/harp",
		},
	}

	gitlabRepos := []map[string]string{
		{
			"name": "gitlab.com/TheHipbot/test",
			"url":  "https://gitlab.com/TheHipbot/test",
		},
	}

	mockStore.
		EXPECT().
		Open().
		Return().
		Times(1)

	mockStore.
		EXPECT().
		ListRemotes().
		Return([]*storage.Remote{
			&storage.Remote{
				Name:     "github.com",
				URL:      "https://github.com",
				Protocol: "ssh",
				Type:     "test",
			},
			&storage.Remote{
				Name:     "gitlab.com",
				URL:      "https://gitlab.com",
				Protocol: "ssh",
				Type:     "test",
			},
		}).
		Times(1)

	mockStore.
		EXPECT().
		SearchRemote("github.com").
		Return(&storage.Remote{
			Name:     "github.com",
			URL:      "https://github.com",
			Protocol: "ssh",
			Type:     "test",
		}, true).
		Times(1)

	mockStore.
		EXPECT().
		SearchRemote("gitlab.com").
		Return(&storage.Remote{
			Name:     "gitlab.com",
			URL:      "https://gitlab.com",
			Protocol: "ssh",
			Type:     "test",
		}, true).
		Times(1)

	suite.mockDriver.
		EXPECT().
		SetHost(gomock.Eq("https://github.com")).
		Return().
		Times(1)

	suite.mockDriver.
		EXPECT().
		SetHost(gomock.Eq("https://gitlab.com")).
		Return().
		Times(1)

	suite.mockDriver.
		EXPECT().
		GetRepos().
		Return(githubRepos, nil).
		Times(1)

	suite.mockDriver.
		EXPECT().
		GetRepos().
		Return(gitlabRepos, nil).
		Times(1)

	getAuthFromStorer(suite.mockCredentialStorer, suite.mockDriver, "github.com")
	getAuthFromStorer(suite.mockCredentialStorer, suite.mockDriver, "gitlab.com")

	for _, r := range append(githubRepos, gitlabRepos...) {
		mockStore.
			EXPECT().
			AddRepository(&storage.Repository{
				Name: r["name"],
				Path: fmt.Sprintf("%s%s", testRepoPath, r["name"]),
			}).
			Return(nil).
			Times(1)
	}

	saveAndCloseStorage(mockStore)
	suite.mockCredentialStorer.
		EXPECT().
		Close().
		Return(nil).
		Times(1)

	remoteRefreshHandler(mockCmd, []string{})
}

// sets up expects on MockStorage for a save then close
func saveAndCloseStorage(mockStorage *mock.MockStorage) {
	gomock.InOrder(
		mockStorage.
			EXPECT().
			Save().
			Return(nil).
			Times(1),
		mockStorage.
			EXPECT().
			Close().
			Return(nil).
			Times(1),
	)
}

func promptForDriverThenSetHost(
	mockPrompter *mock.MockFactory,
	mockSelectPrompt *mock.MockSelectPrompt,
	mockDriver *mock.MockDriver,
	driverIndex int,
	host string) {
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
			Return(driverIndex, "test", nil).
			Times(1),

		mockDriver.
			EXPECT().
			SetHost(gomock.Eq(host)).
			Return().
			Times(1),
	)
}

func promptForAuthThenStore(
	mockPrompter *mock.MockFactory,
	mockInputPrompt *mock.MockInputPrompt,
	mockCredentialStorer *mock.MockCredentialsStorer,
	mockDriver *mock.MockDriver,
	domain string) {

	gomock.InOrder(
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

		// store auth from input
		mockCredentialStorer.
			EXPECT().
			Put(domain, credentials.Credential{
				Type:  "token",
				Token: "1234abcd",
			}).
			Return(nil).
			Times(1),

		mockDriver.
			EXPECT().
			Authenticate(gomock.Eq(remote.Auth{
				Token: "1234abcd",
			})).
			Times(1),
	)
}

func getAuthFromStorer(
	mockCredentialStorer *mock.MockCredentialsStorer,
	mockDriver *mock.MockDriver,
	domain string) {

	gomock.InOrder(
		mockDriver.
			EXPECT().
			AuthType().
			Return("token").
			Times(1),

		mockCredentialStorer.
			EXPECT().
			Get(domain).
			Return(credentials.Credential{
				Type:  "token",
				Token: "1234abcd",
			}, nil).
			Times(1),

		mockDriver.
			EXPECT().
			Authenticate(remote.Auth{
				Token: "1234abcd",
			}).
			Return().
			Times(1),
	)
}

func promptForProtocol(mockPrompter *mock.MockFactory, mockSelectPrompt *mock.MockSelectPrompt, pIndex int, protocol string) {
	gomock.InOrder(
		mockPrompter.
			EXPECT().
			CreateSelectPrompt(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockSelectPrompt).
			Times(1),

		mockSelectPrompt.
			EXPECT().
			Run().
			Return(pIndex, protocol, nil).
			Times(1),
	)
}

func TestRemoteCmdSuite(t *testing.T) {
	suite.Run(t, new(RemoteCmdSuite))
}
