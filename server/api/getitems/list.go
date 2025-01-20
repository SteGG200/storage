package getitems

import (
	"os"
	"time"

	"github.com/SteGG200/storage/logger"
)

type Item struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Date        time.Time `json:"date"`
	IsDirectory bool      `json:"isDirectory"`
}

func listItems(pattern string) (items []Item, err error) {
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
			dirSize, err := getDirectorySize(pattern + "/" + entry.Name())

			if err != nil {
				logger.ErrorLogger.Print(err)
				return nil, err
			}

			size = dirSize
		}

		items = append(items, Item{
			Name:        entry.Name(),
			Size:        size,
			Date:        info.ModTime(),
			IsDirectory: entry.IsDir(),
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
			dirSize, err := getDirectorySize(pattern + "/" + entry.Name())

			if err != nil {
				logger.ErrorLogger.Print(err)
				return 0, err
			}

			size += dirSize
		}
	}

	return
}
