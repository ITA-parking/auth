package util

import (
	"errors"
	"fmt"
	"path/filepath"
)

func GetAbsolutePath(relPath string) (string, error) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return "", errors.Join(errors.New(fmt.Sprintf("Failed to get absolute path for %s", relPath)), err)

	}
	return absPath, nil
}
