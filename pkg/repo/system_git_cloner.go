// +build !gogit

package repo

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/sync/errgroup"
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

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	var g errgroup.Group
	g.Go(func() error {
		return filterGitErrors(os.Stdout, stderrPipe)
	})

	if err := cmd.Start(); err != nil {
		return err
	}

	cmd.Wait()

	return g.Wait()
}

type gitOuput struct {
	gitReader io.Reader
	out       io.Writer
	err       error
}

func filterGitErrors(w io.Writer, r io.Reader) error {
	errRegex := "^fatal: (.+)$"
	re, err := regexp.Compile(errRegex)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var finalErr error
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			message := ""
			if len(matches) > 1 {
				message = matches[1]
			}
			if strings.Contains(message, "already exists and is not an empty directory") {
				finalErr = ErrRepoAlreadyExists
			} else {
				finalErr = ErrCloneRepo
			}
		} else {
			w.Write([]byte(line + "\n"))
		}
	}
	return finalErr
}
