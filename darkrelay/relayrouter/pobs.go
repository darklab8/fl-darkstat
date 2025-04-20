package relayrouter

import (
	"sort"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkrelay/relayfront"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkPobs(
	data *appdata.AppDataRelay,
	build *builder.Builder,
) {
	shared := data.Shared
	configs := data.Configs

	sort.Slice(configs.PoBs, func(i, j int) bool {
		if configs.PoBs[i].Name != "" && configs.PoBs[j].Name == "" {
			return true
		}
		return configs.PoBs[i].Name < configs.PoBs[j].Name
	})

	for _, base := range configs.PoBs {
		sort.Slice(base.ShopItems, func(i, j int) bool {
			if base.ShopItems[i].Category != base.ShopItems[j].Category {
				return base.ShopItems[i].Category < base.ShopItems[j].Category
			}
			return base.ShopItems[i].Name < base.ShopItems[j].Name
		})
	}

	// For sanity checking test regarding data being updated
	// for index, _ := range configs.PoBs {
	// 	configs.PoBs[index].Name += fmt.Sprintf("%d", rand.IntN(100))
	// 	if configs.PoBs[index].Level == nil {
	// 		configs.PoBs[index].Level = ptr.Ptr(10)
	// 	} else {
	// 		configs.PoBs[index].Level = ptr.Ptr(*configs.PoBs[index].Level + 10)
	// 	}
	// }

	build.RegComps(
		builder.NewComponent(
			urls.PoBs,
			relayfront.PoBsT(configs, tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.PoBs),
			relayfront.PoBsT(configs, tab.ShowEmpty(true), shared),
		),
	)
	for _, pob := range configs.PoBs {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(relayfront.PoBDetailedUrl(pob)),
				relayfront.PoBShopItems(pob.Name, pob.ShopItems),
			),
		)
	}

	sort.Slice(configs.PoBGoods, func(i, j int) bool {
		if configs.PoBGoods[i].Category != configs.PoBGoods[j].Category {
			return configs.PoBGoods[i].Category < configs.PoBGoods[j].Category
		}

		return configs.PoBGoods[i].Name < configs.PoBGoods[j].Name
	})

	for _, base := range configs.PoBGoods {
		sort.Slice(base.Bases, func(i, j int) bool {
			if base.Bases[i].Base.Name != "" && base.Bases[j].Base.Name == "" {
				return true
			}
			return base.Bases[i].Base.Name < base.Bases[j].Base.Name
		})
	}
	build.RegComps(
		builder.NewComponent(
			urls.PoBGoods,
			relayfront.PoBGoodsT(configs, tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.PoBGoods),
			relayfront.PoBGoodsT(configs, tab.ShowEmpty(true), shared),
		),
	)
	for _, good := range configs.PoBGoods {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(relayfront.PoBGoodDetailedUrl(good)),
				relayfront.PoBGoodPobs(good),
			),
		)
	}
}
