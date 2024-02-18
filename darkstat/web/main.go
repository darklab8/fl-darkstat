package web

/*
Entrypoint for front and for dev web server?
*/

import (
	"fmt"
	"log"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/web/registry"
)

type Web struct {
}

func NewWeb() *Web {
	w := &Web{}
	registry.Registry.Foreach(func(e *registry.Endpoint) {
		http.HandleFunc(string(e.Url), e.Handler)
	})
	return w
}

func (w *Web) Serve() {
	ip := "0.0.0.0"
	port := 8000
	logus.Log.Info(fmt.Sprintf("launching web server, visit http://localhost:%d", port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil); err != nil {
		log.Fatal(err)
	}
}
