package router

import (
	"context"
	"sort"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkstat/cache"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkBases(
	ctx context.Context,
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	ctx, span := traces.Tracer.Start(ctx, "linker-bases")
	defer span.End()

	sort.Slice(data.Bases, func(i, j int) bool {
		return data.Bases[i].Name < data.Bases[j].Name
	})
	sort.Slice(data.TravelBases, func(i, j int) bool {
		return data.TravelBases[i].Name < data.TravelBases[j].Name
	})
	build.RegComps(
		builder.NewComponent(
			urls.Bases,
			front.BasesT(configs_export.FilterToUserfulBases(data.Bases), front.BaseShowShops, tab.ShowEmpty(false), shared, data, front.BaseOpts{}),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Bases),
			front.BasesT(data.Bases, front.BaseShowShops, tab.ShowEmpty(true), shared, data, front.BaseOpts{}),
		),
		builder.NewComponent(
			urls.Missions,
			front.BasesT(configs_export.FilterToUserfulBases(data.Bases), front.BaseShowMissions, tab.ShowEmpty(false), shared, data, front.BaseOpts{}),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Missions),
			front.BasesT(data.Bases, front.BaseShowMissions, tab.ShowEmpty(true), shared, data, front.BaseOpts{}),
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
				front.BaseMarketGoods(base.Name, base.MarketGoodsPerNick, front.BaseShowShops, shared),
			),
		)

	}

	for _, base := range data.TradeBases {
		// bottom table info trade routes!! Important.
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseTabTradesFrom)),
				front.BaseTradesFrom(base.Name, base, front.BaseTabTradesFrom, shared, data),
			),
		)
		build.RegComps( // move to back?
			builder.NewComponent(
				utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseTabTradesTo)),
				front.BaseTradesTo(base.Name, base, front.BaseTabTradesTo, shared, data),
			),
		)
		// bottom table info for all routes. Important.
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseAllRoutes)),
				front.BaseRoutes(base.Name, base, data, front.BaseAllRoutes, shared),
			),
		)

		// All travel route infocards
		for _, combo_route := range data.GetTravelRoutes(base) {
			build.RegComps(
				builder.NewComponent( // probably move to back, at least for RAM reasons
					utils_types.FilePath(front.RouteUrl(combo_route.Transport.Route)),
					front.TradeRouteInfo2(combo_route.Transport.FromBase, combo_route.Freighter.ToBase, data, shared),
				),
			)
		}
	}

	best_trades := cache.NewCached(func() configs_export.BestTradeDealsOutput {
		best_deals := data.GetBestTradeDeals(ctx, data.TradeBases)

		for _, two_way_deal := range best_deals.TwoWayDeals {
			build.RegComps(
				builder.NewComponent( // probably move to back, at least for RAM reasons
					utils_types.FilePath(front.TwoWayDealDetailedUrl(two_way_deal)),
					front.TwoWayDealDetailed(two_way_deal, shared),
				),
			)
		}
		// recalculate also here bottom tables for two ways tab
		return best_deals
	}, time.Minute*2+time.Second*5)

	// Hackish work around to ensure bottom tables got registered.
	// theoretically speaking they should have worked without awaiting calculations
	best_trades.Get(ctx)

	build.RegComps(
		builder.NewComponent( // move to back?
			urls.TwoWayDeals,
			front.TwoWayTradesTableT(
				best_trades,
				shared,
				data,
				tab.ShowEmpty(false),
			),
		),
		builder.NewComponent( // move to back?
			tab.AllItemsUrl(urls.TwoWayDeals),
			front.TwoWayTradesTableT(
				best_trades,
				shared,
				data,
				tab.ShowEmpty(true),
			),
		),
		builder.NewComponent( // move to back?
			urls.TradeDeals,
			front.TradeDeals(
				best_trades,
				shared,
				data,
				tab.ShowEmpty(false),
			),
		),
		builder.NewComponent( // move to back?
			tab.AllItemsUrl(urls.TradeDeals),
			front.TradeDeals(
				best_trades,
				shared,
				data,
				tab.ShowEmpty(true),
			),
		),
		// top tables of ore and trades
		builder.NewComponent( // move to back?
			urls.TradesFrom,
			front.BasesT(
				configs_export.FilterToUserfulBases(data.TradeBases),
				front.BaseTabTradesFrom,
				tab.ShowEmpty(false),
				shared,
				data,
				front.BaseOpts{BasesWithTradePaths: cache.NewCached(func() []front.BaseWithTradePaths {
					return front.GetBasesWithTradePathsFrom(ctx, configs_export.FilterToUserfulBases(data.TradeBases), data)
				}, time.Minute*2+time.Second*10)},
			),
		),
		builder.NewComponent( // move to back?
			tab.AllItemsUrl(urls.TradesFrom),
			front.BasesT(
				data.TradeBases,
				front.BaseTabTradesFrom,
				tab.ShowEmpty(true),
				shared,
				data,
				front.BaseOpts{BasesWithTradePaths: cache.NewCached(func() []front.BaseWithTradePaths {
					return front.GetBasesWithTradePathsFrom(ctx, data.TradeBases, data)
				}, time.Minute*2+time.Second*15)},
			),
		),
		builder.NewComponent( // move to back?
			urls.TradesTo,
			front.BasesT(
				configs_export.FilterToUserfulBases(data.TradeBases),
				front.BaseTabTradesTo,
				tab.ShowEmpty(false),
				shared,
				data,
				front.BaseOpts{BasesWithTradePaths: cache.NewCached(func() []front.BaseWithTradePaths {
					return front.GetBasesWithTradePathsTo(ctx, configs_export.FilterToUserfulBases(data.TradeBases), data)
				}, time.Minute*2+time.Second*20)},
			),
		),
		builder.NewComponent( // move to back?
			tab.AllItemsUrl(urls.TradesTo),
			front.BasesT(
				data.TradeBases,
				front.BaseTabTradesTo,
				tab.ShowEmpty(true),
				shared,
				data,
				front.BaseOpts{BasesWithTradePaths: cache.NewCached(func() []front.BaseWithTradePaths {
					return front.GetBasesWithTradePathsTo(ctx, data.TradeBases, data)
				}, time.Minute*2+time.Second*25)},
			),
		),
		builder.NewComponent( // move to back?
			urls.Asteroids,
			front.BasesT(
				configs_export.FitlerToUsefulOres(data.MiningOperations),
				front.BaseTabOres,
				tab.ShowEmpty(false),
				shared,
				data,
				front.BaseOpts{BasesWithTradePaths: cache.NewCached(func() []front.BaseWithTradePaths {
					return front.GetBasesWithTradePathsFrom(ctx, configs_export.FitlerToUsefulOres(data.MiningOperations), data)
				}, time.Minute*2+time.Second*30)},
			),
		),
		builder.NewComponent( // move to back?
			tab.AllItemsUrl(urls.Asteroids),
			front.BasesT(
				data.MiningOperations,
				front.BaseTabOres,
				tab.ShowEmpty(true),
				shared,
				data,
				front.BaseOpts{BasesWithTradePaths: cache.NewCached(func() []front.BaseWithTradePaths {
					return front.GetBasesWithTradePathsFrom(ctx, data.MiningOperations, data)
				}, time.Minute*2+time.Second*35)},
			),
		),

		builder.NewComponent(
			urls.TravelRoutes,
			front.BasesT(configs_export.FilterToUserfulBases(data.TravelBases), front.BaseAllRoutes, tab.ShowEmpty(false), shared, data, front.BaseOpts{}),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.TravelRoutes),
			front.BasesT(data.TravelBases, front.BaseAllRoutes, tab.ShowEmpty(true), shared, data, front.BaseOpts{}),
		),
	)

	for _, base := range data.MiningOperations {
		// Ore routes bottom table. Important.
		build.RegComps( // move to back?
			builder.NewComponent(
				utils_types.FilePath(front.BaseDetailedUrl(base, front.BaseTabOres)),
				front.BaseTradesFrom(base.Name, base, front.BaseTabOres, shared, data),
			),
		)

		for _, combo_route := range data.GetBaseTradePathsFrom(ctx, base) { // infocards for mining.
			build.RegComps(
				builder.NewComponent( // probably move to back?
					utils_types.FilePath(front.RouteUrl(combo_route.Transport.Route)),
					front.TradeRouteInfo3(combo_route.Freighter.BuyingGood, combo_route.Freighter.SellingGood, data, shared),
				),
			)
		}

	}
}
