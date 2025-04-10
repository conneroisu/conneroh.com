package assets

import "go.abhg.dev/goldmark/wikilink"

var (
	_ wikilink.Resolver = &resolver{}
)
