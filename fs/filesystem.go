package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

func MoveFile(filename, destination string) error {
	ext := filepath.Ext(filename)
	fmt.Printf("moving a file of type %s\n", ext)

	err := os.Rename(filename, fmt.Sprintf("%s/%s", destination, filepath.Base(filename)))
	if err != nil {
		return err
	}
	return nil
}

func LookForFile(directory string) {
}
