package main

import (
	"context"
	"errors"
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

	"github.com/darklab8/fl-darkstat/configs"
	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp"
	"github.com/darklab8/fl-darkstat/darkapis/darkrpc"
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/metrics"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkmap"
	"github.com/darklab8/fl-darkstat/darkrelay/relayrouter"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/fl-darkstat/docs"
	"github.com/darklab8/go-utils/otlp"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/cantil"
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
	Darkmap Action = "darkmap"
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
	fmt.Println("args=", os.Args[1:])

	docs.SwaggerInfo.Host = strings.ReplaceAll(settings.Env.SiteHost, "https://", "")
	docs.SwaggerInfo.Host = strings.ReplaceAll(docs.SwaggerInfo.Host, "http://", "")
	if strings.Contains(settings.Env.SiteHost, "https") {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	fmt.Println("freelancer folder=", settings.Env.FreelancerFolder, settings.Env)
	defer func() {
		if r := recover(); r != nil {
			logus.Log.Error("Program crashed in main process",
				typelog.Any("recover", r),
				typelog.Any("stack", string(debug.Stack())),
			)
			if !settings.Env.IsDevEnv {
				fmt.Println("going to sleeping")
				time.Sleep(10 * time.Second)
			}
			panic(r)
		}
	}()

	web_darkstat := func(ctx_close context.Context) func() {
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil)) // for pprof
		}()

		start_time_total := time.Now()
		ctx_span, span_boot := traces.Tracer.Start(context.Background(), "bootstrap")
		defer span_boot.End()

		if settings.Env.IsDevEnv {
			f, err := os.Create("web.pprof")
			if err != nil {
				fmt.Println(err)
				return nil

			}
			err = pprof.StartCPUProfile(f)
			logus.Log.CheckError(err, "failed to start pprof")
		}

		start_time_app_data := time.Now()
		app_data := appdata.NewAppData(ctx_span)
		log.Printf("Elapsed start_time_app_data time %s", time.Since(start_time_app_data))

		start_time_relay_data := time.Now()
		_, span := traces.Tracer.Start(ctx_span, "NewRelayData")
		relay_data := appdata.NewRelayData(app_data)
		app_data.Configs.Mapped.Clean()
		span.End()
		log.Printf("Elapsed start_time_relay_data time %s", time.Since(start_time_relay_data))

		start_time_stat_router := time.Now()
		_, span = traces.Tracer.Start(ctx_span, "NewRouter")
		stat_router := router.NewRouter(app_data)
		span.End()
		log.Printf("Elapsed start_time_stat_router time %s", time.Since(start_time_stat_router))
		start_time_stat_router_link := time.Now()
		stat_builder := stat_router.Link(ctx_span)

		log.Printf("Elapsed start_time_stat_router_link time %s", time.Since(start_time_stat_router_link))

		start_time_stat_router_build := time.Now()
		_, span = traces.Tracer.Start(ctx_span, "stat_builder.BuildAll")
		stat_fs := stat_builder.BuildAll(true, nil)
		span.End()
		log.Printf("Elapsed start_time_stat_router_build time %s", time.Since(start_time_stat_router_build))

		_, span = traces.Tracer.Start(ctx_span, "GetRelayFs")
		app_data.Lock()
		relay_fs := GetRelayFs(relay_data)
		app_data.Unlock()
		runtime.GC()
		span.End()

		web_server := darkhttp.RegisterApiRoutes(web.NewWeb(
			[]*builder.Filesystem{stat_fs, relay_fs},
			web.WithMutexableData(app_data),
			web.WithSiteRoot(settings.Env.SiteRoot),
			web.WithAppData(app_data),
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
						start_relay_refresh := time.Now()
						app_data.Lock()
						defer app_data.Unlock()

						// TODO, refactor some day to have it more elegant.
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
						log.Printf("Elapsed start_relay_refresh time %s", time.Since(start_relay_refresh))

					}()
				}
			}()
		}

		relay_server := web.NewWeb(
			[]*builder.Filesystem{relay_fs},
			web.WithMutexableData(app_data),
			web.WithSiteRoot(settings.Env.SiteRoot),
			web.WithAppData(app_data),
			web.WithCtx(ctx_close),
		)

		if settings.Env.IsDevEnv {
			pprof.StopCPUProfile()
		}
		log.Printf("Elapsed web launch time %s", time.Since(start_time_total))

		relay_closer := relay_server.Serve(web.WebServeOpts{Port: ptr.Ptr(8080)})

		rpc_server := darkrpc.NewRpcServer(darkrpc.WithSockSrv(darkrpc.DarkstatRpcSock))
		rpc_server.Serve(app_data)

		metronom := metrics.NewMetronom(web_server.GetMux())
		grpc_server := darkgrpc.NewServer(app_data, darkgrpc.DefaultServerPort, darkgrpc.WithSockAddr(darkgrpc.DarkstatGRpcSock), darkgrpc.WithProm(metronom.Reg))
		go grpc_server.Serve()
		go metronom.Run()

		return func() {
			relay_closer.Close()
			web_closer.Close()
			rpc_server.Close()
			fmt.Println("graceful shutdown is certainly acomplished")
		}
	}

	parser := cantil.NewConsoleParser(
		[]cantil.Action{
			{
				Nickname:    "build",
				Description: "build darkstat to static assets: html, css, js files",
				Func: func(info cantil.ActionInfo) error {
					ctx_span, span_boot := traces.Tracer.Start(context.Background(), "build")
					defer span_boot.End()
					app_data := appdata.NewAppData(ctx_span)
					router.NewRouter(app_data, router.WithStaticAssetsGen()).Link(ctx_span).BuildAll(false, nil)
					return nil
				},
			},
			{
				Nickname:    "web",
				Description: "run as standalone application that serves darkstat from memory",
				Func: func(info cantil.ActionInfo) error {

					ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

					otelShutdown, err := otlp.SetupOTelSDK(ctx) // Set up OpenTelemetry.
					if err != nil {
						return err
					}
					defer func() { // Handle shutdown properly so nothing leaks.
						err = errors.Join(err, otelShutdown(context.Background()))
					}()

					closer := web_darkstat(ctx)

					defer stop()
					<-ctx.Done()
					closer()
					return nil
				},
			},
			{
				Nickname:    "version",
				Description: "get darkstat version",
				Func: func(info cantil.ActionInfo) error {
					fmt.Println("version=", settings.Env.AppVersion)
					return nil
				},
			},
			{
				Nickname:    "health",
				Description: "check darkstat is healthy. Useful for container health checks",
				Func: func(info cantil.ActionInfo) error {
					tr := &http.Transport{
						MaxIdleConns:       10,
						IdleConnTimeout:    10 * time.Second,
						DisableCompression: true,
					}
					client := &http.Client{Transport: tr}
					resp, err := client.Get(fmt.Sprintf("http://localhost:8000/ping?password=%s", settings.Env.DarkcoreEnvVars.Password))
					logus.Log.CheckPanic(err, "failed to health check")
					if resp.StatusCode != 200 {
						logus.Log.Panic("status code is not 200", typelog.Any("code", resp.StatusCode))
					}
					fmt.Println("service is healthy")
					return nil
				},
			},
			{
				Nickname:    "configs",
				Description: "run config parsing to debug configs lib stuff for its data or memory profiling. For situations when unit testing is not enough.",
				Func: func(info cantil.ActionInfo) error {
					configs.CliConfigs()
					return nil
				},
			},
			{
				Nickname:    "darkmap",
				Description: "darkmap group of commands. See `darkmap help` to discovery its commands",
				Func: func(info cantil.ActionInfo) error {
					darkmap.DarkmapCliGroup(info.CmdArgs[1:])
					return nil
				},
			},
		},
		cantil.ParserOpts{
			DefaultAction: ptr.Ptr("web"),
			Enverants:     settings.Enverants,
		},
	)
	err := parser.Run(os.Args[1:])
	logus.Log.CheckError(err, "failed to run parser")
}
