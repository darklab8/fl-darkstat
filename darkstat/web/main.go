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
	"github.com/darklab8/go-typelog/typelog"
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
	port := 8000
	logus.Log.Info("launching listening port at",
		typelog.Int("port", port),
		typelog.String("address", fmt.Sprintf("http://localhost:%d", port)),
	)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
