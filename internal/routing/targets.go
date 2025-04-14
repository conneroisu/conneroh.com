package routing

import (
	"fmt"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
)

// PluralPath is the target of a plural view.
// string
type PluralPath = string

const (
	// PostPluralPath is the target of a plural post view.
	PostPluralPath PluralPath = "posts"
	// ProjectPluralPath is the target of a plural project view.
	ProjectPluralPath PluralPath = "projects"
	// TagsPluralPath is the target of a plural tag view.
	TagsPluralPath PluralPath = "tags"
)

// GetPostURL returns the URL for a post.
func GetPostURL(base string, post *gen.Post) string {
	return fmt.Sprintf("%s/post/%s", base, post.Slug)
}

// GetProjectURL returns the URL for a project.
func GetProjectURL(base string, project *gen.Project) string {
	return fmt.Sprintf("%s/project/%s", base, project.Slug)
}

// GetTagURL returns the URL for a tag.
func GetTagURL(base string, tag *gen.Tag) string {
	return fmt.Sprintf("%s/tag/%s", base, tag.Slug)
}
