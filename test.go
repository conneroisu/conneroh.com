package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/altipla-consulting/directus-go/v2"
)

func main() {
	ctx := context.Background()
	client := directus.NewClient(
		"https://cms-conneroh.fly.dev",
		"8MS1ev0KiHHAJEGvjWjj_0zkzbExUZWs",
		directus.WithBodyLogger(),
		directus.WithLogger(slog.Default()),
	)
	list, err := client.Collections.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, collection := range list {
		log.Println(collection.Collection)
	}
	users, err := client.Users.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		log.Println(user.FirstName + " " + user.LastName)
	}
}
