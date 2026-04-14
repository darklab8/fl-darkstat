package web_utils

import (
	"net/http"
	"strings"
)

func MakeStaticFileResp(requested string, resp http.ResponseWriter) {
	if strings.Contains(requested, ".css") {
		resp.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if strings.Contains(requested, ".html") {
		resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else if strings.Contains(requested, ".js") {
		resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	}
}
