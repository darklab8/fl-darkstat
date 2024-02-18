package web

import (
	"fmt"
	"html"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkstat/common/types"
	"github.com/darklab8/fl-darkstat/darkstat/web/registry"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

const UrlStatic types.Url = "/"

func (w *Web) NewEndpointStatic() *registry.Endpoint {
	return &registry.Endpoint{
		Url: UrlStatic,
		Handler: func(resp http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodGet:

				for path, _ := range w.filesystem.Files {
					fmt.Println("path=", path.ToString())
				}

				requested := req.URL.Path[1:]
				if requested == "" {
					requested = "index.html"
				}
				fmt.Println("requested=", requested)

				content, ok := w.filesystem.Files[utils_types.FilePath(requested)]
				if ok {
					fmt.Fprint(resp, string(content))
				} else {
					resp.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(resp, "content is not found at %s!, %q", req.URL, html.EscapeString(requested))
				}

			default:
				http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	}
}
