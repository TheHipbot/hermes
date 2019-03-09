//go:generate mockgen -package mock -destination ../../mock/mock_storage.go github.com/TheHipbot/hermes/pkg/storage Storage

package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
)

const (
	cacheFormatVersion = "0.0.1"
)

// storer persists the cache
type storer interface {
	io.ReadWriteSeeker
	io.Closer
	Truncate(size int64) error
}

// Storage interface to open and save repo Storage
type Storage interface {
	Open()
	Save() error
	Close() error
	AddRepository(repo *Repository) error
	RemoveRepository(name string) error
	AddRemote(url, name, protocol string) error
	SearchRepositories(needle string) []Repository
	SearchRemote(remote string) (*Remote, bool)
}

type storage struct {
	storer  storer
	Version string             `json:"version"`
	Remotes map[string]*Remote `json:"remotes"`
}

// NewStorage creates a cache then returns it
func NewStorage(storer storer) Storage {
	return &storage{
		storer: storer,
	}
}

// Open the cache from the provided storer
func (s *storage) Open() {
	_, err := s.storer.Seek(0, 0)
	raw, err := ioutil.ReadAll(s.storer)
	var result storage
	if err != nil {
		s.Version = cacheFormatVersion
		s.Remotes = make(map[string]*Remote)
	} else {
		result = storage{}
		if err := json.Unmarshal(raw, &result); err != nil {
			s.Version = cacheFormatVersion
			s.Remotes = make(map[string]*Remote)
		} else {
			s.Version = result.Version
			s.Remotes = result.Remotes
		}
	}
}

// Save cache to ConfigFS
func (s *storage) Save() error {
	raw, err := json.Marshal(s)
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

// Close cache storer
func (s *storage) Close() error {
	return s.storer.Close()
}

// AddRepository adds a repo to the cache
func (s *storage) AddRepository(repo *Repository) error {
	remote := strings.Split(repo.Name, "/")[0]

	if r, ok := s.Remotes[remote]; ok {
		if _, ok := r.Repos[repo.Name]; ok {
			return errors.New("Repo already exists")
		}
		r.Repos[repo.Name] = repo
	} else {
		remoteURL, err := url.Parse(fmt.Sprintf("https://%s", remote))
		if err != nil {
			return err
		}

		s.Remotes[remote] = &Remote{
			Name: remote,
			URL:  remoteURL.String(),
			Repos: map[string]*Repository{
				repo.Name: repo,
			},
		}
	}

	return nil
}

// RemoveRepository a repo from the cache
func (s *storage) RemoveRepository(name string) error {
	remote := strings.Split(name, "/")[0]

	if r, ok := s.Remotes[remote]; ok {
		if _, ok := r.Repos[name]; ok {
			delete(r.Repos, name)
			return nil
		}
	}

	return errors.New("Repo not found")
}

// AddRemote adds a remote to the cache
func (s *storage) AddRemote(url, name, protocol string) error {
	if _, ok := s.Remotes[name]; ok {
		return errors.New("Remote already exists")
	}

	remote := &Remote{
		Name:     name,
		URL:      url,
		Protocol: protocol,
		Repos:    map[string]*Repository{},
	}

	s.Remotes[name] = remote
	return nil
}

// Search will search the cache for any repos that match the
// needle string
func (s *storage) SearchRepositories(needle string) []Repository {
	lowerSearch := strings.ToLower(needle)
	var results []Repository
	for _, remote := range s.Remotes {
		for name, repo := range remote.Repos {
			if strings.Contains(strings.ToLower(name), lowerSearch) {
				results = append(results, *repo)
			}
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	return results
}

// SearchRemote will return the Remote and true if present or an
// empty remote and false if not
func (s *storage) SearchRemote(remote string) (*Remote, bool) {
	ptr, ok := s.Remotes[remote]
	if ok {
		return ptr, ok
	}
	return &Remote{}, ok
}
