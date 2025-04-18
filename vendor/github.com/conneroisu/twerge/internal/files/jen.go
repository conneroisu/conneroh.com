package files

import (
	"fmt"
	"os"

	"github.com/dave/jennifer/jen"
)

// JenFile renders a jen.File to a file at the given path.
func JenFile(
	fil *jen.File,
	path string,
) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()
	return fil.Render(file)
}
