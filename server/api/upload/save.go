package upload

import (
	"errors"
	"io/fs"
	"os"

	"github.com/SteGG200/storage/logger"
)

func createFile(path string, filename string) error {
	_, err := os.Stat(path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	_, err = os.Stat(path + "/" + filename)

	if !errors.Is(err, fs.ErrNotExist) {
		logger.ErrorLogger.Print(fs.ErrExist)
		return fs.ErrExist
	}

	_, err = os.Create(path + "/" + filename)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}

func saveFile(data []byte, path string, filename string) error {
	file, err := os.OpenFile(path+"/"+filename, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	defer file.Close()

	size, err := file.Write(data)

	if size < len(data) {
		logger.WarningLogger.Print("Cannot write entire file.")
	}

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}
