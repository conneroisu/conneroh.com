package assets

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	// VaultLoc is the location of the vault.
	// This is the location of the documents and assets.
	VaultLoc = "internal/data/"
	// AssetsLoc is the location of the assets relative to the vault.
	AssetsLoc = "assets/"
	// PostsLoc is the location of the posts relative to the vault.
	PostsLoc = "posts/"
	// TagsLoc is the location of the tags relative to the vault.
	TagsLoc = "tags/"
	// ProjectsLoc is the location of the projects relative to the vault.
	ProjectsLoc = "projects/"
)

// Pathify returns the slugified path of a document or media asset.
func Pathify(s string) string {
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, AssetsLoc)
	if ok {
		return path
	}

	return strings.TrimSuffix(Slugify(s), filepath.Ext(s))
}

// Slugify returns the path to the document page or media asset page.
func Slugify(s string) string {
	var (
		path string
		ok   bool
	)
	s = strings.TrimSuffix(s, ".md")

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

// Static mapping of file extensions to content types to avoid reflection.
var contentTypes = map[string]string{
	".jpg":   "image/jpeg",
	".jpeg":  "image/jpeg",
	".png":   "image/png",
	".gif":   "image/gif",
	".svg":   "image/svg+xml",
	".webp":  "image/webp",
	".css":   "text/css",
	".txt":   "text/plain",
	".md":    "text/markdown",
	".pdf":   "application/pdf",
	".xml":   "application/xml",
	".zip":   "application/zip",
	".mp3":   "audio/mpeg",
	".mp4":   "video/mp4",
	".webm":  "video/webm",
	".wav":   "audio/wav",
	".mov":   "video/quicktime",
	".ico":   "image/x-icon",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".ttf":   "font/ttf",
	".otf":   "font/otf",
}

// GetContentType returns the content type for a file extension.
func GetContentType(path string) string {
	ext := filepath.Ext(path)
	if contentType, ok := contentTypes[ext]; ok {
		return contentType
	}
	// Fallback to standard library for unknown extensions
	return "application/octet-stream"
}

// BucketPath returns the path to the bucket for a given file path.
func BucketPath(path string) string {
	return "https://conneroisu.fly.storage.tigris.dev/assets/" + path
}

func isVideoType(contentType string) bool {
	switch contentType {
	case "video/mp4", "video/quicktime", "video/webm":
		return true
	default:
		return false
	}
}
