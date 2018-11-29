package storage

// Repository stores a repo and its location on the filesystem
// for use in autocomplete
type Repository struct {
	Name string `json:"name"`
	Path string `json:"repo_path"`
}
