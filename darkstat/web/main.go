package web

/*
Entrypoint for front and for dev web server?
*/

import (
	"log"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkstat/backend"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/front_utils"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.

type Web struct {
}

func NewWeb() *Web {
	w := &Web{}
	w.RegisterFront()
	backend.NewBackend().RegisterBack()
	return w
}

func (w *Web) RegisterFront() {
	app.Route("/", front.NewHome())
	app.RunWhenOnBrowser()
	http.Handle("/", &app.Handler{
		Name: "Main page",
		Styles: []string{
			front_utils.GetStatisRoute("hello.css"),
		},
		Icon: app.Icon{
			// Default:    front_utils.GetStatisRoute("your_icon.png"),
			// AppleTouch: front_utils.GetStatisRoute("your_icon.png"),
			// SVG:        front_utils.GetStatisRoute("your_icon.svg"),
		},
	})
}

func (w *Web) Serve() {

	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatal(err)
	}
}
