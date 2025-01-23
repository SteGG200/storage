package download

import (
	"os"

	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/exception"
)

func readFile(path string) (file *os.File, err error) {
	stat, err := os.Stat(path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	if stat.IsDir() {
		logger.ErrorLogger.Printf("%s is not a file", path)
		return nil, exception.ErrNotAFile
	}

	file, err = os.Open(path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	return file, nil
}
