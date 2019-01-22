package storage

// Remote is a parent node in the cache tree
type Remote struct {
	Name     string            `json:"name"`
	URL      string            `json:"url"`
	Protocol string            `json:"protocol"`
	Port     string            `json:"port"`
	Type     string            `json:"type"`
	Meta     map[string]string `json:"meta"`
	Repos    []Repository      `json:"repos"`
}
