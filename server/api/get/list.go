package get

import (
	"os"

	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/exception"
)

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
			size = 0
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
