package web

import (
	"fmt"
	"html"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

const URLPing core_types.Url = "GET /ping"

// ShowAccount godoc
// @Summary      Check Darkstat health
// @Router       /ping [get]
func NewEndpointPing(w *Web) *registry.Endpoint {
	return &registry.Endpoint{
		Url: URLPing,
		Handler: func(resp http.ResponseWriter, r *http.Request) {
			if w.AppDataMutex != nil {
				w.AppDataMutex.RLock()
				defer w.AppDataMutex.RUnlock()
			}

			_, err := fmt.Fprintf(resp, "pong at %q", html.EscapeString(r.URL.Path))
			logus.Log.CheckError(err, "failed to write in fprintf in ping")

		},
	}
}
