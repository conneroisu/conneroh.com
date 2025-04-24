//go:build !dev
// +build !dev

// Package main contains the main function for the conneroh.com website.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
)

func main() {
	if err := conneroh.Run(
		context.Background(),
		os.Getenv,
	); err != nil {
		fmt.Println(err)
		return
	}
}
