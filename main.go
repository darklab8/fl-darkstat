package main

import (
	"fmt"
	"os"
	"time"

	"github.com/darklab8/fl-darkstat/darkstat/builder"
	"github.com/darklab8/fl-darkstat/darkstat/linker"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/web"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/goutils/utils/time_measure"
)

type Action string

func (a Action) ToStr() string { return string(a) }

const (
	Build   Action = "build"
	Web     Action = "web"
	Version Action = "version"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logus.Log.Error("Program crashed. Sleeping 10 seconds before exit", typelog.Any("recover", r))
			if !(os.Getenv("DEV") == "true") {
				fmt.Println("going to sleeping")
				time.Sleep(10 * time.Second)
			}
			panic(r)
		}
	}()

	var action string
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 1 {
		action = argsWithoutProg[0]
	}
	fmt.Println("act:", action)

	web := func() {
		var fs *builder.Filesystem
		time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
			var linked_build *builder.Builder
			time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
				linked_build = linker.NewLinker().Link()
			}, time_measure.WithMsg("linking stuff linker.NewLinker().Link()"))
			time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
				fs = linked_build.BuildAll()
			}, time_measure.WithMsg("building stuff linked_build.BuildAll()"))
		}, time_measure.WithMsg("total time for web web := func()"))
		web.NewWeb(fs).Serve()
	}

	switch Action(action) {

	case Build:
		linker.NewLinker().Link().BuildAll().RenderToLocal()
	case Web:
		web()
	case Version:
		fmt.Println("version=", settings.GetVersion())
	default:
		web()
	}
}
