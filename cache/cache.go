package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/TheHipbot/hermes/fs"
)

const (
	cacheFormatVersion = "0.0.1"
)

var (
	cache    Cache
	configFS *fs.ConfigFS
)

func init() {
	configFS = fs.NewConfigFS()
	cache = initCache(configFS.ReadCache())
}

func initCache(raw []byte, err error) Cache {
	var result Cache
	if err != nil {
		result = Cache{
			Version: cacheFormatVersion,
		}
	} else {
		result = Cache{}
		if err := json.Unmarshal(raw, &result); err != nil {
			fmt.Print(err)
			result = Cache{
				Version: cacheFormatVersion,
			}
		}
	}
	return result
}

// Cache holds the cache of remotes and their repos
type Cache struct {
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

func (c *Cache) save() error {
	raw, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := configFS.WriteCache(raw); err != nil {
		return err
	}
	return nil
}

// Add a repo to the cache
func Add(name, path string) error {
	repoPath := fmt.Sprintf("%s%s", path, name)
	remote := strings.Split(name, "/")[0]

	if r, ok := cache.Remotes[remote]; ok {
		r.Repos = append(r.Repos, Repo{
			Name: name,
			Path: repoPath,
		})
		return nil
	}

	remoteURL, err := url.Parse(fmt.Sprintf("https://%s", remote))
	if err != nil {
		return err
	}

	cache.Remotes[remote] = &Remote{
		Name: remote,
		URL:  remoteURL.String(),
		Repos: []Repo{
			Repo{
				Name: name,
				Path: repoPath,
			},
		},
	}

	if err := cache.save(); err != nil {
		return err
	}

	return nil
}

// Remove a repo from the cache
func Remove(name string) error {
	found := false
	remote := strings.Split(name, "/")[0]

	if r, ok := cache.Remotes[remote]; ok {
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
func Search(needle string) []Repo {
	var results []Repo
	for _, remote := range cache.Remotes {
		for _, repo := range remote.Repos {
			if strings.Contains(repo.Name, needle) {
				results = append(results, repo)
			}
		}
	}
	return results
}
