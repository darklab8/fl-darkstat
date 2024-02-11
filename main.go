package main

import (
	"flag"
	"fmt"

	"github.com/darklab8/fl-darkstat/darkstat/builder"
	"github.com/darklab8/fl-darkstat/darkstat/web"
)

type Action string

const (
	Build Action = "build"
)

func main() {
	var action string
	flag.StringVar(&action, "act", "undefined", "action to run")
	flag.Parse()
	fmt.Println("act:", action)

	switch Action(action) {

	case Build:
		build := builder.NewBuilder()
		build.Build()
	default:
		// for Clientside must be run as default web
		// Client received empty data
		web.NewWeb().Serve()
	}
}
