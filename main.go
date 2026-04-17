package main

import (
	"context"
	"errors"
	"flag"
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
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp"
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/envers"
	"github.com/darklab8/fl-darkstat/darkcore/metrics"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkcore/static_server"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkmap"
	"github.com/darklab8/fl-darkstat/darkmap/export_map"
	map_urls "github.com/darklab8/fl-darkstat/darkmap/front/urls"
	"github.com/darklab8/fl-darkstat/darkmap/linker"
	map_settings "github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkrelay/relayrouter"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings"

	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/fl-darkstat/docs"
	"github.com/darklab8/fl-darkstat/helpers"
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
	relay_fs := relay_builder.BuildAll(builder.BuildToMemory, builder.NotCleanFolder, nil)
	relay_router = nil
	relay_builder = nil
	return relay_fs
}

func SetOptimalGcForWeb() {
	debug.SetGCPercent(10) // https://go.dev/doc/gc-guide#Memory_limit improve with GOGC
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
	fmt.Println("starting app with args=", os.Args[1:])

	docs.SwaggerInfo.Host = strings.ReplaceAll(settings.Env.SiteHost, "https://", "")
	docs.SwaggerInfo.Host = strings.ReplaceAll(docs.SwaggerInfo.Host, "http://", "")
	if strings.Contains(settings.Env.SiteHost, "https") {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	web_darkstat := func(ctx_close context.Context) func() {
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil)) // for pprof
		}()
		SetOptimalGcForWeb()

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
		app_data := appdata.NewAppData(ctx_span, nil)
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
		stat_fs := stat_builder.BuildAll(builder.BuildToMemory, builder.NotCleanFolder, nil)
		span.End()
		log.Printf("Elapsed start_time_stat_router_build time %s", time.Since(start_time_stat_router_build))

		_, span = traces.Tracer.Start(ctx_span, "GetRelayFs")
		app_data.Lock()
		relay_fs := GetRelayFs(relay_data)
		app_data.Unlock()
		runtime.GC()
		span.End()

		filesystems := []*builder.Filesystem{stat_fs, relay_fs}

		if settings.Env.IsExpermentalMapWithDarkstatOn {
			map_urls.Index = "map.html"
			map_settings.Env.EnableStatRoot = true
			var linked_build *builder.Builder
			linked_build = linker.NewLinker(true, func(l *linker.Linker) {
				l.Export = export_map.NewExport(ctx_span, func(e *export_map.Export) {
					e.Mapped = app_data.Configs.Mapped
				})
			}).Link(context.Background())
			map_fs := linked_build.BuildAll(builder.BuildToMemory, builder.NotCleanFolder, nil)
			filesystems = append(filesystems, map_fs)
		}

		_, web_server := darkhttp.RegisterApiRoutes(web.NewWeb(
			filesystems,
			web.WithMutexableData(app_data),
			web.WithSiteRoot(settings.Env.SiteRoot),
			web.WithAppData(app_data),
		), app_data)
		web_opts := web.WebServeOpts{
			Port: ptr.Ptr(settings.Env.WebPort),
		}
		if settings.Env.EnableUnixSockets {
			web_opts.SockAddress = web.DarkstatHttpSock
		}
		web_closer := web_server.Serve(web_opts)
		log.Printf("Elapsed web launch time %s", time.Since(start_time_total))

		metronom := metrics.NewMetronom(web_server.GetMux())
		go metronom.Run()

		return func() {
			web_closer.Close()
			fmt.Println("graceful shutdown is certainly acomplished")
		}
	}

	for _, enver := range envers.Enverants {
		enver.ValidetNoUnused()
	}

	parser := cantil.NewConsoleParser(
		[]cantil.Action{
			{
				Nickname:    "build",
				Description: "build darkstat to static assets: html, css, js files",
				Func: func(info cantil.ActionInfo) error {
					_, err := StatBuild(
						builder.BuildToFilesystem,
						builder.YesCleanFolder,
						NotIncludePoBs,
						router.YesLinkTravelRoutes,
						nil,
					)
					return err
				},
			},
			{
				Nickname:    "static_alone",
				Description: "run only static assets web server. to serve already existing data built. For dev purposes",
				Func: func(info cantil.ActionInfo) error {
					static_server.StaticServer()
					return nil
				},
			},
			{
				Nickname:    "web",
				Description: "run as standalone application that serves darkstat from memory. DOES NOT AUTOUPDATE for disco, for that use web_cron",
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
				Nickname:    "web_cron",
				Description: "run as standalone application that serves darkstat from memory with full updates periodically. Recommended for disco.",
				Func: func(info cantil.ActionInfo) error {
					ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
					SetOptimalGcForWeb()
					go func() {
						log.Println(http.ListenAndServe("0.0.0.0:6060", nil)) // for pprof
					}()
					out, err := StatBuild(
						builder.BuildToMemory,
						builder.NotCleanFolder,
						YesIncludePobs,
						router.YesLinkTravelRoutes,
						nil,
					)
					logus.Log.CheckError(err, "failed to run stat build")

					var web_server *web.Web
					var api *darkhttp.Api
					web_server = web.NewWeb(
						out.fs,
						web.WithMutexableData(out.app_data),
						web.WithSiteRoot(settings.Env.SiteRoot),
						web.WithAppData(out.app_data),
					)
					api, web_server = darkhttp.RegisterApiRoutes(web_server, out.app_data)
					web_opts := web.WebServeOpts{
						Port: ptr.Ptr(settings.Env.WebPort),
					}
					if settings.Env.EnableUnixSockets {
						web_opts.SockAddress = web.DarkstatHttpSock
					}

					go func() {
						for {
							func() {
								defer func() {
									if r := recover(); r != nil {
										fmt.Print(string(debug.Stack()))
										logus.Log.Error("discovery read update, failed to do",
											typelog.Any("r", r),
											typelog.Any("stack", debug.Stack()),
										)
									}
								}()
								time.Sleep(time.Second * time.Duration(settings.Env.RelayLoopSecs))

								out.app_data.Configs.Mapped.ReadDiscovery(
									context.Background(),
									filefind.FindConfigs(settings.Env.FreelancerFolder),
								)

								out, err = StatBuild(
									builder.BuildToMemory,
									builder.NotCleanFolder,
									YesIncludePobs,
									router.YesLinkTravelRoutes,
									out.app_data.Configs.Mapped,
								)
								logus.Log.CheckError(err, "failed to run stat build")

								mutex := web_server.AppDataMutex
								time_switch_start := time.Now()
								mutex.Lock()
								web_server.SetFS(out.fs)
								web_server.SetMutexableData(out.app_data)
								web_server.SetAppData(out.app_data)
								mutex.Unlock()
								api.SetAppData(out.app_data)
								fmt.Println("switch of web data happened in", time.Since(time_switch_start))
								runtime.GC()
							}()
						}
					}()

					metronom := metrics.NewMetronom(web_server.GetMux())
					go metronom.Run()

					web_closer := web_server.Serve(web_opts)

					defer stop()
					<-ctx.Done()

					web_closer.Close()
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
				Nickname:    "config",
				Description: "get all configs",
				Func: func(info cantil.ActionInfo) error {
					envers.PrintSettings()
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
					resp, err := client.Get(fmt.Sprintf("http://localhost:%d/ping?password=%s", settings.Env.WebPort, settings.Env.DarkcoreEnvVars.Password))
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
				Nickname:    "map",
				Description: "map group of commands. See `map help` to discovery its commands",
				Func: func(info cantil.ActionInfo) error {
					darkmap.DarkmapCliGroup(info.CmdArgs[1:])
					return nil
				},
			},
			{
				Nickname:    "helpers",
				Description: "helpers group of commands. See `helpers help` to discovery its commands",
				Func: func(info cantil.ActionInfo) error {
					helpers.HelpersCliGroup(info.CmdArgs[1:])
					return nil
				},
			},
		},
		cantil.ParserOpts{
			DefaultAction: ptr.Ptr("web"),
			Enverants:     envers.Enverants,
		},
	)

	err := parser.Run(flag.Args())
	logus.Log.CheckError(err, "failed to run parser")
}

type IncludePobsKind int8

const (
	NotIncludePoBs IncludePobsKind = iota
	YesIncludePobs
)

type StatBuildOutput struct {
	fs       []*builder.Filesystem
	app_data *appdata.AppData
}

func StatBuild(
	to_where builder.BuildToWhere,
	clean_folder_at_start builder.CleanFolderKind,
	include_pobs IncludePobsKind,
	link_travel_routes router.LinkTravelRoutesKind,
	mapped *configs_mapped.MappedConfigs,

) (StatBuildOutput, error) {
	var out StatBuildOutput

	ctx_span, span_boot := traces.Tracer.Start(context.Background(), "build")
	defer span_boot.End()
	out.app_data = appdata.NewAppData(ctx_span, mapped)
	build := router.NewRouter(out.app_data, router.WithStaticAssetsGen(), func(l *router.Router) {
		l.LinkTravelRoutesKind = link_travel_routes
	}).Link(ctx_span)

	if include_pobs == YesIncludePobs {
		relay_data := appdata.NewRelayData(out.app_data)
		relay_router := relayrouter.NewRouter(relay_data)
		relay_router.LinkPobs(relay_data, build)
	}

	if settings.Env.IsExpermentalMapWithDarkstatOn {
		map_urls.Index = "map.html"
		map_settings.Env.EnableStatRoot = true
	}

	out.fs = append(out.fs, build.BuildAll(to_where, clean_folder_at_start, nil))

	if settings.Env.IsExpermentalMapWithDarkstatOn {
		linker := linker.NewLinker(false, func(l *linker.Linker) {
			l.Export = export_map.NewExport(ctx_span, func(e *export_map.Export) {
				e.Mapped = out.app_data.Configs.Mapped
			})
		}).Link(ctx_span)
		out.fs = append(out.fs, linker.BuildAll(to_where, builder.NotCleanFolder, nil))
	}
	return out, nil
}
