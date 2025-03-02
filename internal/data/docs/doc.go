// Package docs contains documentation and updater.
package docs

import "embed"

// Posts contains all posts.
//
//go:embed posts/*.md
var Posts embed.FS

// Projects contains all projects.
//
//go:embed projects/*.md
var Projects embed.FS

// Tags contains all tags.
//
//go:embed tags/*.md
var Tags embed.FS
