//go:generate mockgen -package mock -destination ../mock/mock_cache.go github.com/TheHipbot/hermes/cache Cache
package cache

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

type cache struct {
	storer  storer
	Version string             `json:"version"`
	Remotes map[string]*Remote `json:"remotes"`
}

// Cache interface to open and save repo cache
type Cache interface {
	Open()
	Save() error
	Close() error
	Add(name, path string) error
	AddRemote(url, name string) error
	Search(needle string) []Repo
}

// Remote is a parent node in the cache tree
type Remote struct {
	Name     string            `json:"name"`
	URL      string            `json:"url"`
	Protocol string            `json:"protocol"`
	Type     string            `json:"type"`
	Meta     map[string]string `json:"meta"`
	Repos    []Repo            `json:"repos"`
}

// Repo stores a repo and its location on the filesystem
// for use in autocomplete
type Repo struct {
	Name string `json:"name"`
	Path string `json:"repo_path"`
}

// storer persists the cache
type storer interface {
	io.ReadWriteSeeker
	io.Closer
	Truncate(size int64) error
}

// NewCache creates a cache then returns it
func NewCache(storer storer) Cache {
	return &cache{
		storer: storer,
	}
}

// Open the cache from the cache.json file in config
// directory
func (c *cache) Open() {
	_, err := c.storer.Seek(0, 0)
	raw, err := ioutil.ReadAll(c.storer)
	var result cache
	if err != nil {
		c.Version = cacheFormatVersion
		c.Remotes = make(map[string]*Remote)
	} else {
		result = cache{}
		if err := json.Unmarshal(raw, &result); err != nil {
			c.Version = cacheFormatVersion
			c.Remotes = make(map[string]*Remote)
		} else {
			c.Version = result.Version
			c.Remotes = result.Remotes
		}
	}
}

// Save cache to ConfigFS
func (c *cache) Save() error {
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = c.storer.Seek(0, 0)
	if err != nil {
		return err
	}
	p, err := c.storer.Write(raw)
	if err != nil {
		return err
	}
	c.storer.Truncate(int64(p))
	if err != nil {
		return err
	}
	_, err = c.storer.Seek(0, 0)
	return err
}

// Close cache storer
func (c *cache) Close() error {
	return c.storer.Close()
}

// Add a repo to the cache
func (c *cache) Add(name, path string) error {
	repoPath := fmt.Sprintf("%s%s", path, name)
	remote := strings.Split(name, "/")[0]

	if r, ok := c.Remotes[remote]; ok {
		c.Remotes[remote].Repos = append(r.Repos, Repo{
			Name: name,
			Path: repoPath,
		})
	} else {
		remoteURL, err := url.Parse(fmt.Sprintf("https://%s", remote))
		if err != nil {
			return err
		}

		c.Remotes[remote] = &Remote{
			Name: remote,
			URL:  remoteURL.String(),
			Repos: []Repo{
				Repo{
					Name: name,
					Path: repoPath,
				},
			},
		}
	}

	return nil
}

// Remove a repo from the cache
func (c *cache) Remove(name string) error {
	found := false
	remote := strings.Split(name, "/")[0]

	if r, ok := c.Remotes[remote]; ok {
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
func (c *cache) AddRemote(url, name string) error {
	// type Remote struct {
	// 	Name     string            `json:"name"`
	// 	URL      string            `json:"url"`
	// 	Protocol string            `json:"protocol"`
	// 	Type     string            `json:"type"`
	// 	Meta     map[string]string `json:"meta"`
	// 	Repos    []Repo            `json:"repos"`
	// }

	if _, ok := c.Remotes[name]; ok {
		return errors.New("Remote already exists")
	}

	remote := &Remote{
		Name:     name,
		URL:      url,
		Protocol: "http",
		Repos:    []Repo{},
	}

	c.Remotes[name] = remote
	return nil
}

// Search will search the cache for any repos that match the
// needle string
func (c *cache) Search(needle string) []Repo {
	lowerSearch := strings.ToLower(needle)
	var results []Repo
	for _, remote := range c.Remotes {
		for _, repo := range remote.Repos {
			if strings.Contains(strings.ToLower(repo.Name), lowerSearch) {
				results = append(results, repo)
			}
		}
	}
	return results
}
