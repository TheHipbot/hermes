package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	// URL format for github repo requests
	defaultGitHubAPIHost = "https://api.github.com"
	gitHubUserRequestFmt = "/user/repos?access_token=%s&page=%d"
)

func githubCreator() (Driver, error) {
	return &GitHub{
		Host: defaultGitHubAPIHost,
	}, nil
}

// GitHub is a client to github
type GitHub struct {
	Auth
	Host string
}

// SetHost sets github driver host to provided string
func (gh *GitHub) SetHost(host string) {
	gh.Host = host
}

// Authenticate sets Auth object for driver
func (gh *GitHub) Authenticate(a Auth) {
	gh.Auth = a
}

// AuthType sets Auth object for driver
func (gh *GitHub) AuthType() string {
	return authToken
}

// GetRepos gets the repos for the github user
func (gh *GitHub) GetRepos() ([]map[string]string, error) {
	urlFormat := fmt.Sprintf("%s%s", gh.Host, gitHubUserRequestFmt)
	if gh.Auth.Token == "" && gh.Auth.Username == "" {
		return nil, errors.New("Auth is empty")
	}

	page := 1
	accumulator := []map[string]string{}
	return getRepoHelper(fmt.Sprintf(urlFormat, gh.Auth.Token, page), accumulator)
}

func getRepoHelper(url string, acc []map[string]string) ([]map[string]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return acc, err
	}
	defer res.Body.Close()

	repos := make([]map[string]interface{}, 20)
	json.NewDecoder(res.Body).Decode(&repos)
	for _, repo := range repos {
		entry := make(map[string]string, 3)
		url := repo["html_url"].(string)
		entry["url"] = url
		entry["name"] = strings.Split(url, "://")[1]
		acc = append(acc, entry)
	}

	nextURL, err := parseLinkHeader(res.Header.Get("link"))
	if err != nil {
		return acc, nil
	}
	return getRepoHelper(nextURL, acc)
}

func parseLinkHeader(header string) (string, error) {
	rg := regexp.MustCompile("<(.+)>; rel=\"next\",")
	next := rg.FindStringSubmatch(header)
	if len(next) > 0 {
		return next[1], nil
	}
	return "", errors.New("No next link found")
}
