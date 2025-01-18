package router

import (
	"sort"

	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkCommodities(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	sort.Slice(data.Commodities, func(i, j int) bool {
		if data.Commodities[i].Name != "" && data.Commodities[j].Name == "" {
			return true
		}
		return data.Commodities[i].Name < data.Commodities[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.Commodities,
			front.CommoditiesT(data.FilterToUsefulCommodities(data.Commodities), tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Commodities),
			front.CommoditiesT(data.Commodities, tab.ShowEmpty(true), shared),
		),
	)

	for _, base_info := range data.Commodities {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.GoodAtBaseInfoTUrl(base_info)),
				front.GoodAtBaseInfoT(base_info.Name, base_info.Bases, front.ShowAsCommodity(true), shared),
			),
		)
	}
}
