package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/darklab8/fl-darkcore/darkcore/builder"
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

func GetRelayFs(app_data *router.AppData) *builder.Filesystem {
	var args []relayrouter.RouterOpt
	if app_data != nil {
		args = append(args, relayrouter.WithAppData(app_data))
	}
	relay_router := relayrouter.NewRouter(args...)
	relay_builder := relay_router.Link()
	relay_fs := relay_builder.BuildAll(true, nil)
	relay_router = nil
	relay_builder = nil
	return relay_fs
}

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
		_timer_web := timeit.NewTimer("total time for web web := func()")

		_timer_linker := timeit.NewTimer("linking stuff linker.NewLinker().Link()")
		stat_router := router.NewRouter()
		stat_builder := stat_router.Link()
		_timer_linker.Close()

		_timer_buildall := timeit.NewTimer("building stuff linked_build.BuildAll()")
		stat_fs := stat_builder.BuildAll(true, nil)
		_timer_buildall.Close()
		_timer_web.Close()
		go web.NewWeb(stat_fs).Serve(web.WebServeOpts{})

		relay_fs := GetRelayFs(stat_router.AppData)
		runtime.GC()

		go func() {
			for {
				time.Sleep(time.Second * 30)
				relay_fs2 := GetRelayFs(nil)
				relay_fs.Files = relay_fs2.Files
				runtime.GC()
			}
		}()

		web.NewWeb(relay_fs).Serve(web.WebServeOpts{Port: ptr.Ptr(8080)})
	}

	switch Action(action) {

	case Build:
		router.NewRouter().Link().BuildAll(false, nil)
	case Web:
		web_darkstat()
	case Version:
		fmt.Println("version=", settings.Env.AppVersion)
	default:
		web_darkstat()
	}
}
