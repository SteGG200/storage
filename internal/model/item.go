package model

// Item represents a file or directory in the storage system.
type Item struct {
	Name        string `json:"name"`
	Path        string `json:"path"` // Relative path from the storage root
	Size        int64  `json:"size"`
	IsDirectory bool   `json:"isDirectory"`
	ModifiedAt  string `json:"modifiedAt"` // Formatted in RFC3339
}
