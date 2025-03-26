package routing

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
