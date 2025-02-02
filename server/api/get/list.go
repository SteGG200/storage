package get

import (
	"os"
	"path/filepath"
	"time"

	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/exception"
)

type Item struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Date        time.Time `json:"date"`
	IsDirectory bool      `json:"isDirectory"`
}

func listEntries(pattern string) (items []Item, err error) {
	stat, err := os.Stat(pattern)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	if !stat.IsDir() {
		logger.ErrorLogger.Printf("%s is not a directory", pattern)
		return nil, exception.ErrNotADirectory
	}

	entries, err := os.ReadDir(pattern)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	items = make([]Item, 0)

	for _, entry := range entries {
		info, err := entry.Info()

		if err != nil {
			logger.ErrorLogger.Print(err)
			return nil, err
		}

		var size int64

		if !info.IsDir() {
			size = info.Size()
		} else {
			dirSize, err := getDirectorySize(filepath.Join(pattern, info.Name()))

			if err != nil {
				logger.ErrorLogger.Print(err)
				return nil, err
			}

			size = dirSize
		}

		items = append(items, Item{
			Name:        info.Name(),
			Size:        size,
			Date:        info.ModTime(),
			IsDirectory: info.IsDir(),
		})
	}

	return
}

func getDirectorySize(pattern string) (size int64, err error) {
	entries, err := os.ReadDir(pattern)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return 0, err
	}

	size = 0

	for _, entry := range entries {
		info, err := entry.Info()

		if err != nil {
			logger.ErrorLogger.Print(err)
			return 0, err
		}

		if !info.IsDir() {
			size += info.Size()
		} else {
			dirSize, err := getDirectorySize(filepath.Join(pattern, info.Name()))

			if err != nil {
				logger.ErrorLogger.Print(err)
				return 0, err
			}

			size += dirSize
		}
	}

	return
}
