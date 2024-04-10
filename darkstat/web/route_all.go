package web

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/darklab8/fl-darkstat/darkstat/common/types"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/web/registry"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

const UrlStatic types.Url = "/"

func (w *Web) NewEndpointStatic() *registry.Endpoint {
	return &registry.Endpoint{
		Url: UrlStatic,
		Handler: func(resp http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodGet:

				// var log_files []typelog.LogType = make([]typelog.LogType, len(w.filesystem.Files))

				// i := 0
				// for path, _ := range w.filesystem.Files {
				// 	log_files[i] = typelog.String(strconv.Itoa(i), path.ToString())
				// 	i++
				// }
				// logus.Log.Info("acquired files", log_files...)

				requested := req.URL.Path[1:]
				if requested == "" {
					requested = "index.html"
				}

				requested = strings.ReplaceAll(requested, "/", PATH_SEPARATOR)
				logus.Log.Info("having get request",
					typelog.String("requested_path", requested),
					typelog.Int("files_count", len(w.filesystem.Files)),
				)

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
