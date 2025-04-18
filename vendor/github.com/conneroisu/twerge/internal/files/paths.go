package files

import (
	"fmt"
	"path/filepath"
)

// GetPackageName extracts the package name from a Go file at the given path
func GetPackageName(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("file path is empty")
	}

	// Get the directory containing the file
	dir := filepath.Dir(filePath)

	// Check if the directory is a Windows drive or a special case
	if dir == filepath.VolumeName(dir) {
		return "", fmt.Errorf("file path is a windows drive")
	}

	// Check if the directory is "." or "/" or "\"
	if dir == "." || dir == "/" || dir == "\\" {
		return "", fmt.Errorf("file path is a special case that is unsupported")
	}
	// Get the base name of the directory (last component)
	packageName := filepath.Base(dir)

	return packageName, nil
}
