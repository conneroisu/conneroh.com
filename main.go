// Package main contains the main function for the conneroh.com website.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
	"github.com/conneroisu/conneroh.com/internal/logger"
)

func main() {
	var err error
	slog.SetDefault(logger.DefaultProdLogger)
	if err = conneroh.Run(
		context.Background(),
		os.Getenv,
	); err != nil {
		fmt.Println(err)

		return
	}
}
