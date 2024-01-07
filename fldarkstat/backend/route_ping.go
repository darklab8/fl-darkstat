package backend

import (
	"fldarkstat/fldarkstat/settings/types"
	"fmt"
	"html"
	"net/http"
)

const URLPing types.Url = "/ping"

func NewEndpointPing(app *Backend) *Endpoint {
	return &Endpoint{
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
