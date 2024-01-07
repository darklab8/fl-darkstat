package main

import (
	"flag"
	"fldarkstat/fldarkstat/web"
	"fmt"
)

type Action string

const (
	CommandExample Action = "command_example"
)

func main() {
	var action string
	flag.StringVar(&action, "act", "undefined", "action to run")
	flag.Parse()
	fmt.Println("act:", action)

	switch Action(action) {

	case CommandExample:
		// Other command
	default:
		// for Clientside must be run as default web
		// Client received empty data
		web.NewWeb().Serve()
	}
}
