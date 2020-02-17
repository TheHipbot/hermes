package remote

import (
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

var (
	// URL format for gitlab repo requests
	defaultGitlabAPIHost = "https://gitlab.com"
	gitlabUserRequestFmt = "/api/v4/projects?membership=%t&per_page=20&private_token=%s&page=%d"
)

func gitlabCreator(opts *DriverOpts) (Driver, error) {
	newDriver := &GitLab{}
	if opts.Auth == nil {
		return newDriver, ErrAuth
	}
	switch opts.Auth.Type {
	case "token":
		client := gitlab.NewOAuthClient(nil, opts.Auth.Token)
		newDriver.client = client
	}
	if opts.Host != "" {
		newDriver.client.SetBaseURL(opts.Host)
	}
	newDriver.Host = defaultGitlabAPIHost
	newDriver.Opts = opts
	return newDriver, nil
}

// GitLab is a client to gitlab
type GitLab struct {
	Auth
	Host   string
	Opts   *DriverOpts
	client *gitlab.Client
}

// SetHost sets github driver host to provided string
func (gl *GitLab) SetHost(host string) {
	gl.Host = host
	gl.client.SetBaseURL(host)
}

// Authenticate sets Auth object for driver
func (gl *GitLab) Authenticate(a Auth) {
	gl.Auth = a
}

// AuthType sets Auth object for driver
func (gl *GitLab) AuthType() string {
	return authToken
}

// GetRepos gets the repos for the github user
func (gl *GitLab) GetRepos() ([]map[string]string, error) {
	allRepos := []map[string]string{}
	membership := !gl.Opts.AllRepos
	opts := &gitlab.ListProjectsOptions{
		Membership: &membership,
		ListOptions: gitlab.ListOptions{
			PerPage: 40,
		},
	}

	if gl.Auth.Token == "" && gl.Auth.Username == "" {
		return nil, ErrAuth
	}

	for {
		projects, resp, err := gl.client.Projects.ListProjects(opts)
		if err != nil {
			return allRepos, err
		}
		allRepos, err = mapGitLabProjects(allRepos, projects)
		if err != nil {
			return allRepos, err
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func mapGitLabProjects(acc []map[string]string, projects []*gitlab.Project) ([]map[string]string, error) {
	for _, p := range projects {
		entry := make(map[string]string, 4)
		entry["url"] = p.WebURL
		entry["name"] = strings.Split(entry["url"], "://")[1]
		entry["clone_url"] = p.HTTPURLToRepo
		entry["ssh_url"] = p.SSHURLToRepo
		acc = append(acc, entry)
	}
	return acc, nil
}
