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

func (l *Router) LinkExtraItems(
	ctx context.Context,
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	ctx, span := traces.Tracer.Start(ctx, "linker-extra-items")
	defer span.End()
	sort.Slice(data.ExtraItems, func(i, j int) bool {
		if data.ExtraItems[i].Category != data.ExtraItems[j].Category {
			return data.ExtraItems[i].Category < data.ExtraItems[j].Category
		}
		return data.ExtraItems[i].Name < data.ExtraItems[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.ExtraItems,
			front.ExtraItemT(data.FilterToUsefulItems(data.ExtraItems), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.ExtraItems),
			front.ExtraItemT(data.ExtraItems, tab.ShowEmpty(true), shared),
		),
	)

	for _, item := range data.ExtraItems {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.ExtraDetailedUrl(item)),
				front.GoodAtBaseInfoT(item.Name, item.Bases, front.ShowAsCommodity(false), shared),
			),
		)
	}

}
