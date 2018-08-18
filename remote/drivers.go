package remote

import "errors"

// Driver has GetRepos to find repos from remote
type Driver interface {
	GetRepos() ([]map[string]string, error)
	SetHost(host string)
	Authenticate(a Auth)
}

// Auth struct to hold credentials to
// authenticate to remote
type Auth struct {
	Username string
	Password string
	Token    string
}

var (
	creators map[string]func() (Driver, error)
)

func init() {
	creators = map[string]func() (Driver, error){
		"github": githubCreator,
	}
}

// CreateDriver creates a driver
func CreateDriver(t string) (Driver, error) {
	c, ok := creators[t]
	if !ok {
		return nil, errors.New("Driver not implemented")
	}
	return c()
}
