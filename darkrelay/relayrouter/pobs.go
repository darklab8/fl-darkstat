package relayrouter

import (
	"sort"

	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkrelay/relayfront"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkPobs(
	data *router.AppData,
) {
	shared := data.Shared
	build := data.Build
	configs := data.Configs

	sort.Slice(configs.PoBs, func(i, j int) bool {
		if configs.PoBs[i].Name != "" && configs.PoBs[j].Name == "" {
			return true
		}
		return configs.PoBs[i].Name < configs.PoBs[j].Name
	})

	for _, base := range configs.PoBs {
		sort.Slice(base.ShopItems, func(i, j int) bool {
			if base.ShopItems[i].Name != "" && base.ShopItems[j].Name == "" {
				return true
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
			urls.Pobs,
			relayfront.PoBsT(configs, tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Pobs),
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

}
