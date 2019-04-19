//go:generate mockgen -package mock -destination ../../mock/mock_cloner.go github.com/TheHipbot/hermes/pkg/repo Cloner
package repo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	billy "gopkg.in/src-d/go-billy.v4"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

var (
	// ErrRepoAlreadyExists is an error returned when the repo already exists
	ErrRepoAlreadyExists = errors.New("Repository already exists")
)

// Repository struct holds information for a repository
type Repository interface {
	Clone(path string, opts *CloneOptions) error
}

type AuthMethod interface {
	Name() string
	fmt.Stringer
}

type CloneOptions struct {
	URL  string
	Auth AuthMethod
}

type Cloner interface {
	Clone(path string, opts *CloneOptions) error
}

type GitCloner struct {
	Fs billy.Filesystem
}

func (gc *GitCloner) Clone(path string, opts *CloneOptions) error {
	repoFs, _ := gc.Fs.Chroot(path)
	dot, _ := repoFs.Chroot(".git")
	storer := filesystem.NewStorage(dot, cache.NewObjectLRU(cache.DefaultMaxSize))

	_, err := git.Clone(storer, repoFs, &git.CloneOptions{
		URL:      opts.URL,
		Progress: os.Stdout,
		Auth:     opts.Auth,
	})
	return err
}

type CallThroughCloner struct{}

func (c *CallThroughCloner) Clone(path string, opts *CloneOptions) error {
	cmd := exec.Command("git", "clone", opts.URL, path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stdin
	return nil
}
