package assets

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/credited"
	"go.abhg.dev/goldmark/wikilink"
)

// resolver is a wikilink.Resolver that resolves pages and media referenced by
// wikilinks to their destinations.
type resolver struct{ client *credited.AWSClient }

// newResolver creates a new wikilink resolver.
func newResolver(client *credited.AWSClient) *resolver { return &resolver{client: client} }

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (c *resolver) ResolveWikilink(n *wikilink.Node) (destination []byte, err error) {
	targetStr := string(n.Target)
	_, err = c.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("conneroh"),
		Key:    aws.String(targetStr),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object(%s): %w", targetStr, err)
	}
	return fmt.Appendf(nil,
		"https://conneroh.fly.storage.tigris.dev/%s",
		targetStr,
	), nil
}
