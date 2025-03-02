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

// func main() {
// 	db, err := conneroh.NewDb(os.Getenv)
// 	if err != nil {
// 		return
// 	}
// 	defer db.Close()
//
// 	client, err := api.ClientFromEnvironment()
// 	if err != nil {
// 		return
// 	}
// 	id, err := data.UpsertEmbedding(
// 		context.Background(),
// 		db,
// 		client,
// 		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
// 	)
// 	slog.Info("id", slog.String("id", string(id)))
// }
