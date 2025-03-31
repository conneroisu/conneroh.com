package main

import (
	"testing"
)

func TestAssetFilename(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Simple filename",
			path:     "test.png",
			expected: "test.png",
		},
		{
			name:     "Path with directory",
			path:     "images/test.png",
			expected: "test.png",
		},
		{
			name:     "Nested directories",
			path:     "assets/images/profile/test.png",
			expected: "test.png",
		},
		{
			name:     "Path with no extension",
			path:     "images/document",
			expected: "document",
		},
		{
			name:     "Path with multiple extensions",
			path:     "images/archive.tar.gz",
			expected: "archive.tar.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset := Asset{
				Path: tt.path,
				Data: []byte("test data"),
			}

			if got := asset.Filename(); got != tt.expected {
				t.Errorf("Asset.Filename() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAssetCreation(t *testing.T) {
	// Test creating an asset with path and data
	path := "images/test.png"
	data := []byte("test image data")

	asset := Asset{
		Path: path,
		Data: data,
	}

	if asset.Path != path {
		t.Errorf("Asset.Path = %v, want %v", asset.Path, path)
	}

	if string(asset.Data) != string(data) {
		t.Errorf("Asset.Data = %v, want %v", string(asset.Data), string(data))
	}
}
