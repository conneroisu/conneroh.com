package routing

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/internal/hx"
)

var (
	// defaultMakeAPIErrFn is the default api error function.
	defaultMakeAPIErrFn = func(err error, r *http.Request) {
		// print the type of the error
		slog.Error(
			"api error",
			"err",
			err.Error(),
			"method",
			r.Method,
			"path",
			r.URL.Path,
		)
	}
)

// APIFunc is a function that handles an API request and returns an error.
type APIFunc func(http.ResponseWriter, *http.Request) error

// Make returns a handler that calls the given function.
func Make(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			defaultMakeAPIErrFn(err, r)
		}
	}
}

// MorphableHandler returns a handler that checks for the presence of the
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

// BytesHandler returns a handler that writes the given bytes to the response.
func BytesHandler(b []byte) http.HandlerFunc {
	return Make(func(w http.ResponseWriter, r *http.Request) error {
		_, err := w.Write(b)
		return err
	})
}
