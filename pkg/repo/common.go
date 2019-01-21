package repo

import (
	billy "gopkg.in/src-d/go-billy.v4"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// Repository struct holds information for a repository
type Repository interface {
	Clone(path string) error
}

type cloner interface {
	clone(storer *filesystem.Storage, tree billy.Filesystem, opts *git.CloneOptions) error
}

type gitCloner struct{}

func (gc *gitCloner) clone(storer *filesystem.Storage, tree billy.Filesystem, opts *git.CloneOptions) error {
	_, err := git.Clone(storer, tree, opts)
	return err
}
