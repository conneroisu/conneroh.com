package search

import "net/http"

// Params is the parameters for searchable lists.
type Params struct {
	PostReq    []string
	ProjectReq []string
	TagReq     []string
}

// FromRequest returns a Params struct from the request.
func FromRequest(r *http.Request) Params {
	return Params{
		PostReq:    r.Form["post"],
		ProjectReq: r.Form["project"],
		TagReq:     r.Form["tag"],
	}
}
