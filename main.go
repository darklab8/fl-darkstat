package main

import (
	"fmt"
	"os"
	"time"

	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkrelay/relayrouter"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
	"github.com/darklab8/go-utils/utils/timeit"
)

type Action string

func (a Action) ToStr() string { return string(a) }

const (
	Build   Action = "build"
	Web     Action = "web"
	Version Action = "version"
	Relay   Action = "relay"
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

	web_darkstat := func() {
		timer_web := timeit.NewTimer("total time for web web := func()")

		timer_linker := timeit.NewTimer("linking stuff linker.NewLinker().Link()")
		stat_router := router.NewRouter().Link()
		timer_linker.Close()

		timer_buildall := timeit.NewTimer("building stuff linked_build.BuildAll()")
		stat_fs := stat_router.BuildAll(true)
		timer_buildall.Close()
		timer_web.Close()
		web.NewWeb(stat_fs).Serve(web.WebServeOpts{})
	}

	switch Action(action) {

	case Build:
		router.NewRouter().Link().BuildAll(false)
	case Web:
		web_darkstat()
	case Version:
		fmt.Println("version=", settings.Env.AppVersion)
	case Relay:
		stat_router := router.NewRouter()
		stat_builder := stat_router.Link()
		stat_fs := stat_builder.BuildAll(true)
		go web.NewWeb(stat_fs).Serve(web.WebServeOpts{})

		relay_router := relayrouter.NewRouter(relayrouter.WithAppData(stat_router.AppData))
		relay_builder := relay_router.Link()
		relay_fs := relay_builder.BuildAll(true)
		web.NewWeb(relay_fs).Serve(web.WebServeOpts{Port: ptr.Ptr(8080)})
	default:
		web_darkstat()
	}
}
