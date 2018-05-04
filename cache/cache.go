package cache

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	cacheFormatVersion = "0.0.1"
)

var (
	cache Cache
)

// Cache holds the cache of remotes and their repos
type Cache struct {
	Version string             `json:"version"`
	Remotes map[string]*Remote `json:"remotes`
}

// Remote is a parent node in the cache tree
type Remote struct {
	Name  string   `json:"name"`
	URL   *url.URL `json:"url"`
	Repos []Repo   `json:"repos"`
}

// Repo stores a repo and its location on the filesystem
// for use in autocomplete
type Repo struct {
	Name string `json:"name"`
	Path string `json:"repo_path"`
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
		URL:  remoteURL,
		Repos: []Repo{
			Repo{
				Name: name,
				Path: repoPath,
			},
		},
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
