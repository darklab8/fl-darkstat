package web

import (
	"fmt"
	"html"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkstat/settings/types"
	"github.com/darklab8/fl-darkstat/darkstat/web/registry"
)

const URLPing types.Url = "/ping"

func init() {
	registry.Registry.Register(NewEndpointPing())
}

func NewEndpointPing() *registry.Endpoint {
	return &registry.Endpoint{
		Url: URLPing,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				fmt.Fprintf(w, "pong!, %q", html.EscapeString(r.URL.Path))
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	}
}
