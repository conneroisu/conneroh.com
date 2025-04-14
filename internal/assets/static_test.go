package assets_test

import (
	"net/url"
	"testing"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
)

// TestPaths tests that all paths are valid by creating a localhost:8080/{path}
// and checking that the url package can parse it.
func TestPaths(t *testing.T) {
	for _, p := range gen.AllPosts {
		t.Run(p.Slug, func(t *testing.T) {
			t.Logf("testing %s", p.PagePath())
			_, err := url.ParseRequestURI("http://localhost:8080" + p.PagePath())
			if err != nil {
				t.Errorf("invalid path: %s", p.PagePath())
			}
		})
	}
	for _, p := range gen.AllProjects {
		t.Run(p.Slug, func(t *testing.T) {
			t.Logf("testing %s", p.PagePath())
			_, err := url.ParseRequestURI("http://localhost:8080" + p.PagePath())
			if err != nil {
				t.Errorf("invalid path: %s", p.PagePath())
			}
		})
	}
	for _, tag := range gen.AllTags {
		t.Run(tag.Slug, func(t *testing.T) {
			t.Logf("testing %s", tag.PagePath())
			_, err := url.ParseRequestURI("http://localhost:8080" + tag.PagePath())
			if err != nil {
				t.Errorf("invalid path: %s", tag.PagePath())
			}
		})
	}
}
