package storage

// Repository stores a repo and its location on the filesystem
// for use in autocomplete
type Repository struct {
	Name     string `json:"name"`
	Path     string `json:"repo_path"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
}
