// +build gogit

package repo

import (
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"os"

	billy "gopkg.in/src-d/go-billy.v4"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
)

func init() {
	RegisterCloner("git", func() (Cloner, error) {
		return &GitCloner{
			Fs: appFs,
		}, nil
	})
}

// GitCloner is a Cloner which uses go-git to clone
// git repositories
type GitCloner struct {
	Fs billy.Filesystem
}

// Clone clones a repository
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
