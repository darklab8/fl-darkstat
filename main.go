package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/darkstat/linker"
	"github.com/darklab8/fl-darkstat/darkstat/web"
)

type Action string

func (a Action) ToStr() string { return string(a) }

const (
	Build Action = "build"
	Web   Action = "web"
)

func main() {
	var action string
	flag.StringVar(&action, "act", string(Web),
		fmt.Sprintln("action to run. Possible choices...",
			strings.Join([]string{Build.ToStr(), Web.ToStr()}, ", ")),
	)
	flag.Parse()
	fmt.Println("act:", action)

	switch Action(action) {

	case Build:
		linker.NewLinker().Link().BuildAll().RenderToLocal()
	case Web:
		fs := linker.NewLinker().Link().BuildAll()
		web.NewWeb(fs).Serve()
	}
}
