package web

/*
Entrypoint for front and for dev web server?
*/

import (
	"fmt"
	"log"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkstat/builder"
	"github.com/darklab8/fl-darkstat/darkstat/web/registry"
)

type Web struct {
	filesystem *builder.Filesystem
	registry   *registry.Registion
}

func NewWeb(filesystem *builder.Filesystem) *Web {
	w := &Web{
		filesystem: filesystem,
		registry:   registry.NewRegister(),
	}

	w.registry.Register(w.NewEndpointStatic())

	w.registry.Register(w.NewEndpointPing())

	return w
}

func (w *Web) Serve() {
	w.registry.Foreach(func(e *registry.Endpoint) {
		http.HandleFunc(string(e.Url), e.Handler)
	})

	ip := "0.0.0.0"
	port := 8000
	fmt.Printf("launching web server, visit http://localhost:%d to check it!\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil); err != nil {
		log.Fatal(err)
	}
}
