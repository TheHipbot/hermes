//go:generate mockgen -package mock -destination ../mock/mock_driver.go github.com/TheHipbot/hermes/remote Driver
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

var (
	creators map[string]func() (Driver, error)

	// auth types
	authToken = "token"

	// ErrNotImplemented error when driver requested is not implemented
	ErrNotImplemented = errors.New("Driver not implemented")
	// ErrInvalidOpts error when options for driver are invald
	ErrInvalidOpts = errors.New("Invalid driver options")
)

func init() {
	creators = map[string]func() (Driver, error){
		"github": githubCreator,
	}
}

// NewDriver creates a driver
func NewDriver(t string) (Driver, error) {
	c, ok := creators[t]
	if !ok {
		return nil, ErrNotImplemented
	}
	return c()
}

// RegisterDriver registers a new driver from the given name and driver returning
// function
func RegisterDriver(name string, creator func() (Driver, error)) {
	creators[name] = creator
}
