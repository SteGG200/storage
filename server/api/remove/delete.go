package remove

import (
	"os"

	"github.com/SteGG200/storage/logger"
)

func removeItem(path string) error {
	_, err := os.Stat(path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	err = os.RemoveAll(path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}
