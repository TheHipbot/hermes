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

// Cache holds the cache of remotes and their repos
type Cache struct {
	storer  Storer
	Version string             `json:"version"`
	Remotes map[string]*Remote `json:"remotes"`
}

// Remote is a parent node in the cache tree
type Remote struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Repos []Repo `json:"repos"`
}

// Repo stores a repo and its location on the filesystem
// for use in autocomplete
type Repo struct {
	Name string `json:"name"`
	Path string `json:"repo_path"`
}

// Storer persists the cache
type Storer interface {
	io.ReadWriteSeeker
	io.Closer
	Truncate(size int64) error
}

// NewCache creates a cache then returns it
func NewCache(storer Storer) *Cache {
	return &Cache{
		storer: storer,
	}
}

// Open the cache from the cache.json file in config
// directory
func (c *Cache) Open() {
	_, err := c.storer.Seek(0, 0)
	raw, err := ioutil.ReadAll(c.storer)
	var result Cache
	if err != nil {
		c.Version = cacheFormatVersion
		c.Remotes = make(map[string]*Remote)
	} else {
		result = Cache{}
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
func (c *Cache) Save() error {
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
func (c *Cache) Close() error {
	return c.storer.Close()
}

// Add a repo to the cache
func (c *Cache) Add(name, path string) error {
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
func (c *Cache) Remove(name string) error {
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

// Search will search the cache for any repos that match the
// needle string
func (c *Cache) Search(needle string) []Repo {
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
