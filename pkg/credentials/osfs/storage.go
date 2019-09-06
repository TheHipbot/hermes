package osfs

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/TheHipbot/hermes/pkg/credentials"
)

// Storage stores credentials to
// the storer
type Storage struct {
	storer storer
}

type storer interface {
	io.ReadWriteSeeker
	io.Closer
	Truncate(size int64) error
}

// NewFSStorer takes a storer as an argument, returns
// a *osfs.Storage with the storer
func NewFSStorer(storer storer) *Storage {
	return &Storage{
		storer: storer,
	}
}

func (s *Storage) open() (map[string]credentials.Credential, error) {
	_, err := s.storer.Seek(0, 0)
	raw, err := ioutil.ReadAll(s.storer)

	var result map[string]credentials.Credential
	if err != nil {
		return result, err
	}

	if err := yaml.Unmarshal(raw, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (s *Storage) save(creds map[string]credentials.Credential) error {
	raw, err := yaml.Marshal(creds)
	if err != nil {
		return err
	}

	_, err = s.storer.Seek(0, 0)
	if err != nil {
		return err
	}
	p, err := s.storer.Write(raw)
	if err != nil {
		return err
	}
	s.storer.Truncate(int64(p))
	if err != nil {
		return err
	}
	_, err = s.storer.Seek(0, 0)
	return err
}

// Get returns the credential of the given key from the
// the credential file
func (s *Storage) Get(key string) (credentials.Credential, error) {
	creds, err := s.open()
	if err != nil {
		return credentials.Credential{}, err
	} else if creds == nil {
		creds = make(map[string]credentials.Credential)
	}

	cred, ok := creds[key]
	if !ok {
		return credentials.Credential{}, credentials.ErrCredentialNotFound
	}
	return cred, nil
}

// Put stores the given credential on the given key and writes the
// credential to the file
func (s *Storage) Put(key string, cred credentials.Credential) error {
	creds, err := s.open()
	if err != nil {
		return err
	} else if creds == nil {
		creds = make(map[string]credentials.Credential)
	}

	creds[key] = cred
	if err := s.save(creds); err != nil {
		return err
	}
	return nil
}

// Delete removes the stored credential with the given key
func (s *Storage) Delete(key string) error {
	return nil
}

// Close will close file used to store credentials
func (s *Storage) Close() error {
	return s.storer.Close()
}
