package appdata

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/core_static"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front/static"
	"github.com/darklab8/fl-darkstat/darkstat/front/static_front"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

type AppData struct {
	Configs *configs_export.Exporter
	Shared  *types.SharedData

	mu sync.RWMutex
}

func (a *AppData) Lock()    { a.mu.Lock() }
func (a *AppData) Unlock()  { a.mu.Unlock() }
func (a *AppData) RLock()   { a.mu.RLock() }
func (a *AppData) RUnlock() { a.mu.RUnlock() }

type IsDiscovery bool

func NewBuilder(is_discovery IsDiscovery) *builder.Builder {
	var build *builder.Builder
	timer_building_creation := timeit.NewTimer("building creation")

	tractor_tab_name := settings.Env.TractorTabName
	if is_discovery {
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
		SiteHost:       settings.Env.SiteHost,
		SiteRoot:       siteRoot,
		SiteUrl:        settings.Env.SiteUrl,
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

func NewMapped(ctx context.Context) *configs_mapped.MappedConfigs {
	ctx, span := traces.Tracer.Start(ctx, "NewMapped")
	defer span.End()

	var mapped *configs_mapped.MappedConfigs
	freelancer_folder := settings.Env.FreelancerFolder

	timeit.NewTimerF(func() {
		mapped = configs_mapped.NewMappedConfigs()
	}, timeit.WithMsg("MappedConfigs creation"))
	logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))
	mapped.Read(ctx, freelancer_folder)
	return mapped
}

func NewAppData(ctx context.Context) *AppData {
	ctx, span := traces.Tracer.Start(ctx, "NewAppData")
	defer span.End()

	mapped := NewMapped(ctx)
	configs := configs_export.NewExporter(mapped)

	var data *configs_export.Exporter
	timeit.NewTimerMF("exporting data", func() { data = configs.Export(ctx, configs_export.ExportOptions{}) })

	var shared *types.SharedData = &types.SharedData{
		AverageTradeLaneSpeed: mapped.GetAvgTradeLaneSpeed(),
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
		shared.Infocarder = configs.Infocarder
	})

	shared.CraftableBaseName = mapped.CraftableBaseName()

	return &AppData{
		Configs: data,
		Shared:  shared,
	}
}

func NewRelayData(app_data *AppData) *AppDataRelay {
	return &AppDataRelay{
		Configs: app_data.Configs.ExporterRelay,
		Shared:  app_data.Shared,
		mu:      &app_data.mu,
	}
}

type AppDataRelay struct {
	Configs *configs_export.ExporterRelay
	Shared  *types.SharedData

	mu *sync.RWMutex
}
