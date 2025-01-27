package create

import (
	"errors"
	"io/fs"
	"os"

	"github.com/SteGG200/storage/logger"
)

func createFolder(path string, foldername string) error {
	_, err := os.Stat(path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	_, err = os.Stat(path + "/" + foldername)

	if !errors.Is(err, fs.ErrNotExist) {
		logger.ErrorLogger.Print(fs.ErrExist)
		return fs.ErrExist
	}

	err = os.Mkdir(path+"/"+foldername, 0777)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}
