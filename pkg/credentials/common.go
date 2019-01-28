//go:generate mockgen -package mock -destination ../../mock/mock_credentials.go github.com/TheHipbot/hermes/pkg/credentials Storer
package credentials

import "errors"

var (
	// ErrCredentialNotFound is an error that will be returned when
	// the requested credential was not found in the Storer
	ErrCredentialNotFound = errors.New("Credential not found")

	// ErrCredentialStorerError is an error returned when the credential
	// storer throws an error
	ErrCredentialStorerError = errors.New("Credential storer error")
)

// Credential stores a credential and its type
type Credential struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
	Type     string `yaml:"type"`
}

// Storer an interface to store and retrieve Credentials
type Storer interface {
	Get(key string) (Credential, error)
	Put(key string, cred Credential) error
	Close() error
}

func NewMemStorer() Storer {
	return &storage{
		credentials: make(map[string]Credential),
	}
}

type storage struct {
	credentials map[string]Credential
}

func (s *storage) Get(key string) (Credential, error) {
	if s.credentials != nil {
		if cred, ok := s.credentials[key]; !ok {
			return Credential{}, ErrCredentialNotFound
		} else {
			return cred, nil
		}
	}

	return Credential{}, ErrCredentialStorerError
}

func (s *storage) Put(key string, cred Credential) error {
	if s.credentials != nil {
		s.credentials[key] = cred
		return nil
	}

	return ErrCredentialStorerError
}

func (s *storage) Close() error {
	return nil
}
