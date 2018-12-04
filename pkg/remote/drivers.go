//go:generate mockgen -package mock -destination ../../mock/mock_driver.go github.com/TheHipbot/hermes/pkg/remote Driver
package remote

import "errors"

// Driver has GetRepos to find repos from remote
type Driver interface {
	GetRepos() ([]map[string]string, error)
	SetHost(host string)
	Authenticate(a Auth)
	AuthType() string
}

// Auth struct to hold credentials to
// authenticate to remote
type Auth struct {
	Username string
	Password string
	Token    string
}

// DriverOpts provide options for which repos to get
type DriverOpts struct {
	AllRepos bool
	Starred  bool
}

var (
	creators = map[string]func(opts *DriverOpts) (Driver, error){}

	// auth types
	authToken = "token"

	// ErrNotImplemented error when driver requested is not implemented
	ErrNotImplemented = errors.New("Driver not implemented")
	// ErrInvalidOpts error when options for driver are invald
	ErrInvalidOpts = errors.New("Invalid driver options")
)

func init() {
	RegisterDriver("github", githubCreator)
	RegisterDriver("gitlab", gitlabCreator)
}

// NewDriver creates a driver
func NewDriver(t string, opts *DriverOpts) (Driver, error) {
	c, ok := creators[t]
	if !ok {
		return nil, ErrNotImplemented
	}
	return c(opts)
}

// RegisterDriver registers a new driver from the given name and driver returning
// function
func RegisterDriver(name string, creator func(opts *DriverOpts) (Driver, error)) {
	creators[name] = creator
}
