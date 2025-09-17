package router

/*
Links data from exported fl-configs
into stuff rendered by fl-darkstat
Technically it is "Router"
*/

import (
	"context"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Router struct {
	AppData              *appdata.AppData
	is_static_assets_gen bool
}

type RouterOpt func(l *Router)

func NewRouter(AppData *appdata.AppData, opts ...RouterOpt) *Router {
	l := &Router{AppData: AppData}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func WithAppData(AppData *appdata.AppData) RouterOpt {
	return func(l *Router) { l.AppData = AppData }
}

func WithStaticAssetsGen() RouterOpt {
	return func(l *Router) { l.is_static_assets_gen = true }
}

func (l *Router) Link(ctx context.Context) *builder.Builder {
	ctx, span := traces.Tracer.Start(ctx, "linker-main")
	defer span.End()
	shared := l.AppData.Shared
	configs := l.AppData.Configs
	build := appdata.NewBuilder(configs.Mapped.Discovery != nil)

	defer timeit.NewTimer("link, internal measure").Close()

	timeit.NewTimerMF("linking main stuff", func() {
		ctx, span := traces.Tracer.Start(ctx, "linker-all")
		defer span.End()

		l.LinkBases(ctx, build, configs, shared)
		l.LinkFactions(ctx, build, configs, shared)
		l.LinkShips(ctx, build, configs, shared)
		l.LinkGuns(ctx, build, configs, shared)
		l.LinkCommodities(ctx, build, configs, shared)
		l.LinkAmmo(ctx, build, configs, shared)
		l.LinkMines(ctx, build, configs, shared)
		l.LinkShields(ctx, build, configs, shared)
		l.LinkThrusters(ctx, build, configs, shared)
		l.LinkTractors(ctx, build, configs, shared)
		l.LinkEngines(ctx, build, configs, shared)
		l.LinkCounterMeasures(ctx, build, configs, shared)
		l.LinkScanners(ctx, build, configs, shared)
		l.LinkExtraItems(ctx, build, configs, shared)

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
		configs.GetInfocardsDict(func(infocards infocarder.Infocards) {
			for nickname, infocard := range infocards {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(tab.InfocardURL(nickname)),
						tab.Infocard(infocard),
					),
				)
			}
		})
	})

	return build
}
