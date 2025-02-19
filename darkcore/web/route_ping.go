package web

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

const URLPing core_types.Url = "GET /ping"

// ShowAccount godoc
// @Summary      Check Darkstat health
// @Router       /ping [get]
func NewEndpointPing(w *Web) *registry.Endpoint {
	return &registry.Endpoint{
		Url: URLPing,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			filter_nicknames := r.URL.Query()["filter_nicknames"]
			fmt.Fprintf(w, "pong at %q", html.EscapeString(r.URL.Path))
			fmt.Println(len(filter_nicknames), filter_nicknames)

			filter_nicknames = []string{}
			filter_nicknames = strings.Split(r.URL.Query().Get("filter_nicknames"), ",")
			fmt.Println(len(filter_nicknames), filter_nicknames)
		},
	}
}
