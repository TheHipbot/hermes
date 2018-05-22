package prompt

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/TheHipbot/hermes/fs"
	"github.com/manifoldco/promptui"

	"github.com/stretchr/testify/suite"
)

type PromptRepoSuite struct {
	suite.Suite
}

type prompterMock struct {
	mock.Mock
}

func (p *prompterMock) CreateSelectPrompt(label string, items interface{}, tmpls *promptui.SelectTemplates) Prompt {
	args := p.Called(label, items, tmpls)
	return args.Get(0).(Prompt)
}

func (s *PromptRepoSuite) TestPrompterCreateSelectRepo() {
	repos := []fs.Repo{
		fs.Repo{
			Name: "github.com/TheHipbot/hermes",
			Path: "/test-repos/github.com/TheHipbot/hermes",
		},
		fs.Repo{
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

func (s *PromptRepoSuite) TestNewRepoSelectPrompt() {
	prompter := new(prompterMock)
	repos := []fs.Repo{
		fs.Repo{
			Name: "github.com/TheHipbot/hermes",
			Path: "/test-repos/github.com/TheHipbot/hermes",
		},
		fs.Repo{
			Name: "github.com/TheHipbot/hermes",
			Path: "/test-repos/github.com/TheHipbot/dockerfiles",
		},
	}
	prompter.
		On("CreateSelectPrompt", "Select a repo", repos, selectRepoTemplates).
		Return(&promptui.Select{
			Label:     "Select a repo",
			Items:     repos,
			Templates: selectRepoTemplates,
		}).
		Once()

	res := NewRepoSelectPrompt(prompter, repos)
	s.IsType(res, &promptui.Select{}, "Should be a promptui prompt type")
	selectP := res.(*promptui.Select)
	s.Equal(selectP.Label, "Select a repo", "Should return prompt with the correct label")
	s.Equal(selectP.Items, repos, "Should return prompt with the correct items")
	s.Equal(selectP.Templates, selectRepoTemplates, "Should return prompt with the correct templates")
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(PromptRepoSuite))
}
