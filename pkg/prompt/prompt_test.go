package prompt

import (
	"testing"

	"github.com/TheHipbot/hermes/cache"
	"github.com/stretchr/testify/mock"

	"github.com/manifoldco/promptui"

	"github.com/stretchr/testify/suite"
)

type PromptRepoSuite struct {
	suite.Suite
}

type prompterMock struct {
	mock.Mock
}

func (p *prompterMock) CreateSelectPrompt(label string, items interface{}, tmpls *promptui.SelectTemplates) SelectPrompt {
	args := p.Called(label, items, tmpls)
	return args.Get(0).(SelectPrompt)
}

func (p *prompterMock) CreateInputPrompt(label string) InputPrompt {
	args := p.Called(label)
	return args.Get(0).(InputPrompt)
}

func (s *PromptRepoSuite) TestPrompterCreateSelectRepo() {
	repos := []cache.Repo{
		cache.Repo{
			Name: "github.com/TheHipbot/hermes",
			Path: "/test-repos/github.com/TheHipbot/hermes",
		},
		cache.Repo{
			Name: "github.com/TheHipbot/hermes",
			Path: "/test-repos/github.com/TheHipbot/dockerfiles",
		},
	}
	p := &Prompter{}
	s.Equal(p.CreateSelectPrompt("blah", repos, selectRepoTemplates), &promptui.Select{
		Label:     "blah",
		Items:     repos,
		Templates: selectRepoTemplates,
	})
}

func (s *PromptRepoSuite) TestCreateRepoSelectPrompt() {
	prompter := new(prompterMock)
	repos := []cache.Repo{
		cache.Repo{
			Name: "github.com/TheHipbot/hermes",
			Path: "/test-repos/github.com/TheHipbot/hermes",
		},
		cache.Repo{
			Name: "github.com/TheHipbot/hermes",
			Path: "/test-repos/github.com/TheHipbot/dockerfiles",
		},
	}
	prompter.
		On("CreateSelectPrompt", selectRepoLabel, repos, selectRepoTemplates).
		Return(&promptui.Select{
			Label:     selectRepoLabel,
			Items:     repos,
			Templates: selectRepoTemplates,
		}).
		Once()

	res := CreateRepoSelectPrompt(prompter, repos)
	s.IsType(res, &promptui.Select{}, "Should be a promptui prompt type")
	selectP := res.(*promptui.Select)
	s.Equal(selectP.Label, selectRepoLabel, "Should return prompt with the correct label")
	s.Equal(selectP.Items, repos, "Should return prompt with the correct items")
	s.Equal(selectP.Templates, selectRepoTemplates, "Should return prompt with the correct templates")
}

func (s *PromptRepoSuite) TestCreateDriverSelectPrompt() {
	prompter := new(prompterMock)
	types := []string{
		"github",
		"gitlab",
		"bitbucket",
	}
	prompter.
		On("CreateSelectPrompt", "Select remote server type", types, selectRepoTemplates).
		Return(&promptui.Select{
			Label:     "Select a repo",
			Items:     types,
			Templates: selectRepoTemplates,
		}).
		Once()

	res := CreateDriverSelectPrompt(prompter, types)
	s.IsType(res, &promptui.Select{}, "Should be a promptui prompt type")
	selectP := res.(*promptui.Select)
	s.Equal(selectP.Label, "Select a repo", "Should return prompt with the correct label")
	s.Equal(selectP.Items, types, "Should return prompt with the correct items")
	s.Equal(selectP.Templates, selectRepoTemplates, "Should return prompt with the correct templates")
}

func (s *PromptRepoSuite) TestCreateTokenInputPrompt() {
	prompter := new(prompterMock)
	prompter.
		On("CreateInputPrompt", "Enter auth token").
		Return(&promptui.Prompt{
			Label: "Enter auth token",
		}).
		Once()

	res := CreateTokenInputPrompt(prompter)
	s.IsType(res, &promptui.Prompt{}, "Should be a promptui prompt type")
	selectP := res.(*promptui.Prompt)
	s.Equal(selectP.Label, "Enter auth token", "Should return prompt with the correct label")
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(PromptRepoSuite))
}
