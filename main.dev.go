//go:build dev
// +build dev

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
	"github.com/conneroisu/conneroh.com/internal/logger"
)

// main is the main function for dev mode version of the conneroh.com website.
// It compares hashes.
func main() {
	var err error
	slog.SetDefault(logger.DefaultLogger)
	// Run the `update`

	// Run the `server`
	if err = conneroh.Run(
		context.Background(),
		os.Getenv,
	); err != nil {
		fmt.Println(err)
		return
	}
}
