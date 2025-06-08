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

func (l *Router) LinkCounterMeasures(
	ctx context.Context,
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	ctx, span := traces.Tracer.Start(ctx, "linker-countermeasures")
	defer span.End()
	sort.Slice(data.CMs, func(i, j int) bool {
		if data.CMs[i].Name != "" && data.CMs[j].Name == "" {
			return true
		}
		return data.CMs[i].Name < data.CMs[j].Name
	})
	build.RegComps(
		builder.NewComponent(
			urls.CounterMeasures,
			front.CounterMeasureT(data.FilterToUsefulCounterMeasures(data.CMs), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.CounterMeasures),
			front.CounterMeasureT(data.CMs, tab.ShowEmpty(true), shared),
		),
	)

	for _, cm := range data.CMs {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.CounterMeasreDetailedUrl(cm)),
				front.GoodAtBaseInfoT(cm.Name, cm.Bases, front.ShowAsCommodity(false), shared),
			),
		)
	}

}
