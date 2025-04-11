package routing

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/internal/hx"
)

// MorphableHandler returns a handler that checks for the presence of the
// hx-trigger header and serves either the full or morphed view.
func MorphableHandler(
	full templ.Component,
	morph templ.Component,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get(hx.HdrRequest)
		if header == "" {
			templ.Handler(full).ServeHTTPStreamed(w, r)
		} else {
			templ.Handler(morph).ServeHTTP(w, r)
		}
	}
}
