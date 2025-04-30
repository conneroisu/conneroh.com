package routing

import "net/http"

const (
	// SlugVarName is the name/path-key of the slug variable.
	SlugVarName = "slug"
)

// Slug returns the slug of the tag from the request.
func Slug(r *http.Request) string {
	return r.PathValue(SlugVarName)
}
