package main

import (
	"fmt"
	"os"
	"time"

	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkstat/linker"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/fl-data-discovery/autopatcher"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/timeit"
)

type Action string

func (a Action) ToStr() string { return string(a) }

const (
	Build   Action = "build"
	Web     Action = "web"
	Version Action = "version"
	Patch   Action = "patch"
)

func main() {
	fmt.Println("freelancer folder=", settings.Env.FreelancerFolder, settings.Env)
	defer func() {
		if r := recover(); r != nil {
			logus.Log.Error("Program crashed. Sleeping 10 seconds before exit", typelog.Any("recover", r))
			if !settings.Env.IsDevEnv {
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
		timer_web := timeit.NewTimer("total time for web web := func()")

		var linked_build *builder.Builder
		timer_linker := timeit.NewTimer("linking stuff linker.NewLinker().Link()")
		linked_build = linker.NewLinker().Link()
		timer_linker.Close()

		timer_buildall := timeit.NewTimer("building stuff linked_build.BuildAll()")
		fs = linked_build.BuildAll(true)
		timer_buildall.Close()
		timer_web.Close()
		web.NewWeb(fs).Serve()
	}

	switch Action(action) {

	case Build:
		linker.NewLinker().Link().BuildAll(false)
	case Web:
		web()
	case Version:
		fmt.Println("version=", settings.Env.AppVersion)
	case Patch:
		autopatcher.RunAutopatcher()
	default:
		web()
	}
}
