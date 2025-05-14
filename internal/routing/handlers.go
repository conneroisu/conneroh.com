package routing

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/internal/hx"
)

// APIFunc is a function that handles an API request and returns an error.
type APIFunc func(http.ResponseWriter, *http.Request) error

// Make returns a handler that calls the given function.
func Make(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			slog.Error(
				"api error",
				"err",
				err.Error(),
				"method",
				r.Method,
				"path",
				r.URL.Path,
			)
			target := r.Header.Get(hx.HdrTarget)
			if target != "" {
				http.Redirect(w, r, "/500", http.StatusFound) // 302 Found or http.StatusTemporaryRedirect (307)
			}
		}
	}
}

// MorphableHandler returns a handler that checks for the presence of the
// htmx request header and either serves the morphed component or the
// component wrapper.
func MorphableHandler(
	wrapper func(comp templ.Component) templ.Component,
	morph templ.Component,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get(hx.HdrRequest)
		if header == "" {
			templ.Handler(wrapper(morph)).ServeHTTP(w, r)
		} else {
			templ.Handler(morph).ServeHTTP(w, r)
		}
	}
}
