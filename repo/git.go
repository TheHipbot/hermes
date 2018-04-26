package repo

import (
	"net/url"
	"os"

	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

var (
	appFs billy.Filesystem
)

func init() {
	appFs = osfs.New("")
}

// GitRepository holds info for git repos and
// implements the Repository interface
type GitRepository struct {
	Name string
	URL  *url.URL
}

// Clone git repository to path
func (gr *GitRepository) Clone(path string) error {
	repoFs, _ := appFs.Chroot(path)
	dot, _ := repoFs.Chroot(".git")
	storer, err := filesystem.NewStorage(dot)
	if err != nil {
		os.Exit(1)
	}

	// TODO: catch repo already exists error in commands?
	if _, err := git.Clone(storer, repoFs, &git.CloneOptions{
		URL:      gr.URL.String(),
		Progress: os.Stdout,
	}); err != nil && err != git.ErrRepositoryAlreadyExists {
		return err
	}
	return nil
}
