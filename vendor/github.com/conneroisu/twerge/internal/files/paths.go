package files

import (
	"path/filepath"
	"strings"
)

// GetPackageName returns the package name from a file path.
func GetPackageName(path string) string {
	// Get the package name from the file path
	fileName := filepath.Base(path)
	pkgEnd := strings.LastIndex(fileName, ".")
	if pkgEnd == -1 {
		return "main"
	}
	return fileName[:pkgEnd]
}
