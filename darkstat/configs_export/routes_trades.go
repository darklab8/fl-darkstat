package configs_export

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/darklab8/fl-darkstat/darkstat/cache"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
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

func newTradePathExporter(
	e *Exporter,
	Bases []*Base,
	MiningOperations []*Base,
) *TradePathExporter {
	var sell_locations_by_commodity *cache.Cached[map[CommodityKey][]*MarketGood]

	sell_locations_by_commodity = cache.NewCached(func() map[CommodityKey][]*MarketGood {
		BasesFromPobs := e.PoBsToBases(e.GetPoBs())

		var commodity_bases []*Base = []*Base{}
		commodity_bases = append(append(Bases, BasesFromPobs...), MiningOperations...)

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

func (e *TradePathExporter) GetBaseTradePathsFrom(base *Base) []*ComboTradeRoute {
	var TradeRoutes []*ComboTradeRoute

	for _, buying_good := range base.MarketGoodsPerNick {
		if buying_good.Category != "commodity" {
			continue
		}
		if !buying_good.BaseSells {
			continue
		}
		commodity_key := GetCommodityKey(buying_good.Nickname, buying_good.ShipClass)
		commodity_selling_bases := e.sell_locations_by_base.Get()[commodity_key]
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

func (e *TradePathExporter) GetBaseTradePathsTo(base *Base) []*ComboTradeRoute {
	var TradeRoutes []*ComboTradeRoute

	for _, selling_good := range base.MarketGoodsPerNick {
		if selling_good.Category != "commodity" {
			continue
		}

		commodity_key := GetCommodityKey(selling_good.Nickname, selling_good.ShipClass)
		commodity_selling_bases := e.sell_locations_by_base.Get()[commodity_key]
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
	ProfitPerTimeForKiloVolumes float64
	ProfitWeight                float64
}

const LimitBestPaths = 1500

func (e *TradePathExporter) GetBestTradeDeals(bases []*Base) []*TradeDeal {
	var trade_deals []*TradeDeal

	for index, base := range bases {
		fmt.Println("base_", index, "is processed for trade detals")
		trade_routes := e.GetBaseTradePathsFrom(base)
		for _, trade_route := range trade_routes {
			profit_per_time := trade_route.Transport.GetProffitPerTime()
			kilo_volume := math.Min(10, KiloVolumesDeliverable(trade_route.Transport.BuyingGood, trade_route.Transport.SellingGood))

			if kilo_volume < 5 {
				continue
			}
			profit_per_time_for_kilo_volumes := kilo_volume * profit_per_time
			time_s := GetTimeS(trade_route.Transport.Route.g, trade_route.Transport.BuyingGood, trade_route.Transport.SellingGood)

			var time_weight float64
			time_weight = math.Min(time_s, 600) / 600

			trade_route.Transport.GetProffitPerTime()
			trade_deals = append(trade_deals, &TradeDeal{
				ComboTradeRoute:             trade_route,
				ProfitPerTimeForKiloVolumes: profit_per_time_for_kilo_volumes,
				ProfitWeight:                profit_per_time*math.Min(10, kilo_volume)/10 + profit_per_time*time_weight,
			})
		}
		if len(trade_deals) > LimitBestPaths+500 {
			sort.Slice(trade_deals, func(i, j int) bool {
				return trade_deals[i].ProfitWeight > trade_deals[j].ProfitWeight
			})
			trade_deals = trade_deals[:LimitBestPaths]
		}
	}
	sort.Slice(trade_deals, func(i, j int) bool {
		return trade_deals[i].ProfitWeight > trade_deals[j].ProfitWeight
	})
	trade_deals = trade_deals[:LimitBestPaths]
	return trade_deals
}

type BaseBestPathTimes struct {
	TransportProfitPerTime *float64
	FrigateProfitPerTime   *float64
	FreighterProfitPerTime *float64
}

func (e *TradePathExporter) GetBaseBestPathFrom(base *Base) *BaseBestPathTimes {
	var result *BaseBestPathTimes = &BaseBestPathTimes{}
	for _, buying_good := range base.MarketGoodsPerNick {
		if buying_good.Category != "commodity" {
			continue
		}
		if !buying_good.BaseSells {
			continue
		}
		commodity_key := GetCommodityKey(buying_good.Nickname, buying_good.ShipClass)
		commodity_selling_bases := e.sell_locations_by_base.Get()[commodity_key]

		if buying_good == nil {
			continue
		}
		for _, selling_good := range commodity_selling_bases {
			if selling_good.ShipClass != nil || buying_good.ShipClass != nil {
				continue
			}

			e.GetVolumedMarketGoods(buying_good, selling_good, func(copied_buying_good, copied_selling_good *MarketGood) {
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

func (e *TradePathExporter) GetBaseBestPathTo(base *Base) *BaseBestPathTimes {
	var result *BaseBestPathTimes = &BaseBestPathTimes{}
	for _, selling_good := range base.MarketGoodsPerNick {
		if selling_good.Category != "commodity" {
			continue
		}

		commodity_key := GetCommodityKey(selling_good.Nickname, selling_good.ShipClass)
		commodity_selling_bases := e.sell_locations_by_base.Get()[commodity_key]

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
