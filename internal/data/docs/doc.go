// Package docs contains documentation and updater.
package docs

import (
	"embed"
)

// TODO: Remove these embeds?

// Posts contains all posts.
//
//go:embed posts/*
var Posts embed.FS

// Projects contains all projects.
//
//go:embed projects/*
var Projects embed.FS

// Tags contains all tags.
//
//go:embed tags/*
var Tags embed.FS

// Assets contains all media assets.
//
//go:embed assets/*
var Assets embed.FS
