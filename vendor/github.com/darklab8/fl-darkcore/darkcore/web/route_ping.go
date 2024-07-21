package web

import (
	"fmt"
	"html"
	"net/http"

	"github.com/darklab8/fl-darkcore/darkcore/core_types"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
)

const URLPing core_types.Url = "/ping"

func (w *Web) NewEndpointPing() *registry.Endpoint {
	return &registry.Endpoint{
		Url: URLPing,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				fmt.Fprintf(w, "pong at %q", html.EscapeString(r.URL.Path))
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	}
}
