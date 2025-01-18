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

func (l *Router) LinkTractors(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	sort.Slice(data.Tractors, func(i, j int) bool {
		if data.Tractors[i].Name != "" && data.Tractors[j].Name == "" {
			return true
		}
		return data.Tractors[i].Name < data.Tractors[j].Name
	})
	build.RegComps(
		builder.NewComponent(
			urls.Tractors,
			front.TractorsT(data.FilterToUsefulTractors(data.Tractors), tab.ShowEmpty(false), front.TractorModShop, shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Tractors),
			front.TractorsT(data.Tractors, tab.ShowEmpty(true), front.TractorModShop, shared),
		),
		builder.NewComponent(
			urls.IDRephacks,
			front.TractorsT(data.FilterToUsefulTractors(data.Tractors), tab.ShowEmpty(false), front.TractorIDRephacks, shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.IDRephacks),
			front.TractorsT(data.Tractors, tab.ShowEmpty(true), front.TractorIDRephacks, shared),
		),
	)

	for _, tractor := range data.Tractors {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.TractorDetailedUrl(tractor, front.TractorModShop)),
				front.GoodAtBaseInfoT(tractor.Name, tractor.Bases, front.ShowAsCommodity(false), shared),
			),
			builder.NewComponent(
				utils_types.FilePath(front.TractorDetailedUrl(tractor, front.TractorIDRephacks)),
				front.IDRephacksT(tractor),
			),
		)
	}
}
