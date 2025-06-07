package validate

import (
	"fmt"
	"os"
)

func ValidateFile(filePath string) error {
	info, err := os.Stat(filePath)
	// if the file does not exist or cannot be accessed, return an error
	if err != nil {
		return fmt.Errorf("file does not exist or cannot be accessed: %w", err)
	}
	// if the file is a directory, return an error
	if info.IsDir() {
		return fmt.Errorf("file is a directory, not a file")
	}
	//
	return nil
}