package remote

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"
)

func getRepoHelper(url string, acc []map[string]string, mapper func(map[string]interface{}) map[string]string) ([]map[string]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return acc, err
	}
	defer res.Body.Close()

	repos := make([]map[string]interface{}, 20)
	json.NewDecoder(res.Body).Decode(&repos)
	for _, repo := range repos {
		// entry := make(map[string]string, 3)
		// url := repo["html_url"].(string)
		// entry["url"] = url
		// entry["name"] = strings.Split(url, "://")[1]
		acc = append(acc, mapper(repo))
	}

	nextURL, err := parseLinkHeader(res.Header.Get("link"))
	if err != nil {
		return acc, nil
	}
	return getRepoHelper(nextURL, acc, mapper)
}

func parseLinkHeader(header string) (string, error) {
	rg := regexp.MustCompile("<(.+)>; rel=\"next\",")
	next := rg.FindStringSubmatch(header)
	if len(next) > 0 {
		return next[1], nil
	}
	return "", errors.New("No next link found")
}
