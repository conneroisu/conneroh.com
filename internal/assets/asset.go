package assets

import "fmt"

// Asset is a struct for embedding assets.
type Asset struct {
	Slug string
	Path string
	Data []byte
}

// URL returns the url of the asset.
func (a *Asset) URL() string {
	return fmt.Sprintf(
		"https://conneroh.fly.storage.tigris.dev/%s",
		a.Path,
	)
}
