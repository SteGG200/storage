package rename

import (
	"os"
	"path/filepath"

	"github.com/SteGG200/storage/logger"
)

func renameItem(oldPath, newName string) error {
	_, err := os.Stat(oldPath)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	parentsDir := filepath.Dir(oldPath)

	newPath := filepath.Join(parentsDir, newName)

	err = os.Rename(oldPath, newPath)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}
