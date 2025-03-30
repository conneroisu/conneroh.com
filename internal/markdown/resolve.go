package markdown

import (
	"strings"

	"go.abhg.dev/goldmark/wikilink"
)

// CustomResolver is a wikilink.Resolver that uses the default wikilink
// resolver but also resolves the "custom" wikilink type.
//
// It implements the wikilink.Resolver interface.
//
//	type Resolver interface {
//		// ResolveWikilink returns the address of the page that the provided
//		// wikilink points to. The destination will be URL-escaped before
//		// being placed into a link.
//		//
//		// If ResolveWikilink returns a non-nil error, rendering will be
//		// halted.
//		//
//		// If ResolveWikilink returns a nil destination and error, the
//		// Renderer will omit the link and render its contents as a regular
//		// string.
//		ResolveWikilink(*Node) (destination []byte, err error)
//	}
//
// ```
// ─────────────────────────────────────────────────────────────────────────────────────────────────────
// Resolver resolves pages referenced by wikilinks to their destinations.
// ─────────────────────────────────────────────────────────────────────────────────────────────────────
type CustomResolver struct {
	Assets []Asset
}

// NewCustomResolver creates a new wikilink resolver.
func NewCustomResolver(assets []Asset) *CustomResolver {
	return &CustomResolver{Assets: assets}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
//
// If ResolveWikilink returns a non-nil error, rendering will be
// halted.
//
// If ResolveWikilink returns a nil destination and error, the
// Renderer will omit the link and render its contents as a regular
// string.
func (c *CustomResolver) ResolveWikilink(n *wikilink.Node) (destination []byte, err error) {
	targetStr := string(n.Target)
	for _, asset := range c.Assets {
		if targetStr == asset.Path || strings.HasSuffix(targetStr, asset.Filename()) {
			return []byte(asset.URL()), nil
		}
	}

	return nil, nil
}
