package router

/*
Links data from exported fl-configs
into stuff rendered by fl-darkstat
Technically it is "Router"
*/

import (
	"fmt"
	"sync"
	"time"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkcore/darkcore/core_static"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/static"
	"github.com/darklab8/fl-darkstat/darkstat/front/static_front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Router struct {
	AppData *AppData
}

type RouterOpt func(l *Router)

func NewRouter(AppData *AppData, opts ...RouterOpt) *Router {
	l := &Router{AppData: AppData}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func WithAppData(AppData *AppData) RouterOpt {
	return func(l *Router) { l.AppData = AppData }
}

type AppData struct {
	Mapped  *configs_mapped.MappedConfigs
	Build   *builder.Builder
	Configs *configs_export.Exporter
	Shared  *types.SharedData

	mu sync.Mutex
}

func (a *AppData) Lock()   { a.mu.Lock() }
func (a *AppData) Unlock() { a.mu.Unlock() }

func NewBuilder(mapped *configs_mapped.MappedConfigs) *builder.Builder {
	var build *builder.Builder
	timer_building_creation := timeit.NewTimer("building creation")

	tractor_tab_name := settings.Env.TractorTabName
	if mapped.Discovery != nil {
		tractor_tab_name = "IDs"
	}
	staticPrefix := "static/"
	siteRoot := settings.Env.SiteRoot
	params := &types.GlobalParams{
		Buildpath: "",
		Theme:     types.ThemeLight,
		Themes: []string{
			siteRoot + urls.Index.ToString(),
			siteRoot + urls.DarkIndex.ToString(),
			siteRoot + urls.VanillaIndex.ToString(),
		},
		TractorTabName: tractor_tab_name,
		SiteUrl:        settings.Env.SiteUrl,
		SiteRoot:       siteRoot,
		StaticRoot:     siteRoot + staticPrefix,
		Heading:        settings.Env.AppHeading,
		Timestamp:      time.Now().UTC(),

		RelayHost: settings.Env.RelayHost,
		RelayRoot: settings.Env.RelayRoot,
	}

	static_files := []builder.StaticFile{
		builder.NewStaticFileFromCore(core_static.HtmxJS),
		builder.NewStaticFileFromCore(core_static.HtmxPreloadJS),
		builder.NewStaticFileFromCore(core_static.SortableJS),
		builder.NewStaticFileFromCore(core_static.ResetCSS),
		builder.NewStaticFileFromCore(core_static.FaviconIco),

		builder.NewStaticFileFromCore(static_front.CommonCSS),
		builder.NewStaticFileFromCore(static_front.CustomCSS),
		builder.NewStaticFileFromCore(static_front.CustomJS),
		builder.NewStaticFileFromCore(static_front.CustomJSResizer),
		builder.NewStaticFileFromCore(static_front.CustomJSFiltering),
		builder.NewStaticFileFromCore(static_front.CustomJSFilteringRoutes),
		builder.NewStaticFileFromCore(static_front.CustomJSShared),
		builder.NewStaticFileFromCore(static_front.CustomJSSharedDiscovery),
		builder.NewStaticFileFromCore(static_front.CustomJSSharedVanilla),
	}

	for _, file := range static.StaticFilesystem.Files {
		static_files = append(static_files, builder.NewStaticFileFromCore(file))
	}

	build = builder.NewBuilder(params, static_files)
	timer_building_creation.Close()
	return build
}

func NewMapped() *configs_mapped.MappedConfigs {
	var mapped *configs_mapped.MappedConfigs
	freelancer_folder := settings.Env.FreelancerFolder

	timeit.NewTimerF(func() {
		mapped = configs_mapped.NewMappedConfigs()
	}, timeit.WithMsg("MappedConfigs creation"))
	logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))
	mapped.Read(freelancer_folder)
	return mapped
}

func NewAppData() *AppData {
	mapped := NewMapped()
	configs := configs_export.NewExporter(mapped)
	build := NewBuilder(mapped)

	var data *configs_export.Exporter
	timeit.NewTimerMF("exporting data", func() { data = configs.Export(configs_export.ExportOptions{}) })

	var shared *types.SharedData = &types.SharedData{
		Mapped: mapped,
	}

	timeit.NewTimerMF("filtering to useful stuff", func() {
		if mapped.FLSR != nil {
			shared.FLSRData = types.FLSRData{
				ShowFLSR: true,
			}
		}

		if mapped.Discovery != nil {
			shared.DiscoveryData = types.DiscoveryData{
				ShowDisco:         true,
				Ids:               configs.Tractors,
				TractorsByID:      configs.TractorsByID,
				Config:            mapped.Discovery.Techcompat,
				LatestPatch:       mapped.Discovery.LatestPatch,
				OrderedTechcompat: *configs_export.NewOrderedTechCompat(configs),
			}
		}
		fmt.Println("attempting to access l.configs.Infocards")
		shared.Infocards = configs.Infocards
	})

	shared.CraftableBaseName = mapped.CraftableBaseName()

	return &AppData{
		Build:   build,
		Configs: data,
		Shared:  shared,
		Mapped:  mapped,
	}
}

func (a *AppData) Refresh() {
	updated := NewAppData()
	a.Build = updated.Build
	a.Mapped = updated.Mapped
	a.Shared = updated.Shared
	a.Configs = updated.Configs
}

func (l *Router) Link() *builder.Builder {
	shared := l.AppData.Shared
	configs := l.AppData.Configs
	build := l.AppData.Build

	defer timeit.NewTimer("link, internal measure").Close()

	timeit.NewTimerMF("linking main stuff", func() {

		l.LinkBases(build, configs, shared)
		l.LinkFactions(build, configs, shared)
		l.LinkShips(build, configs, shared)
		l.LinkGuns(build, configs, shared)
		l.LinkCommodities(build, configs, shared)
		l.LinkAmmo(build, configs, shared)
		l.LinkMines(build, configs, shared)
		l.LinkShields(build, configs, shared)
		l.LinkThrusters(build, configs, shared)
		l.LinkTractors(build, configs, shared)
		l.LinkEngines(build, configs, shared)
		l.LinkCounterMeasures(build, configs, shared)
		l.LinkScanners(build, configs, shared)

		build.RegComps(
			builder.NewComponent(
				"index_"+"docs.html",
				front.DocsEntry(types.ThemeLight, shared),
			),
			builder.NewComponent(
				tab.AllItemsUrl(urls.Docs),
				front.DocsT(tab.ShowEmpty(true), shared),
			),
			builder.NewComponent(
				urls.Docs,
				front.DocsT(tab.ShowEmpty(false), shared),
			),
			builder.NewComponent(
				urls.Index,
				front.Index(types.ThemeLight, shared),
			),
			builder.NewComponent(
				urls.DarkIndex,
				front.Index(types.ThemeDark, shared),
			),
			builder.NewComponent(
				urls.VanillaIndex,
				front.Index(types.ThemeVanilla, shared),
			),
		)
	})

	timeit.NewTimerMF("linking most of stuff", func() {
		for nickname, infocard := range configs.Infocards {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(tab.InfocardURL(nickname)),
					tab.Infocard(infocard),
				),
			)
		}
	})

	return build
}
