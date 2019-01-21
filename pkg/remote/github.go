package remote

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// URL format for github repo requests
	defaultGitHubAPIHost = "https://api.github.com"
	gitHubUserRequestFmt = "/user/repos?access_token=%s&page=%d"
)

func githubCreator(opts *DriverOpts) (Driver, error) {
	return &GitHub{
		Host: defaultGitHubAPIHost,
		Opts: opts,
	}, nil
}

// GitHub is a client to github
type GitHub struct {
	Auth
	Host string
	Opts *DriverOpts
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
	urlFormat := fmt.Sprintf("%s%s", gh.Host, gitHubUserRequestFmt)
	if gh.Auth.Token == "" && gh.Auth.Username == "" {
		return nil, errors.New("Auth is empty")
	}

	page := 1
	accumulator := []map[string]string{}
	return getRepoHelper(fmt.Sprintf(urlFormat, gh.Auth.Token, page), accumulator, func(item map[string]interface{}) map[string]string {
		entry := make(map[string]string, 3)
		url := item["html_url"].(string)
		entry["url"] = url
		entry["name"] = strings.Split(url, "://")[1]
		return entry
	})
}
