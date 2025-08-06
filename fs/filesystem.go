package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

func MoveFile(file_path, destination string) error {
	ext := filepath.Ext(file_path)
	fmt.Printf("moving a file of type %s\n", ext)

	err := os.Rename(file_path, fmt.Sprintf("%s/%s", destination, filepath.Base(file_path)))
	if err != nil {
		return err
	}
	return nil
}

func LookForFile(directory string) {
}
