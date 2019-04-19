package repo

import (
	"errors"
	"os"
	"regexp"

	"github.com/kevinburke/ssh_config"
	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	sshgit "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var (
	appFs         billy.Filesystem
	defaultCloner Cloner
)

func init() {
	appFs = osfs.New("")
	defaultCloner = &GitCloner{
		Fs: appFs,
	}
}

// GitRepository holds info for git repos and
// implements the Repository interface
type GitRepository struct {
	Fs       billy.Filesystem
	Name     string
	URL      string
	Protocol string
	Cloner   Cloner
}

// NewGitRepository creates a GitRepository
func NewGitRepository(name, url string) *GitRepository {
	return &GitRepository{
		Fs:     appFs,
		Name:   name,
		URL:    url,
		Cloner: defaultCloner,
	}
}

// Clone git repository to path
func (gr *GitRepository) Clone(path string) error {

	opts := &CloneOptions{
		URL: gr.URL,
	}

	switch gr.Protocol {
	case "ssh":
		var hostname string
		r, err := regexp.Compile(`.+@([a-zA-z.\-0-9]+)[:/].+`)
		if err != nil {
			return errors.New("Regular expression should compile")
		}
		cps := r.FindStringSubmatch(gr.URL)
		if len(cps) < 2 {
			return errors.New("Could not find hostname in URL")
		}
		hostname = cps[1]
		a, err := getSSHAuth(hostname)
		opts.Auth = a
	}

	err := gr.Cloner.Clone(path, opts)
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
			pk, err := sshgit.NewPublicKeysFromFile("git", keyPath, "")
			if err != nil {
				return &sshgit.PublicKeys{}, err
			}
			pk.HostKeyCallbackHelper = sshgit.HostKeyCallbackHelper{
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
			return pk, nil
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
