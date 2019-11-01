// +build !gogit

package repo

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func init() {
	RegisterCloner("git", func() (Cloner, error) {
		return &CallThroughCloner{}, nil
	})
}

// CallThroughCloner is a cloner which calls through to
// the git binary installed on the system for cloning git
// repositories
type CallThroughCloner struct{}

// Clone clones a repository
func (c *CallThroughCloner) Clone(path string, opts *CloneOptions) error {
	cmd := exec.Command("git", "clone", "--progress", opts.URL, path)
	cmd.Stdout = os.Stdout
	var errBuffer bytes.Buffer
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		raw, err := ioutil.ReadAll(&errBuffer)
		if err != nil {
			return ErrCloneRepo
		}
		errString := string(raw)
		if strings.Contains(errString, "already exists and is not an empty directory") {
			return ErrRepoAlreadyExists
		}
	}

	return nil
}
