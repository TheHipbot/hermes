package repo

import (
	"errors"
	"os"

	"github.com/kevinburke/ssh_config"
	homedir "github.com/mitchellh/go-homedir"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

var (
	appFs         billy.Filesystem
	defaultCloner cloner
)

func init() {
	appFs = osfs.New("")
	defaultCloner = &gitCloner{}
}

// GitRepository holds info for git repos and
// implements the Repository interface
type GitRepository struct {
	Fs       billy.Filesystem
	Name     string
	URL      string
	Protocol string
	cloner   cloner
}

// NewGitRepository creates a GitRepository
func NewGitRepository(name, url string) *GitRepository {
	return &GitRepository{
		Fs:     appFs,
		Name:   name,
		URL:    url,
		cloner: defaultCloner,
	}
}

// Clone git repository to path
func (gr *GitRepository) Clone(path string) error {
	repoFs, _ := gr.Fs.Chroot(path)
	dot, _ := repoFs.Chroot(".git")
	storer := filesystem.NewStorage(dot, cache.NewObjectLRU(cache.DefaultMaxSize))

	opts := &git.CloneOptions{
		URL:      gr.URL,
		Progress: os.Stdout,
	}

	switch gr.Protocol {
	case "ssh":
		a, _ := getSSHAuth("github.com")
		opts.Auth = a
	}

	err := gr.cloner.clone(storer, repoFs, opts)
	return err
}

func getSSHAuth(host string) (transport.AuthMethod, error) {
	pathsToCheck := []string{
		ssh_config.Get(host, "IdentityFile"),
		"~/.ssh/id_rsa",
		"/etc/ssh/id_rsa",
	}

	for _, path := range pathsToCheck {
		if keyPath, err := pathIfExists(path); err == nil {
			return ssh.NewPublicKeysFromFile("git", keyPath, "")
		}
	}

	return nil, errors.New("No ssh key found")
}

func pathIfExists(path string) (string, error) {
	keyPath, err := homedir.Expand(path)
	if err != nil {
		return "", err
	}
	if _, err := appFs.Stat(keyPath); os.IsNotExist(err) {
		return "", err
	}
	return keyPath, nil
}
