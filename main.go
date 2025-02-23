package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp"
	"github.com/darklab8/fl-darkstat/darkapis/darkrpc"
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkrelay/relayrouter"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
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
	Configs Action = "configs"
)

func GetRelayFs(app_data *appdata.AppDataRelay) *builder.Filesystem {
	relay_router := relayrouter.NewRouter(app_data)
	relay_builder := relay_router.Link()
	relay_fs := relay_builder.BuildAll(true, nil)
	relay_router = nil
	relay_builder = nil
	return relay_fs
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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if settings.Env.IsCPUProfilerEnabled {
		// task profiler:cpu after that
		f, err := os.Create("prof.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

	}
	if settings.Env.IsMemProfilerEnabled {
		// task profiler:mem after that
		go func() {
			time.Sleep(time.Second * 30)
			f, _ := os.Create("mem.pprof")
			pprof.WriteHeapProfile(f)
			f.Close()
		}()
	}

	docs.SwaggerInfo.Host = strings.ReplaceAll(settings.Env.SiteHost, "https://", "")
	docs.SwaggerInfo.Host = strings.ReplaceAll(docs.SwaggerInfo.Host, "http://", "")
	if strings.Contains(settings.Env.SiteHost, "https") {
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

	web_darkstat := func() func() {
		app_data := appdata.NewAppData()
		relay_data := appdata.NewRelayData(app_data)
		app_data.Configs.Mapped.Clean()

		stat_router := router.NewRouter(app_data)
		stat_builder := stat_router.Link()

		stat_fs := stat_builder.BuildAll(true, nil)

		app_data.Lock()
		relay_fs := GetRelayFs(relay_data)
		app_data.Unlock()
		runtime.GC()

		web_server := darkhttp.RegisterApiRoutes(web.NewWeb(
			[]*builder.Filesystem{stat_fs, relay_fs},
			web.WithMutexableData(app_data),
			web.WithSiteRoot(settings.Env.SiteRoot),
		), app_data)
		web_closer := web_server.Serve(web.WebServeOpts{SockAddress: web.DarkstatHttpSock})

		if app_data.Configs.IsDiscovery {
			go func() {
				for {
					func() {
						defer func() {
							if r := recover(); r != nil {
								fmt.Println("Recovered in f, trying to update app data", r, string(debug.Stack()))
							}
						}()
						time.Sleep(time.Second * time.Duration(settings.Env.RelayLoopSecs))
						app_data.Lock()
						defer app_data.Unlock()

						// TODO minimize usage of data here.
						relay_data.Configs.Mapped.Discovery.PlayerOwnedBases.Refresh()
						relay_data.Configs.PoBs = relay_data.Configs.GetPoBs()
						relay_data.Configs.PoBGoods = relay_data.Configs.GetPoBGoods(app_data.Configs.PoBs)
						relay_fs2 := GetRelayFs(relay_data)
						for key, _ := range relay_fs.Files {
							delete(relay_fs.Files, key)
						}
						relay_fs.Files = relay_fs2.Files
						logus.Log.Info("refreshed content")
						runtime.GC()
					}()
				}
			}()
		}

		relay_server := web.NewWeb(
			[]*builder.Filesystem{relay_fs},
			web.WithMutexableData(app_data),
			web.WithSiteRoot(settings.Env.SiteRoot),
		)
		relay_closer := relay_server.Serve(web.WebServeOpts{Port: ptr.Ptr(8080)})

		rpc_server := darkrpc.NewRpcServer(darkrpc.WithSockSrv(darkrpc.DarkstatRpcSock))
		rpc_server.Serve(app_data)

		grpc_server := darkgrpc.NewServer(app_data, darkgrpc.DefaultServerPort, darkgrpc.WithSockAddr(darkgrpc.DarkstatGRpcSock))
		go grpc_server.Serve()

		return func() {
			relay_closer.Close()
			web_closer.Close()
			rpc_server.Close()
			fmt.Println("graceful shutdown is certainly acomplished")
		}
	}

	switch Action(action) {

	case Build:
		app_data := appdata.NewAppData()
		router.NewRouter(app_data).Link().BuildAll(false, nil)
	case Web:
		closer := web_darkstat()

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		<-ctx.Done()
		closer()
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
	case Configs:
		main_configs()
	default:

		closer := web_darkstat()

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		<-ctx.Done()
		closer()
		fmt.Println("trying gracefully shutting down")
	}

}
