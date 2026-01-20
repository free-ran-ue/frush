package util

import (
	"fmt"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FileRead(path string) ([]byte, error) {
	if !FileExists(path) {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return os.ReadFile(path)
}
