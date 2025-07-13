package create

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/exception"
)

func createFolder(path string, foldername string) error {
	_, err := os.Stat(path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	_, err = os.Stat(filepath.Join(path, foldername))

	if !errors.Is(err, fs.ErrNotExist) {
		logger.ErrorLogger.Print(exception.ErrFolderExist)
		return exception.ErrFolderExist
	}

	err = os.Mkdir(filepath.Join(path, foldername), 0744)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}
