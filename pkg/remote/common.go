package remote

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"
)

var (
	// ErrAuth is an error returned when an authentication error status
	// code is returned from server
	ErrAuth = errors.New("bad authentication provided")

	// ErrEndOfRepos is an error returned when there is no next Link header
	// indicating that there are no additional pages of repos
	ErrEndOfRepos = errors.New("no next link found, no more repos available")

	// ErrParsingResponse is returned when a response is received from the remote
	// but that response could not be parsed
	ErrParsingResponse = errors.New("response from remote could not be parsed")

	// ErrRemoteRequest returned when there is an error with creating or making
	// the request to the remote
	ErrRemoteRequest = errors.New("bad request to remote")
)

func getRepoHelper(url string, acc []map[string]string, mapper func(map[string]interface{}) map[string]string) ([]map[string]string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return acc, ErrRemoteRequest
	}
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return acc, ErrRemoteRequest
	} else if res.StatusCode >= 400 {
		if res.StatusCode == 401 || res.StatusCode == 403 {
			return acc, ErrAuth
		}
		return acc, ErrRemoteRequest
	}

	defer res.Body.Close()

	repos := make([]map[string]interface{}, 20)
	if err = json.NewDecoder(res.Body).Decode(&repos); err != nil {
		return acc, ErrParsingResponse
	}

	for _, repo := range repos {
		acc = append(acc, mapper(repo))
	}

	nextURL, err := parseLinkHeader(res.Header.Get("link"))
	if err == ErrEndOfRepos {
		return acc, nil
	} else if err != nil {
		return acc, err
	}

	return getRepoHelper(nextURL, acc, mapper)
}

func parseLinkHeader(header string) (string, error) {
	rg := regexp.MustCompile(".*<(.+)>;(?: ?)rel=\"next\"(?:|.+)$")
	next := rg.FindStringSubmatch(header)
	if len(next) > 1 {
		return next[1], nil
	}
	return "", ErrEndOfRepos
}
