package main

import (
	"flag"
	"fldarkstat/fldarkstat/builder"
	"fldarkstat/fldarkstat/web"
	"fmt"
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
