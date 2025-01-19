package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkrelay/relayrouter"
	"github.com/darklab8/fl-darkstat/darkstat/api"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/fl-darkstat/docs"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type Action string

func (a Action) ToStr() string { return string(a) }

const (
	Build   Action = "build"
	Web     Action = "web"
	Version Action = "version"
	Relay   Action = "relay"
	Health  Action = "health"
)

func GetRelayFs(app_data *router.AppData) *builder.Filesystem {
	relay_router := relayrouter.NewRouter(app_data)
	relay_builder := relay_router.Link()
	relay_fs := relay_builder.BuildAll(true, nil)
	relay_router = nil
	relay_builder = nil
	return relay_fs
}

type Account struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"account name"`
}

// @title Darkstat API
// @version 1.0
// @description Darkstat API exposed info in json format.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://darklab8.github.io/blog/index.html#contacts
// @contact.email dark.dreamflyer@gmail.com

// @license.name AGPL3
// @license.url https://raw.githubusercontent.com/darklab8/fl-darkstat/refs/heads/master/LICENSE

// @BasePath /
func main() {

	// // for CPU profiling only stuff.
	// f, err := os.Create("prof.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	// for Memory profiling stuff
	// go func() {
	// 	time.Sleep(time.Second * 30)
	// 	f, _ := os.Create("mem.pprof")
	// 	pprof.WriteHeapProfile(f)
	// 	f.Close()
	// }()

	docs.SwaggerInfo.Host = strings.ReplaceAll(settings.Env.SiteUrl, "https://", "")
	docs.SwaggerInfo.Host = strings.ReplaceAll(docs.SwaggerInfo.Host, "http://", "")
	if strings.Contains(settings.Env.SiteUrl, "https") {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

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
		app_data := router.NewAppData()

		stat_router := router.NewRouter(app_data)
		stat_builder := stat_router.Link()

		stat_fs := stat_builder.BuildAll(true, nil)

		app_data.Lock()
		relay_fs := GetRelayFs(stat_router.AppData)
		app_data.Unlock()
		runtime.GC()

		go api.RegisterApiRoutes(web.NewWeb(
			[]*builder.Filesystem{stat_fs, relay_fs},
			web.WithMutexableData(app_data),
			web.WithSiteRoot(settings.Env.SiteRoot),
		), app_data).Serve(web.WebServeOpts{})

		if settings.IsRelayActive(app_data.Mapped) {
			go func() {
				for {
					time.Sleep(time.Second * time.Duration(settings.Env.RelayLoopSecs))
					app_data.Lock()
					app_data.Mapped.Discovery.PlayerOwnedBases.Refresh()
					app_data.Configs.PoBs = app_data.Configs.GetPoBs()
					app_data.Configs.PoBGoods = app_data.Configs.GetPoBGoods(app_data.Configs.PoBs)
					relay_fs2 := GetRelayFs(stat_router.AppData)
					relay_fs.Files = relay_fs2.Files
					logus.Log.Info("refreshed content")
					runtime.GC()
					app_data.Unlock()
				}
			}()
		}

		web.NewWeb(
			[]*builder.Filesystem{relay_fs},
			web.WithMutexableData(app_data),
			web.WithSiteRoot(settings.Env.SiteRoot),
		).Serve(web.WebServeOpts{Port: ptr.Ptr(8080)})
	}

	switch Action(action) {

	case Build:
		app_data := router.NewAppData()
		router.NewRouter(app_data).Link().BuildAll(false, nil)
	case Web:
		web_darkstat()
	case Version:
		fmt.Println("version=", settings.Env.AppVersion)
	case Health:
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(fmt.Sprintf("http://localhost:8000/index.html?password=%s", settings.Env.DarkcoreEnvVars.Password))
		logus.Log.CheckPanic(err, "failed to health check")
		if resp.StatusCode != 200 {
			logus.Log.Panic("status code is not 200", typelog.Any("code", resp.StatusCode))
		}
		fmt.Println("service is healthy")
	default:
		web_darkstat()
	}

}
