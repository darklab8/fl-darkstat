package router

import (
	"sort"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkBases(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	sort.Slice(data.Bases, func(i, j int) bool {
		return data.Bases[i].Name < data.Bases[j].Name
	})
	sort.Slice(data.TravelBases, func(i, j int) bool {
		return data.TravelBases[i].Name < data.TravelBases[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			urls.Bases,
			front.BasesT(configs_export.FilterToUserfulBases(data.Bases), front.BaseShowShops, tab.ShowEmpty(false), shared, data),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Bases),
			front.BasesT(data.Bases, front.BaseShowShops, tab.ShowEmpty(true), shared, data),
		),
		builder.NewComponent(
			urls.Missions,
			front.BasesT(configs_export.FilterToUserfulBases(data.Bases), front.BaseShowMissions, tab.ShowEmpty(false), shared, data),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Missions),
			front.BasesT(data.Bases, front.BaseShowMissions, tab.ShowEmpty(true), shared, data),
		),
	)

	for _, base := range data.Bases {
		if base.Missions != nil {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseShowMissions)),
					front.BaseMissions(base.Name, *base.Missions, front.BaseShowMissions),
				),
			)
		}

		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseShowShops)),
				front.BaseMarketGoods(base.Name, base.MarketGoodsPerNick, front.BaseShowShops),
			),
		)

	}

	for _, base := range data.TradeBases {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseTabTrades)),
				front.BaseTrades(base.Name, base, front.BaseTabTrades, shared, data),
			),
		)
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseAllRoutes)),
				front.BaseRoutes(base.Name, base, data, front.BaseAllRoutes, shared),
			),
		)
		for _, combo_route := range data.GetBaseTradePaths(base) {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.RouteUrl(combo_route.Transport.Route)),
					front.TradeRouteInfo3(combo_route.Freighter.BuyingGood, combo_route.Freighter.SellingGood, data, shared),
				),
			)
		}

		for _, combo_route := range data.GetTravelRoutes(base) {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.RouteUrl(combo_route.Transport.Route)),
					front.TradeRouteInfo2(combo_route.Transport.FromBase, combo_route.Freighter.ToBase, data, shared),
				),
			)
		}
	}

	build.RegComps(
		builder.NewComponent(
			urls.Trades,
			front.BasesT(configs_export.FilterToUserfulBases(data.TradeBases), front.BaseTabTrades, tab.ShowEmpty(false), shared, data),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Trades),
			front.BasesT(data.TradeBases, front.BaseTabTrades, tab.ShowEmpty(true), shared, data),
		),
		builder.NewComponent(
			urls.Asteroids,
			front.BasesT(configs_export.FitlerToUsefulOres(data.MiningOperations), front.BaseTabOres, tab.ShowEmpty(false), shared, data),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Asteroids),
			front.BasesT(data.MiningOperations, front.BaseTabOres, tab.ShowEmpty(true), shared, data),
		),
		builder.NewComponent(
			urls.TravelRoutes,
			front.BasesT(configs_export.FilterToUserfulBases(data.TravelBases), front.BaseAllRoutes, tab.ShowEmpty(false), shared, data),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.TravelRoutes),
			front.BasesT(data.TravelBases, front.BaseAllRoutes, tab.ShowEmpty(true), shared, data),
		),
	)

	for _, base := range data.MiningOperations {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseTabOres)),
				front.BaseTrades(base.Name, base, front.BaseTabOres, shared, data),
			),
		)

		for _, combo_route := range data.GetBaseTradePaths(base) {

			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.RouteUrl(combo_route.Transport.Route)),
					front.TradeRouteInfo3(combo_route.Freighter.BuyingGood, combo_route.Freighter.SellingGood, data, shared),
				),
			)
		}

	}
}
