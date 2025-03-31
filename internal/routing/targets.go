package routing

import (
	"fmt"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
)

// SingleTarget is the target of a single view.
// string
type SingleTarget = string

const (
	// SingleTargetPost is the target of a single post view.
	SingleTargetPost SingleTarget = "post"
	// SingleTargetProject is the target of a single project view.
	SingleTargetProject SingleTarget = "project"
	// SingleTargetTag is the target of a single tag view.
	SingleTargetTag SingleTarget = "tag"
)

// PluralTarget is the target of a plural view.
// string
type PluralTarget = string

const (
	// PluralTargetPost is the target of a plural post view.
	PluralTargetPost PluralTarget = "posts"
	// PluralTargetProject is the target of a plural project view.
	PluralTargetProject PluralTarget = "projects"
	// PluralTargetTag is the target of a plural tag view.
	PluralTargetTag PluralTarget = "tags"
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
