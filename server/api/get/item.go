package get

import "time"

type Item struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Date        time.Time `json:"date"`
	IsDirectory bool      `json:"isDirectory"`
	Path        string    `json:"path"`
}
