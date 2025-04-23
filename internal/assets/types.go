package assets

import (
	"mime"
	"path/filepath"
	"slices"
)

var (
	AllowedAssetTypes = []string{
		"image/png",
		"image/jpeg",
		"image/gif",
		"image/webp",
		"image/avif",
		"image/tiff",
		"image/svg+xml",
		"application/pdf",
	}
	AllowedDocumentTypes = []string{
		"text/markdown",
		"text/html",
	}
)

// IsAllowedMediaType returns true if the provided path is an allowed asset type.
func IsAllowedMediaType(path string) bool {
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	return slices.Contains(AllowedAssetTypes, contentType)
}

// IsAllowedAsset returns true if the provided path is an allowed asset.
func IsAllowedAsset(path string) bool {
	return IsAllowedMediaType(path) || IsAllowedDocumentType(path)
}

// IsAllowedDocumentType returns true if the provided path is an allowed document type.
func IsAllowedDocumentType(path string) bool {
	return slices.Contains(AllowedDocumentTypes, mime.TypeByExtension(filepath.Ext(path)))
}
