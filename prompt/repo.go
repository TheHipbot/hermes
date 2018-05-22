package prompt

import (
	"github.com/TheHipbot/hermes/fs"
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
func NewRepoSelectPrompt(f Factory, repos []fs.Repo) Prompt {
	return f.CreateSelectPrompt("Select a repo", repos, selectRepoTemplates)
}
