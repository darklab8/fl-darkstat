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

func (l *Router) LinkScanners(
	ctx context.Context,
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	ctx, span := traces.Tracer.Start(ctx, "linker-scanners")
	defer span.End()

	sort.Slice(data.Scanners, func(i, j int) bool {
		if data.Scanners[i].Name != "" && data.Scanners[j].Name == "" {
			return true
		}
		return data.Scanners[i].Name < data.Scanners[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.Scanners,
			front.ScannersT(data.FilterToUserfulScanners(data.Scanners), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Scanners),
			front.ScannersT(data.Scanners, tab.ShowEmpty(true), shared),
		),
	)

	for _, item := range data.Scanners {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.ScannerDetailedUrl(item)),
				front.GoodAtBaseInfoT(item.Name, item.Bases, front.ShowAsCommodity(false), shared),
			),
		)
	}
}
