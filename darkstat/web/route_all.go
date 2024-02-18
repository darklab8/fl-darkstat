package web

import (
	"fmt"
	"html"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkstat/settings/types"
	"github.com/darklab8/fl-darkstat/darkstat/web/registry"
)

const UrlStatic types.Url = "/"

func init() {
	registry.Registry.Register(NewEndpointStatic())
}

func NewEndpointStatic() *registry.Endpoint {
	return &registry.Endpoint{
		Url: UrlStatic,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				fmt.Fprintf(w, "static at %s!, %q", r.URL, html.EscapeString(r.URL.Path))
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	}
}
