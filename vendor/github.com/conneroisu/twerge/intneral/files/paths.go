package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetPackageName extracts the package name from a Go file at the given path
func GetPackageName(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Scan through the file looking for the package declaration
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || line == "" {
			continue
		}

		// Check if the line declares a package
		if strings.HasPrefix(line, "package ") {
			// Extract the package name
			packageName := strings.TrimSpace(strings.TrimPrefix(line, "package "))

			// Remove any trailing comments
			if idx := strings.Index(packageName, "//"); idx >= 0 {
				packageName = strings.TrimSpace(packageName[:idx])
			}

			return packageName, nil
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	// If no package declaration was found
	return "", fmt.Errorf("no package declaration found in file")
}
