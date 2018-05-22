package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

const (
	cacheFormatVersion = "0.0.1"
)

var (
	configFS *ConfigFS
	once     sync.Once
)

func init() {
	configFS = NewConfigFS()
}

// Cache holds the cache of remotes and their repos
type Cache struct {
	cfs     *ConfigFS
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

// NewCache creates a cache then returns it
func NewCache() *Cache {
	return &Cache{
		cfs: configFS,
	}
}

// Open the cache from the cache.json file in config
// directory
func (c *Cache) Open() {
	raw, err := c.cfs.ReadCache()
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

	if err := c.cfs.WriteCache(raw); err != nil {
		return err
	}
	return nil
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
