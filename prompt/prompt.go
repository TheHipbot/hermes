//go:generate mockgen -package mock -destination ../mock/mock_prompt.go github.com/TheHipbot/hermes/prompt InputPrompt,SelectPrompt,Factory

package prompt

import (
	"github.com/manifoldco/promptui"
)

var (
	selectRepoTemplates = &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U00002714 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | white }}",
		Selected: "{{ .Name | green }}",
	}
	selectRepoLabel   = "Select a repo"
	selectDriverLabel = "Select remote server type"
	inputKeyLabel     = "Enter auth token"
)

// SelectPrompt is a user prompt which can be Run
type SelectPrompt interface {
	Run() (int, string, error)
}

// InputPrompt is a user prompt which can be Run
type InputPrompt interface {
	Run() (string, error)
}

// Factory makes prompts
type Factory interface {
	CreateSelectPrompt(label string, items interface{}, tmpls *promptui.SelectTemplates) SelectPrompt
	CreateInputPrompt(label string) InputPrompt
}

// Prompter is an implementation of Factory which creates prompts
// from the github.com/manifoldco/promptui library
type Prompter struct{}

// Repo stores a repo and its location
type Repo struct {
	Name string `json:"name"`
	Path string `json:"repo_path"`
}

// CreateSelectPrompt creates a select prompt
func (b *Prompter) CreateSelectPrompt(label string, items interface{}, tmpls *promptui.SelectTemplates) SelectPrompt {
	return &promptui.Select{
		Label:     label,
		Items:     items,
		Templates: tmpls,
	}
}

// CreateInputPrompt creates a input prompt
func (b *Prompter) CreateInputPrompt(label string) InputPrompt {
	return &promptui.Prompt{
		Label: label,
	}
}

// CreateRepoSelectPrompt returns a Prompt to select a Repo cache entry from
// a given list and return the selected repo
func CreateRepoSelectPrompt(f Factory, repos interface{}) SelectPrompt {
	return f.CreateSelectPrompt(selectRepoLabel, repos, selectRepoTemplates)
}

//CreateDriverSelectPrompt returns prompt for driver
func CreateDriverSelectPrompt(f Factory, drivers interface{}) SelectPrompt {
	return f.CreateSelectPrompt(selectDriverLabel, drivers, selectRepoTemplates)
}

//CreateTokenInputPrompt returns prompt for auth key
func CreateTokenInputPrompt(f Factory) InputPrompt {
	return f.CreateInputPrompt(inputKeyLabel)
}
