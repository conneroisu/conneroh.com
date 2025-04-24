package main_test

import (
	"database/sql"
	"testing"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

func TestGetPosts(t *testing.T) {
	sqlDB, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		t.Fatal(eris.Wrap(err, "error opening database"))
	}
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	assets.RegisterModels(db)

	posts := []*assets.Post{}
	err = db.NewSelect().
		Model(&posts).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(t.Context())

	if err != nil {
		t.Fatal(eris.Wrap(err, "error getting posts"))
	}
	if len(posts) == 0 {
		t.Fatal("expected 0 posts, got", len(posts))
	}
	t.Log("posts", posts)
	for _, post := range posts {
		if post.Title == "A Reflective Journey - Navigating Your Cumulative Experience at Iowa State University" {
			if len(post.Tags) < 2 {
				t.Fatal("expected at least 2 tags, got", len(post.Tags))
			}
		}
		t.Log("post", post)
		t.Log("post.Tags", post.Tags)
		t.Log("post.Posts", post.Posts)
		t.Log("post.Projects", post.Projects)
	}
}
