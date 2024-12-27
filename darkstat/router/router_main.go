package router

/*
Links data from exported fl-configs
into stuff rendered by fl-darkstat
Technically it is "Router"
*/

import (
	"fmt"
	"time"

	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkcore/darkcore/core_static"
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
	mapped  *configs_mapped.MappedConfigs
	configs *configs_export.Exporter
}

type RouterOpt func(l *Router)

func NewLinker(opts ...RouterOpt) *Router {
	l := &Router{}
	for _, opt := range opts {
		opt(l)
	}

	timeit.NewTimerF(func() {
		freelancer_folder := settings.Env.FreelancerFolder
		if l.configs == nil {
			l.mapped = configs_mapped.NewMappedConfigs()
			logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))
			l.mapped.Read(freelancer_folder)
			l.configs = configs_export.NewExporter(l.mapped)
		}
	}, timeit.WithMsg("MappedConfigs creation"))
	return l
}

func (l *Router) Link() *builder.Builder {
	var build *builder.Builder
	defer timeit.NewTimer("link, internal measure").Close()
	timer_building_creation := timeit.NewTimer("building creation")
	tractor_tab_name := settings.Env.TractorTabName
	if l.mapped.Discovery != nil {
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
		SiteRoot:       siteRoot,
		StaticRoot:     siteRoot + staticPrefix,
		Heading:        settings.Env.AppHeading,
		Timestamp:      time.Now().UTC(),
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

	var data *configs_export.Exporter
	timeit.NewTimerMF("exporting data", func() { data = l.configs.Export() })

	var shared *types.SharedData = &types.SharedData{
		Mapped: l.mapped,
	}

	timeit.NewTimerMF("filtering to useful stuff", func() {
		if l.mapped.FLSR != nil {
			shared.FLSRData = types.FLSRData{
				ShowFLSR: true,
			}
		}

		if l.mapped.Discovery != nil {
			shared.DiscoveryData = types.DiscoveryData{
				ShowDisco:         true,
				Ids:               l.configs.Tractors,
				TractorsByID:      l.configs.TractorsByID,
				Config:            l.mapped.Discovery.Techcompat,
				LatestPatch:       l.mapped.Discovery.LatestPatch,
				OrderedTechcompat: *configs_export.NewOrderedTechCompat(l.configs),
			}
		}
		fmt.Println("attempting to access l.configs.Infocards")
		shared.Infocards = l.configs.Infocards
	})

	shared.CraftableBaseName = l.configs.CraftableBaseName()

	timeit.NewTimerMF("linking main stuff", func() {

		l.LinkBases(build, data, shared)
		l.LinkFactions(build, data, shared)
		l.LinkShips(build, data, shared)
		l.LinkGuns(build, data, shared)
		l.LinkCommodities(build, data, shared)
		l.LinkAmmo(build, data, shared)
		l.LinkMines(build, data, shared)
		l.LinkShields(build, data, shared)
		l.LinkThrusters(build, data, shared)
		l.LinkTractors(build, data, shared)
		l.LinkEngines(build, data, shared)
		l.LinkCounterMeasures(build, data, shared)
		l.LinkScanners(build, data, shared)

		build.RegComps(
			builder.NewComponent(
				"index_"+"docs.html",
				front.DocsEntry(types.ThemeLight, shared),
			),
			builder.NewComponent(
				urls.HashesIndex,
				front.HashesEntry(types.ThemeLight, shared),
			),
			builder.NewComponent(
				tab.AllItemsUrl(urls.HashesIndex),
				front.HashesEntry(types.ThemeLight, shared),
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

			builder.NewComponent(
				urls.Hashes,
				front.HashesT(data.Hashes, tab.ShowEmpty(false), shared),
			),
			builder.NewComponent(
				tab.AllItemsUrl(urls.Hashes),
				front.HashesT(data.Hashes, tab.ShowEmpty(true), shared),
			),
		)
	})

	timeit.NewTimerMF("linking most of stuff", func() {
		for nickname, infocard := range data.Infocards {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.InfocardURL(nickname)),
					tab.Infocard(infocard),
				),
			)
		}
	})

	return build
}
