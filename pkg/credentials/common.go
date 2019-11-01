//go:generate mockgen -package mock -destination ../../mock/mock_credentials.go -mock_names Storer=MockCredentialsStorer github.com/TheHipbot/hermes/pkg/credentials Storer

package credentials

import "errors"

var (
	// ErrCredentialNotFound is an error that will be returned when
	// the requested credential was not found in the Storer
	ErrCredentialNotFound = errors.New("credential not found")

	// ErrCredentialStorerError is an error returned when the credential
	// storer throws an error
	ErrCredentialStorerError = errors.New("credential storer error")
)

// Credential stores a credential and its type
type Credential struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token,omitempty"`
	Type     string `yaml:"type"`
}

// Storer an interface to store and retrieve Credentials
type Storer interface {
	Get(key string) (Credential, error)
	Put(key string, cred Credential) error
	Delete(key string) error
	Close() error
}

// NewMemStorer returns an in memory implementation of the
// Storer interface
func NewMemStorer() Storer {
	return &storage{
		credentials: make(map[string]Credential),
	}
}

type storage struct {
	credentials map[string]Credential
}

// Get retrieves the Credential of the given key or
// returns an error if unable to do so
func (s *storage) Get(key string) (Credential, error) {
	if s.credentials != nil {
		cred, ok := s.credentials[key]
		if !ok {
			return Credential{}, ErrCredentialNotFound
		}
		return cred, nil
	}

	return Credential{}, ErrCredentialStorerError
}

// Put stores the credential of the given key or returns an
// error if unable to do so
func (s *storage) Put(key string, cred Credential) error {
	if s.credentials != nil {
		s.credentials[key] = cred
		return nil
	}

	return ErrCredentialStorerError
}

func (s *storage) Delete(key string) error {
	return nil
}

// Close for the in memory is no-op
// no closing necessary
func (s *storage) Close() error {
	return nil
}
