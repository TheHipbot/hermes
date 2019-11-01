//go:generate mockgen -package mock -destination ../../mock/mock_cloner.go github.com/TheHipbot/hermes/pkg/repo Cloner

package repo

import (
	"errors"
	"fmt"
)

var (
	creators = map[string]func() (Cloner, error){}

	// ErrRepoAlreadyExists is an error returned when the repo already exists
	ErrRepoAlreadyExists = errors.New("repository already exists")
	// ErrCloneRepo when there is a normal error cloning repo
	ErrCloneRepo = errors.New("error cloning repo")
)

// Repository struct holds information for a repository
type Repository interface {
	Clone(path string, opts *CloneOptions) error
}

// AuthMethod is the method to authenticate
// for cloners
type AuthMethod interface {
	Name() string
	fmt.Stringer
}

// CloneOptions is for packaging various
// options for cloning repositories
type CloneOptions struct {
	URL  string
	Auth AuthMethod
}

// Cloner is an interface for cloning repositories
type Cloner interface {
	Clone(path string, opts *CloneOptions) error
}

// RegisterCloner takes a name for the cloner type and a function
// which creates an instance of that cloner
func RegisterCloner(name string, creator func() (Cloner, error)) {
	creators[name] = creator
}

// NewCloner return a cloner from the given type and error
func NewCloner(name string) (Cloner, error) {
	if c, ok := creators[name]; ok {
		return c()
	}
	return nil, errors.New("Ivalid Cloner Type")
}
