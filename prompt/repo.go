//go:generate mockgen -destination ../mock/mock_prompt.go github.com/TheHipbot/hermes/prompt Prompt,Factory

package prompt

import (
	"github.com/TheHipbot/hermes/cache"
	"github.com/manifoldco/promptui"
)

var (
	selectRepoTemplates = &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U00002714 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | white }}",
		Selected: "{{ .Name | green }}",
	}
)

// Prompt is a user prompt which can be Run
type Prompt interface {
	Run() (int, string, error)
}

// Factory makes prompts
type Factory interface {
	CreateSelectPrompt(label string, items interface{}, tmpls *promptui.SelectTemplates) Prompt
}

// Prompter is an implementation of Factory which creates prompts
// from the github.com/manifoldco/promptui library
type Prompter struct{}

// CreateSelectPrompt creates a select prompt
func (b *Prompter) CreateSelectPrompt(label string, items interface{}, tmpls *promptui.SelectTemplates) Prompt {
	return &promptui.Select{
		Label:     label,
		Items:     items,
		Templates: tmpls,
	}
}

// NewRepoSelectPrompt returns a Prompt to select a Repo cache entry from
// a given list and return the selected repo
func NewRepoSelectPrompt(f Factory, repos []cache.Repo) Prompt {
	return f.CreateSelectPrompt("Select a repo", repos, selectRepoTemplates)
}
