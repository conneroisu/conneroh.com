package routing

import (
	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/internal/data/master"
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

// PluralStaticFn is a function that returns a static component for a plural view.
type PluralStaticFn func(
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostsSlugMap *map[string]master.FullPost,
	fullProjectsSlugMap *map[string]master.FullProject,
	fullTagsSlugMap *map[string]master.FullTag,
) templ.Component
