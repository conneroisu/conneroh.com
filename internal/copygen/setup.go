// Package copygen is a code generator for the assets package.
package copygen

import "github.com/conneroisu/conneroh.com/internal/assets"

/* Copygen defines the functions that are generated. */
type Copygen interface {
	// depth domain.Post 2
	ToPost(*assets.Doc) *assets.Post
	// depth domain.Project 2
	ToProject(*assets.Doc) *assets.Project
	// depth domain.Tag 2
	ToTag(*assets.Doc) *assets.Tag
	// depth domain.Employment 2
	ToEmployment(*assets.Doc) *assets.Employment
}
