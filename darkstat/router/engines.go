package router

import (
	"context"
	"sort"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkEngines(
	ctx context.Context,
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	ctx, span := traces.Tracer.Start(ctx, "linker-engines")
	defer span.End()

	sort.Slice(data.Engines, func(i, j int) bool {
		if data.Engines[i].Name != "" && data.Engines[j].Name == "" {
			return true
		}
		return data.Engines[i].Name < data.Engines[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.Engines,
			front.Engines(data.FilterToUsefulEngines(data.Engines), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Engines),
			front.Engines(data.Engines, tab.ShowEmpty(true), shared),
		),
	)

	for _, engine := range data.Engines {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.EngineDetailedUrl(engine)),
				front.GoodAtBaseInfoT(engine.Name, engine.Bases, front.ShowAsCommodity(false), shared),
			),
		)
	}
}
