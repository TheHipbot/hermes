package remote

import "errors"

var (
	// URL format for gitlab repo requests
	defaultGitlabAPIHost = "https://gitlab.com"
	gitlabUserRequestFmt = "/api/v4/projects?membership=true&per_page=100&private_token=%s&page=%d"
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
	// urlFormat := fmt.Sprintf("%s%s", gl.Host, gitlabUserRequestFmt)
	if gl.Auth.Token == "" && gl.Auth.Username == "" {
		return nil, errors.New("Auth is empty")
	}

	// page := 1
	accumulator := []map[string]string{}
	// return getRepoHelper(fmt.Sprintf(urlFormat, gl.Auth.Token, page), accumulator)
	return accumulator, nil
}
