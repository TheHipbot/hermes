package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
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

type Storage struct {
	storer  storer
	Version string             `json:"version"`
	Remotes map[string]*Remote `json:"remotes"`
}

// NewStorage creates a cache then returns it
func NewStorage(storer storer) *Storage {
	return &Storage{
		storer: storer,
	}
}

// Open the cache from the cache.json file in config
// directory
func (s *Storage) Open() {
	_, err := s.storer.Seek(0, 0)
	raw, err := ioutil.ReadAll(s.storer)
	var result Storage
	if err != nil {
		s.Version = cacheFormatVersion
		s.Remotes = make(map[string]*Remote)
	} else {
		result = Storage{}
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
func (s *Storage) Save() error {
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
func (s *Storage) Close() error {
	return s.storer.Close()
}

// AddRepo a repo to the cache
func (s *Storage) AddRepo(name, path string) error {
	repoPath := fmt.Sprintf("%s%s", path, name)
	remote := strings.Split(name, "/")[0]

	if r, ok := s.Remotes[remote]; ok {
		s.Remotes[remote].Repos = append(r.Repos, Repository{
			Name: name,
			Path: repoPath,
		})
	} else {
		remoteURL, err := url.Parse(fmt.Sprintf("https://%s", remote))
		if err != nil {
			return err
		}

		s.Remotes[remote] = &Remote{
			Name: remote,
			URL:  remoteURL.String(),
			Repos: []Repository{
				Repository{
					Name: name,
					Path: repoPath,
				},
			},
		}
	}

	return nil
}

// RemoveRepo a repo from the cache
func (s *Storage) RemoveRepo(name string) error {
	found := false
	remote := strings.Split(name, "/")[0]

	if r, ok := s.Remotes[remote]; ok {
		for i, repo := range r.Repos {
			if strings.Compare(repo.Name, name) == 0 {
				r.Repos = append(r.Repos[:i], r.Repos[i+1:]...)
				found = true
				break
			}
		}
	}

	if !found {
		return errors.New("Repo not found")
	}

	return nil
}

// AddRemote adds a remote to the cache
func (s *Storage) AddRemote(url, name string) error {
	// type Remote struct {
	// 	Name     string            `json:"name"`
	// 	URL      string            `json:"url"`
	// 	Protocol string            `json:"protocol"`
	// 	Type     string            `json:"type"`
	// 	Meta     map[string]string `json:"meta"`
	// 	Repos    []Repo            `json:"repos"`
	// }

	if _, ok := s.Remotes[name]; ok {
		return errors.New("Remote already exists")
	}

	remote := &Remote{
		Name:     name,
		URL:      url,
		Protocol: "http",
		Repos:    []Repository{},
	}

	s.Remotes[name] = remote
	return nil
}

// Search will search the cache for any repos that match the
// needle string
func (s *Storage) Search(needle string) []Repository {
	lowerSearch := strings.ToLower(needle)
	var results []Repository
	for _, remote := range s.Remotes {
		for _, repo := range remote.Repos {
			if strings.Contains(strings.ToLower(repo.Name), lowerSearch) {
				results = append(results, repo)
			}
		}
	}
	return results
}
