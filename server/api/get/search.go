package get

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/exception"
)

func searchItem(root, searchValue string) (items []Item, err error) {
	stat, err := os.Stat(root)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	if !stat.IsDir() {
		logger.ErrorLogger.Printf("%s is not a directory", root)
		return nil, exception.ErrNotADirectory
	}

	pattern, err := regexp.Compile(fmt.Sprintf("(?i).*(%s).*", searchValue))

	if err != nil {
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	items = make([]Item, 0)
	stack := make([]Item, 0)

	var searchFunc func(path string) error

	searchFunc = func(path string) error {
		entries, err := os.ReadDir(path)

		if err != nil {
			logger.ErrorLogger.Print(err)
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				if pattern.MatchString(entry.Name()) {
					info, err := entry.Info()

					if err != nil {
						logger.ErrorLogger.Print(err)
						return err
					}

					relativePath, err := filepath.Rel(root, path)

					if err != nil {
						logger.ErrorLogger.Print(err)
						return err
					}

					stack = append(stack, Item{
						Name:        info.Name(),
						Size:        int64(0),
						Date:        info.ModTime(),
						IsDirectory: true,
						Path:        relativePath,
					})

					err = searchFunc(filepath.Join(path, info.Name()))

					if err != nil {
						logger.InfoLogger.Print(err)
						return err
					}

					temp := stack[len(stack)-1]

					items = append(items, temp)

					stack = stack[:len(stack)-1]

					if len(stack) > 0 {
						stack[len(stack)-1].Size += temp.Size
					}
				} else {
					err = searchFunc(filepath.Join(path, entry.Name()))

					if err != nil {
						logger.InfoLogger.Print(err)
						return err
					}
				}
			} else {
				info, err := entry.Info()

				if err != nil {
					logger.ErrorLogger.Print(err)
					return err
				}

				relativePath, err := filepath.Rel(root, path)

				if err != nil {
					logger.ErrorLogger.Print(err)
					return err
				}

				if pattern.MatchString(entry.Name()) {
					items = append(items, Item{
						Name:        info.Name(),
						Size:        info.Size(),
						Date:        info.ModTime(),
						IsDirectory: false,
						Path:        relativePath,
					})
				}

				if len(stack) > 0 {
					stack[len(stack)-1].Size += info.Size()
				}
			}
		}

		return nil
	}

	err = searchFunc(root)

	return
}
