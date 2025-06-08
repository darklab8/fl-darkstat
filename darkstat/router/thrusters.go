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

func (l *Router) LinkThrusters(
	ctx context.Context,
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	ctx, span := traces.Tracer.Start(ctx, "linker-thrusters")
	defer span.End()

	sort.Slice(data.Thrusters, func(i, j int) bool {
		if data.Thrusters[i].Name != "" && data.Thrusters[j].Name == "" {
			return true
		}
		return data.Thrusters[i].Name < data.Thrusters[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.Thrusters,
			front.ThrusterT(data.FilterToUsefulThrusters(data.Thrusters), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Thrusters),
			front.ThrusterT(data.Thrusters, tab.ShowEmpty(true), shared),
		),
	)

	for _, thruster := range data.Thrusters {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.ThrusterDetailedUrl(thruster)),
				front.GoodAtBaseInfoT(thruster.Name, thruster.Bases, front.ShowAsCommodity(false), shared),
			),
		)
	}
}
