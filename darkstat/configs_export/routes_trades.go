package configs_export

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sort"
	"time"

	"github.com/darklab8/fl-darkstat/darkstat/cache"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-utils/utils/ptr"
)

type TradeRoute struct {
	Route       *Route
	BuyingGood  *MarketGood
	SellingGood *MarketGood
}

func NewTradeRoute(g *GraphResults, buying_good *MarketGood, selling_good *MarketGood) *TradeRoute {
	if g == nil {
		return &TradeRoute{Route: &Route{is_disabled: true}}
	}

	route := &TradeRoute{
		Route:       NewRoute(g, buying_good.BaseNickname.ToStr(), selling_good.BaseNickname.ToStr()),
		BuyingGood:  buying_good,
		SellingGood: selling_good,
	}

	return route
}

func (t *TradeRoute) GetProffitPerTime() float64 {
	return GetProffitPerTime(t.Route.g, t.BuyingGood, t.SellingGood)
}

func GetTimeS(g *GraphResults, BuyingGood *MarketGood, SellingGood *MarketGood) float64 {
	time_ms := trades.GetTimeMs2(g.Graph, g.Time, BuyingGood.BaseNickname.ToStr(), SellingGood.BaseNickname.ToStr())
	time_s := float64(time_ms)/trades.PrecisionMultipiler + float64(trades.BaseDockingDelay)
	return time_s
}

// memory optimized version of GetProffitPerTime
func GetProffitPerTime(g *GraphResults, BuyingGood *MarketGood, SellingGood *MarketGood) float64 {
	if g == nil {
		return 0
	}
	if SellingGood.GetPriceBaseBuysFor()-BuyingGood.PriceBaseSellsFor == 0 {
		return 0
	}
	ProffitPerV := float64(SellingGood.GetPriceBaseBuysFor()-BuyingGood.PriceBaseSellsFor) / float64(SellingGood.Volume)
	time_s := GetTimeS(g, BuyingGood, SellingGood)
	return ProffitPerV / time_s
}

type ComboTradeRoute struct {
	Transport *TradeRoute
	Frigate   *TradeRoute
	Freighter *TradeRoute
}

type TradePathExporter struct {
	*Exporter
	sell_locations_by_base *cache.Cached[map[CommodityKey][]*MarketGood]
}

func NewTradePathExporter(
	e *Exporter,
	Bases []*Base,
	MiningOperations []*Base,
) *TradePathExporter {
	var sell_locations_by_commodity *cache.Cached[map[CommodityKey][]*MarketGood]

	sell_locations_by_commodity = cache.NewCached(func() map[CommodityKey][]*MarketGood {
		BasesFromPobs := e.PoBsToBases(e.GetPoBs())

		commodity_bases := append(append(Bases, BasesFromPobs...), MiningOperations...)

		sell_locations_by_commodity := make(map[CommodityKey][]*MarketGood)
		for _, base := range commodity_bases {
			for _, market_good := range base.MarketGoodsPerNick {
				commodity_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
				sell_locations_by_commodity[commodity_key] = append(sell_locations_by_commodity[commodity_key], market_good)
			}
		}
		return sell_locations_by_commodity
	}, time.Minute*2)

	tp := &TradePathExporter{
		Exporter:               e,
		sell_locations_by_base: sell_locations_by_commodity,
	}
	return tp
}

func (e *TradePathExporter) GetBaseTradePathsFiltered(TradeRoutes []*ComboTradeRoute) []*ComboTradeRoute {
	sort.Slice(TradeRoutes, func(i, j int) bool {
		return TradeRoutes[i].Transport.GetProffitPerTime() > TradeRoutes[j].Transport.GetProffitPerTime()
	})
	return TradeRoutes
}

func (e *TradePathExporter) GetVolumedMarketGoods(buying_good *MarketGood, selling_good *MarketGood, callback func(*MarketGood, *MarketGood)) {
	if commodity, ok := e.Mapped.Equip().CommoditiesMap[buying_good.Nickname]; ok {
		// then it is commodity that can be duplicated through volumes
		for _, volume_info := range commodity.Volumes {
			copied_buying_good := GetPtrStructCopy(buying_good)
			copied_buying_good.Volume = volume_info.Volume.Get()
			copied_buying_good.ShipClass = volume_info.GetShipClass()
			// copied.OriginalVolume = commodity.OriginalVolume.Volume.Get()
			copied_selling_good := GetPtrStructCopy(selling_good)
			copied_selling_good.Volume = volume_info.Volume.Get()
			copied_selling_good.ShipClass = volume_info.GetShipClass()

			callback(copied_buying_good, copied_selling_good)
		}
	}
}

func (e *TradePathExporter) GetBaseTradePathsFrom(ctx context.Context, base *Base) []*ComboTradeRoute {
	var TradeRoutes []*ComboTradeRoute

	for _, buying_good := range base.MarketGoodsPerNick {
		if buying_good.Category != "commodity" {
			continue
		}
		if !buying_good.BaseSells {
			continue
		}
		commodity_key := GetCommodityKey(buying_good.Nickname, buying_good.ShipClass)
		commodity_selling_bases := e.sell_locations_by_base.Get(ctx)[commodity_key]
		for _, selling_good_at_base := range commodity_selling_bases {
			if selling_good_at_base.ShipClass != nil || buying_good.ShipClass != nil {
				continue
			}

			e.GetVolumedMarketGoods(buying_good, selling_good_at_base, func(copied_buying_good, copied_selling_good *MarketGood) {
				trade_route := &ComboTradeRoute{
					Transport: NewTradeRoute(e.Transport, copied_buying_good, copied_selling_good),
					Frigate:   NewTradeRoute(e.Frigate, copied_buying_good, copied_selling_good),
					Freighter: NewTradeRoute(e.Freighter, copied_buying_good, copied_selling_good),
				}
				if trade_route.Transport.GetProffitPerTime() <= 0 {
					return
				}
				kilo_volumes := KiloVolumesDeliverable(buying_good, selling_good_at_base)
				if kilo_volumes < 5 {
					return
				}
				TradeRoutes = append(TradeRoutes, trade_route)
			})

			// If u need to limit to specific min distance
			// if trade_route.Transport.GetTime() < 60*10*350 {
			// 	continue
			// }

			// fmt.Println("path for", trade_route.Transport.BuyingGood.BaseNickname, trade_route.Transport.SellingGood.BaseNickname)
			// fmt.Println("trade_route.Transport.GetPaths().length", len(trade_route.Transport.GetPaths()))
		}
	}

	return TradeRoutes
}

func (e *TradePathExporter) GetBaseTradePathsTo(ctx context.Context, base *Base) []*ComboTradeRoute {
	var TradeRoutes []*ComboTradeRoute

	for _, selling_good := range base.MarketGoodsPerNick {
		if selling_good.Category != "commodity" {
			continue
		}

		commodity_key := GetCommodityKey(selling_good.Nickname, selling_good.ShipClass)
		commodity_selling_bases := e.sell_locations_by_base.Get(ctx)[commodity_key]
		for _, buying_good := range commodity_selling_bases {
			if !buying_good.BaseSells {
				continue
			}
			if selling_good.ShipClass != nil || buying_good.ShipClass != nil {
				continue
			}
			if buying_good.FactionName == "Mining Field" {
				continue
			}

			e.GetVolumedMarketGoods(buying_good, selling_good, func(copied_buying_good, copied_selling_good *MarketGood) {
				trade_route := &ComboTradeRoute{
					Transport: NewTradeRoute(e.Transport, copied_buying_good, copied_selling_good),
					Frigate:   NewTradeRoute(e.Frigate, copied_buying_good, copied_selling_good),
					Freighter: NewTradeRoute(e.Freighter, copied_buying_good, copied_selling_good),
				}
				if trade_route.Transport.GetProffitPerTime() <= 0 {
					return
				}
				kilo_volumes := KiloVolumesDeliverable(buying_good, selling_good)
				if kilo_volumes < 5 {
					return
				}
				TradeRoutes = append(TradeRoutes, trade_route)
			})

		}
	}
	return TradeRoutes
}

type TradeDeal struct {
	*ComboTradeRoute
	FreighterInfo OneWayRouteInfo
	TransportInfo OneWayRouteInfo
}

const LimitBestPaths = 2000

type BestTradeDealsOutput struct {
	OneWayDeals []*TradeDeal
	TwoWayDeals []*TwoWayDeal
}

type TwoWayDeal struct {
	Route1        *TradeDeal
	Route2        *TradeDeal
	TransportInfo RouteInfo
	FrigateInfo   RouteInfo
	FreighterInfo RouteInfo
}

type OneWayRouteInfo struct {
	KiloVolumes                 float64
	ProfitPerTimeForKiloVolumes float64
	ProfitWeight                float64
	TimeS                       float64
}

func OneWayRouteInfoF(trade_route *TradeRoute) OneWayRouteInfo {
	var result OneWayRouteInfo
	profit_per_time := trade_route.GetProffitPerTime()
	max_importance_of_kilo_volumes := float64(50)
	result.KiloVolumes = math.Min(max_importance_of_kilo_volumes, KiloVolumesDeliverable(trade_route.BuyingGood, trade_route.SellingGood))
	result.ProfitPerTimeForKiloVolumes = result.KiloVolumes * profit_per_time
	result.TimeS = GetTimeS(trade_route.Route.g, trade_route.BuyingGood, trade_route.SellingGood)
	result.ProfitWeight = profit_per_time * result.TimeS * result.KiloVolumes
	return result
}

func (e *TradePathExporter) GetBestTradeDeals(ctx context.Context, bases []*Base) BestTradeDealsOutput {
	var result BestTradeDealsOutput
	var trade_deals []*TradeDeal

	len_bases := len(bases)
	for index, base := range bases {

		if settings.Env.DarkstatDisablePobsForBestTrades && base.IsPob {
			continue
		}

		if index%100 == 0 {
			fmt.Println("base_", index, "/", len_bases, " is processed for trade detals")
		}
		trade_routes := e.GetBaseTradePathsFrom(ctx, base)
		for _, trade_route := range trade_routes {

			if settings.Env.DarkstatDisablePobsForBestTrades && trade_route.Transport.SellingGood.PoBGood != nil {
				continue
			}

			route_info := OneWayRouteInfoF(trade_route.Transport)
			if route_info.KiloVolumes < 10 {
				continue
			}
			if route_info.TimeS*trades.PrecisionMultipiler > float64(trades.INFthreshold) {
				continue
			}
			trade_deals = append(trade_deals, &TradeDeal{
				ComboTradeRoute: trade_route,
				FreighterInfo:   OneWayRouteInfoF(trade_route.Freighter),
				TransportInfo:   OneWayRouteInfoF(trade_route.Transport),
			})
		}
		if len(trade_deals) > LimitBestPaths+500 {
			sort.Slice(trade_deals, func(i, j int) bool {
				return trade_deals[i].TransportInfo.ProfitWeight > trade_deals[j].TransportInfo.ProfitWeight
			})
			trade_deals = trade_deals[:LimitBestPaths]

			sort.Slice(trade_deals, func(i, j int) bool {
				return trade_deals[i].TransportInfo.ProfitPerTimeForKiloVolumes > trade_deals[j].TransportInfo.ProfitPerTimeForKiloVolumes
			})
			trade_deals = trade_deals[:LimitBestPaths-LimitBestPaths/10]
		}

		if index%100 == 0 {
			runtime.GC()
		}
	}

	// Final sorting
	sort.Slice(trade_deals, func(i, j int) bool {
		return trade_deals[i].Transport.GetProffitPerTime() > trade_deals[j].Transport.GetProffitPerTime()
	})
	if len(trade_deals) > LimitBestPaths+500 {
		trade_deals = trade_deals[:LimitBestPaths]
	}
	runtime.GC()

	result.OneWayDeals = trade_deals

	fmt.Println("TWO WAYS: starting calculating two way best trade routes")
	for route_index1, trade_route1 := range result.OneWayDeals {

		if route_index1%100 == 0 {
			fmt.Printf("TWO WAYS: processed %d out of %d\n", route_index1, len(result.OneWayDeals))
		}
		for _, trade_route2 := range result.OneWayDeals {

			route_info := trade_route_info(trade_route1.Transport, trade_route2.Transport)
			if route_info.Route1ConnectTime > TwoWayLimitConnnectingTimeS {
				continue
			}
			if route_info.Route2ConnectTime > TwoWayLimitConnnectingTimeS {
				continue
			}

			result.TwoWayDeals = append(result.TwoWayDeals, &TwoWayDeal{
				Route1:        trade_route1,
				Route2:        trade_route2,
				TransportInfo: trade_route_info(trade_route1.Transport, trade_route2.Transport),
				FrigateInfo:   trade_route_info(trade_route1.Frigate, trade_route2.Frigate),
				FreighterInfo: trade_route_info(trade_route1.Freighter, trade_route2.Freighter),
			})

			if len(result.TwoWayDeals) > TwoWayLimitRoutes+500 {
				sort.Slice(result.TwoWayDeals, func(i, j int) bool {
					return result.TwoWayDeals[i].TransportInfo.ProfitPerTime > result.TwoWayDeals[j].TransportInfo.ProfitPerTime
				})
				result.TwoWayDeals = result.TwoWayDeals[:TwoWayLimitRoutes]

				runtime.GC()
			}

		}
	}
	sort.Slice(result.TwoWayDeals, func(i, j int) bool {
		return result.TwoWayDeals[i].TransportInfo.ProfitPerTime > result.TwoWayDeals[j].TransportInfo.ProfitPerTime
	})
	if len(result.TwoWayDeals) > TwoWayLimitRoutes {
		result.TwoWayDeals = result.TwoWayDeals[:TwoWayLimitRoutes]
	}
	fmt.Println("TWO WAYS: finished calculating two way best trade routes, found=", len(result.TwoWayDeals))

	runtime.GC()

	return result
}

var TwoWayLimitRoutes = 2000
var TwoWayLimitConnnectingTimeS = float64(180)

type RouteInfo struct {
	ProfitPerTime     float64
	Route1Time        float64
	Route2Time        float64
	Route1ConnectTime float64
	Route2ConnectTime float64
	TotalProfit       float64
	TwoWayTime        float64
}

func trade_route_info(trade_route1 *TradeRoute, trade_route2 *TradeRoute) RouteInfo {
	var route_info RouteInfo
	// time of routes
	route_info.Route1Time = GetTimeS(trade_route1.Route.g, trade_route1.BuyingGood, trade_route1.SellingGood)
	route_info.Route2Time = GetTimeS(trade_route1.Route.g, trade_route2.BuyingGood, trade_route2.SellingGood)

	// Calculate profit per second
	profit1 := trade_route1.GetProffitPerTime() * route_info.Route1Time
	profit2 := trade_route2.GetProffitPerTime() * route_info.Route2Time

	// time to connect between routes
	route_info.Route1ConnectTime = GetTimeS(trade_route1.Route.g, trade_route1.SellingGood, trade_route2.BuyingGood)
	route_info.Route2ConnectTime = GetTimeS(trade_route2.Route.g, trade_route2.SellingGood, trade_route1.BuyingGood)

	route_info.TotalProfit = (profit1 + profit2)
	route_info.TwoWayTime = (route_info.Route1Time + route_info.Route2Time + route_info.Route1ConnectTime + route_info.Route2ConnectTime)

	route_info.ProfitPerTime = route_info.TotalProfit / route_info.TwoWayTime
	return route_info
}

type BaseBestPathTimes struct {
	TransportProfitPerTime *float64
	FrigateProfitPerTime   *float64
	FreighterProfitPerTime *float64
}

func (e *TradePathExporter) GetBaseBestPathFrom(ctx context.Context, base *Base) *BaseBestPathTimes {
	var result *BaseBestPathTimes = &BaseBestPathTimes{}
	for _, buying_good := range base.MarketGoodsPerNick {
		if buying_good.Category != "commodity" {
			continue
		}
		if !buying_good.BaseSells {
			continue
		}
		commodity_key := GetCommodityKey(buying_good.Nickname, buying_good.ShipClass)
		commodity_selling_bases := e.sell_locations_by_base.Get(ctx)[commodity_key]

		if buying_good == nil {
			continue
		}
		for _, selling_good := range commodity_selling_bases {
			if selling_good.ShipClass != nil || buying_good.ShipClass != nil {
				continue
			}

			e.GetVolumedMarketGoods(buying_good, selling_good, func(copied_buying_good, copied_selling_good *MarketGood) {
				time_s := GetTimeS(e.Transport, buying_good, selling_good)
				if time_s*trades.PrecisionMultipiler > float64(trades.INFthreshold) {
					return
				}

				TransportProfitPerTime := GetProffitPerTime(e.Transport, buying_good, selling_good)
				FrigateProfitPerTime := GetProffitPerTime(e.Frigate, buying_good, selling_good)
				FreighterProfitPerTime := GetProffitPerTime(e.Freighter, buying_good, selling_good)
				if TransportProfitPerTime <= 0 {
					return
				}

				kilo_volumes := KiloVolumesDeliverable(buying_good, selling_good)
				if kilo_volumes < 5 {
					return
				}

				if result.TransportProfitPerTime == nil {
					result.TransportProfitPerTime = ptr.Ptr(TransportProfitPerTime)
				} else if TransportProfitPerTime > *result.TransportProfitPerTime {
					result.TransportProfitPerTime = ptr.Ptr(TransportProfitPerTime)
				}

				if result.FrigateProfitPerTime == nil {
					result.FrigateProfitPerTime = ptr.Ptr(FrigateProfitPerTime)
				} else if FrigateProfitPerTime > *result.FrigateProfitPerTime {
					result.FrigateProfitPerTime = ptr.Ptr(FrigateProfitPerTime)
				}

				if result.FreighterProfitPerTime == nil {
					result.FreighterProfitPerTime = ptr.Ptr(FreighterProfitPerTime)
				} else if FreighterProfitPerTime > *result.FreighterProfitPerTime {
					result.FreighterProfitPerTime = ptr.Ptr(FreighterProfitPerTime)
				}
			})

		}
	}
	return result
}

func (e *TradePathExporter) GetBaseBestPathTo(ctx context.Context, base *Base) *BaseBestPathTimes {
	var result *BaseBestPathTimes = &BaseBestPathTimes{}
	for _, selling_good := range base.MarketGoodsPerNick {
		if selling_good.Category != "commodity" {
			continue
		}

		commodity_key := GetCommodityKey(selling_good.Nickname, selling_good.ShipClass)
		commodity_selling_bases := e.sell_locations_by_base.Get(ctx)[commodity_key]

		if selling_good == nil {
			continue
		}
		for _, buying_good := range commodity_selling_bases {
			if !buying_good.BaseSells {
				continue
			}
			if selling_good.ShipClass != nil || buying_good.ShipClass != nil {
				continue
			}
			if buying_good.FactionName == "Mining Field" {
				continue
			}
			e.GetVolumedMarketGoods(buying_good, selling_good, func(copied_buying_good, copied_selling_good *MarketGood) {
				time_s := GetTimeS(e.Transport, buying_good, selling_good)
				if time_s*trades.PrecisionMultipiler > float64(trades.INFthreshold) {
					return
				}

				TransportProfitPerTime := GetProffitPerTime(e.Transport, buying_good, selling_good)
				FrigateProfitPerTime := GetProffitPerTime(e.Frigate, buying_good, selling_good)
				FreighterProfitPerTime := GetProffitPerTime(e.Freighter, buying_good, selling_good)
				if TransportProfitPerTime <= 0 {
					return
				}

				kilo_volumes := KiloVolumesDeliverable(buying_good, selling_good)
				if kilo_volumes < 5 {
					return
				}

				if result.TransportProfitPerTime == nil {
					result.TransportProfitPerTime = ptr.Ptr(TransportProfitPerTime)
				} else if TransportProfitPerTime > *result.TransportProfitPerTime {
					result.TransportProfitPerTime = ptr.Ptr(TransportProfitPerTime)
				}

				if result.FrigateProfitPerTime == nil {
					result.FrigateProfitPerTime = ptr.Ptr(FrigateProfitPerTime)
				} else if FrigateProfitPerTime > *result.FrigateProfitPerTime {
					result.FrigateProfitPerTime = ptr.Ptr(FrigateProfitPerTime)
				}

				if result.FreighterProfitPerTime == nil {
					result.FreighterProfitPerTime = ptr.Ptr(FreighterProfitPerTime)
				} else if FreighterProfitPerTime > *result.FreighterProfitPerTime {
					result.FreighterProfitPerTime = ptr.Ptr(FreighterProfitPerTime)
				}
			})
		}
	}
	return result
}
