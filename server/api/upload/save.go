package upload

import (
	"errors"
	"io/fs"
	"os"
)

func saveFile(data []byte, filename string) error {
	if _, err := os.Stat(filename); !errors.Is(err, fs.ErrNotExist) {
		return fs.ErrExist
	}

	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	file.Write(data)

	return nil
}
