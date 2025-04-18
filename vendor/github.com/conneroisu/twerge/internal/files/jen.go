package files

import (
	"fmt"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
)

// WriteJen writes a jen.File to a file.
func WriteJen(jenFile *jen.File, path string) error {
	var buf strings.Builder
	err := jenFile.Render(&buf)
	if err != nil {
		return fmt.Errorf("error generating code: %w", err)
	}
	fil, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = fil.WriteString(buf.String())
	if err != nil {
		return err
	}
	err = fil.Close()
	if err != nil {
		return err
	}
	return nil
}
