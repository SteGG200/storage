package fsutil

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SteGG200/storage/internal/model"
	"github.com/charlievieth/fastwalk"
)

// Search searches for files/directories inside searchRoot (which must be absolute) whose names contain the query.
// It skips symlinks and hardlinks. It returns paths relative to storageRoot.
func Search(storageRoot, searchRoot, query string) ([]model.Item, error) {
	var (
		mu      sync.Mutex
		results []model.Item
	)

	queryLower := strings.ToLower(query)

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip items with errors during search
			return nil
		}

		if path == searchRoot {
			return nil
		}

		// Skip symlinks
		if d.Type()&fs.ModeSymlink != 0 {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		name := d.Name()
		if strings.Contains(strings.ToLower(name), queryLower) {
			info, err := d.Info()
			if err != nil {
				return nil
			}

			// Reject hardlinks
			if !info.IsDir() {
				if err := RejectSymlinksAndHardlinks(path); err != nil {
					return nil
				}
			}

			relPath, err := filepath.Rel(storageRoot, path)
			if err != nil {
				return nil
			}
			relPath = filepath.ToSlash(relPath)

			item := model.Item{
				Name:        name,
				Path:        relPath,
				Size:        info.Size(),
				IsDirectory: info.IsDir(),
				ModifiedAt:  info.ModTime().Format(time.RFC3339),
			}

			mu.Lock()
			results = append(results, item)
			mu.Unlock()
		}

		return nil
	}

	err := fastwalk.Walk(nil, searchRoot, walkFn)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// ListDir lists direct contents of targetDir (which must be absolute).
// It skips symlinks and hardlinks. It returns paths relative to storageRoot.
func ListDir(storageRoot, targetDir string) ([]model.Item, error) {
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, err
	}

	results := make([]model.Item, 0, len(entries))
	for _, entry := range entries {
		path := filepath.Join(targetDir, entry.Name())

		if entry.Type()&fs.ModeSymlink != 0 {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if !info.IsDir() {
			if err := RejectSymlinksAndHardlinks(path); err != nil {
				continue
			}
		}

		relPath, err := filepath.Rel(storageRoot, path)
		if err != nil {
			continue
		}
		relPath = filepath.ToSlash(relPath)

		results = append(results, model.Item{
			Name:        entry.Name(),
			Path:        relPath,
			Size:        info.Size(),
			IsDirectory: info.IsDir(),
			ModifiedAt:  info.ModTime().Format(time.RFC3339),
		})
	}

	return results, nil
}
