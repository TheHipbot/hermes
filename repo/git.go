package repo

import (
	"net/url"
	"os"

	"gopkg.in/src-d/go-billy.v4/osfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// GitRepository holds info for git repos and
// implements the Repository interface
type GitRepository struct {
	Name string
	URL  *url.URL
}

// Clone git repository to path
func (gr *GitRepository) Clone(path string) error {
	workDir := osfs.New(path)
	dot, _ := workDir.Chroot(".git")
	storer, err := filesystem.NewStorage(dot)
	if err != nil {
		os.Exit(1)
	}

	if _, err := git.Clone(storer, workDir, &git.CloneOptions{
		URL:      gr.URL.String(),
		Progress: os.Stdout,
	}); err != nil && err != git.ErrRepositoryAlreadyExists {
		return err
	}
	return nil
}
