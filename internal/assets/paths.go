package assets

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	// VaultLoc is the location of the vault.
	// This is the location of the documents and assets.
	VaultLoc = "internal/data/docs/"
	// AssetsLoc is the location of the assets relative to the vault.
	AssetsLoc = "assets/"
	// PostsLoc is the location of the posts relative to the vault.
	PostsLoc = "posts/"
	// TagsLoc is the location of the tags relative to the vault.
	TagsLoc = "tags/"
	// ProjectsLoc is the location of the projects relative to the vault.
	ProjectsLoc = "projects/"
)

// Slugify returns the slugified path of a document or media asset.
func Slugify(s string) string {
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, AssetsLoc)
	if ok {
		return path
	}

	return strings.TrimSuffix(Pathify(s), filepath.Ext(s))
}

// Pathify returns the path to the document page or media asset page.
func Pathify(s string) string {
	var (
		path string
		ok   bool
	)

	path, ok = strings.CutPrefix(s, PostsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, ProjectsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, TagsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, AssetsLoc)
	if ok {
		return path
	}
	panic(fmt.Errorf("failed to pathify %s", s))
}
