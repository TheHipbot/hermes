package remote

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

var (
	// URL format for github repo requests
	defaultGitHubAPIHost = "https://api.github.com"
	gitHubUserRequestFmt = "/user/repos?access_token=%s&page=%d"
)

func githubCreator(opts *DriverOpts) (Driver, error) {
	newDriver := &GitHub{}
	if opts.Auth == nil {
		return newDriver, ErrAuth
	}
	switch opts.Auth.Type {
	case "token":
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: opts.Auth.Token,
		})
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
		newDriver.client = client
	}
	newDriver.Host = defaultGitHubAPIHost
	newDriver.Opts = opts
	return newDriver, nil
}

// GitHub is a client to github
type GitHub struct {
	Auth
	Host   string
	Opts   *DriverOpts
	client *github.Client
}

// SetHost sets github driver host to provided string
func (gh *GitHub) SetHost(host string) {
	match, err := regexp.MatchString("^(https?://)?github.com", host)
	if match || err != nil {
		gh.Host = defaultGitHubAPIHost
	} else {
		gh.Host = host
	}
}

// Authenticate sets Auth object for driver
func (gh *GitHub) Authenticate(a Auth) {
	gh.Auth = a
}

// AuthType sets Auth object for driver
func (gh *GitHub) AuthType() string {
	return authToken
}

// GetRepos gets the repos for the github user
func (gh *GitHub) GetRepos() ([]map[string]string, error) {
	allRepos := []map[string]string{}
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 40,
		},
	}

	for {
		repos, resp, err := gh.client.Repositories.List(context.Background(), "", opts)
		if err != nil {
			return allRepos, err
		}
		allRepos, err = mapGitHubRepos(allRepos, repos)
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

func mapGitHubRepos(acc []map[string]string, repos []*github.Repository) ([]map[string]string, error) {
	for _, r := range repos {
		entry := make(map[string]string, 4)
		entry["url"] = r.GetHTMLURL()
		entry["name"] = strings.Split(entry["url"], "://")[1]
		entry["clone_url"] = r.GetCloneURL()
		entry["ssh_url"] = r.GetSSHURL()
		acc = append(acc, entry)
	}
	return acc, nil
}
