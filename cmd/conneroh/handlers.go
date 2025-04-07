package conneroh

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

const hName = "HX-Trigger-Name"

func searchHandler(
	target routing.PluralTarget,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("search")
		header := r.Header.Get(hName)
		if header == "" {
			switch target {
			case routing.PluralTargetPost:
				templ.Handler(layouts.Page(views.List(target, &gen.AllPosts, nil, nil, query))).ServeHTTP(w, r)
			case routing.PluralTargetProject:
				templ.Handler(layouts.Page(views.List(target, nil, &gen.AllProjects, nil, query))).ServeHTTP(w, r)
			case routing.PluralTargetTag:
				templ.Handler(layouts.Page(views.List(target, nil, nil, &gen.AllTags, query))).ServeHTTP(w, r)
			}
			return
		}
		switch target {
		case routing.PluralTargetPost:
			templ.Handler(views.Results(target, &gen.AllPosts, nil, nil)).ServeHTTP(w, r)
		case routing.PluralTargetProject:
			templ.Handler(views.Results(target, nil, &gen.AllProjects, nil)).ServeHTTP(w, r)
		case routing.PluralTargetTag:
			templ.Handler(views.Results(target, nil, nil, &gen.AllTags)).ServeHTTP(w, r)
		}
	}
}

func filter(
	embs []gen.Embedded,
	query string,
) []gen.Embedded {
	return embs
}
