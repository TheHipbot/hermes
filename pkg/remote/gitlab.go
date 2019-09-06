package remote

import (
	"fmt"
	"strings"
)

var (
	// URL format for gitlab repo requests
	defaultGitlabAPIHost = "https://gitlab.com"
	gitlabUserRequestFmt = "/api/v4/projects?membership=%t&per_page=20&private_token=%s&page=%d"
)

func gitlabCreator(opts *DriverOpts) (Driver, error) {
	return &Gitlab{
		Host: defaultGitlabAPIHost,
		Opts: opts,
	}, nil
}

// Gitlab is a client to gitlab
type Gitlab struct {
	Auth
	Host string
	Opts *DriverOpts
}

// SetHost sets github driver host to provided string
func (gl *Gitlab) SetHost(host string) {
	gl.Host = host
}

// Authenticate sets Auth object for driver
func (gl *Gitlab) Authenticate(a Auth) {
	gl.Auth = a
}

// AuthType sets Auth object for driver
func (gl *Gitlab) AuthType() string {
	return authToken
}

// GetRepos gets the repos for the github user
func (gl *Gitlab) GetRepos() ([]map[string]string, error) {
	urlFormat := fmt.Sprintf("%s%s", gl.Host, gitlabUserRequestFmt)
	if gl.Auth.Token == "" && gl.Auth.Username == "" {
		return nil, ErrAuth
	}

	page := 1
	accumulator := []map[string]string{}
	return getRepoHelper(fmt.Sprintf(urlFormat, !gl.Opts.AllRepos, gl.Auth.Token, page), accumulator, func(item map[string]interface{}) map[string]string {
		entry := make(map[string]string, 3)
		url := item["web_url"].(string)
		entry["url"] = url
		entry["name"] = strings.Split(url, "://")[1]
		entry["clone_url"] = item["http_url_to_repo"].(string)
		entry["ssh_url"] = item["ssh_url_to_repo"].(string)
		return entry
	})
}
