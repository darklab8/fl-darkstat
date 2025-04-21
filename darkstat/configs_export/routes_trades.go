package configs_export

import (
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

// memory optimized version of GetProffitPerTime
func GetProffitPerTime(g *GraphResults, BuyingGood *MarketGood, SellingGood *MarketGood) float64 {
	if g == nil {
		return 0
	}
	if SellingGood.GetPriceBaseBuysFor()-BuyingGood.PriceBaseSellsFor == 0 {
		return 0
	}
	ProffitPerV := float64(SellingGood.GetPriceBaseBuysFor()-BuyingGood.PriceBaseSellsFor) / float64(SellingGood.Volume)
	time_ms := trades.GetTimeMs2(g.Graph, g.Time, BuyingGood.BaseNickname.ToStr(), SellingGood.BaseNickname.ToStr())
	time_s := float64(time_ms)/trades.PrecisionMultipiler + float64(trades.BaseDockingDelay)
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
	}, time.Minute)

	tp := &TradePathExporter{
		Exporter:               e,
		sell_locations_by_base: sell_locations_by_commodity,
	}
	return tp
}

func (e *TradePathExporter) GetBaseTradePathsFiltered(base *Base) []*ComboTradeRoute {
	TradeRoutes := e.GetBaseTradePaths(base)
	sort.Slice(TradeRoutes, func(i, j int) bool {
		return TradeRoutes[i].Transport.GetProffitPerTime() > TradeRoutes[j].Transport.GetProffitPerTime()
	})
	return TradeRoutes
}

var (
	KiloVolume    float64 = 1000
	MaxKilVolumes float64 = 999
)

func KiloVolumesDeliverable(buying_good *MarketGood, selling_good *MarketGood) float64 {
	if buying_good.PoBGood == nil && selling_good.PoBGood == nil {
		return MaxKilVolumes
	}

	if buying_good.PoBGood != nil {
		if buying_good.PoBGood.Quantity <= buying_good.PoBGood.MinStock {
			return 0
		}

		return (float64(buying_good.PoBGood.Quantity-buying_good.PoBGood.MinStock) * buying_good.Volume) / KiloVolume
	}

	if selling_good.PoBGood != nil {
		if selling_good.PoBGood.Quantity >= selling_good.PoBGood.MaxStock {
			return 0
		}

		return (float64(selling_good.PoBGood.MaxStock-selling_good.PoBGood.Quantity) * selling_good.Volume) / KiloVolume
	}

	a := (float64(buying_good.PoBGood.Quantity-buying_good.PoBGood.MinStock) * buying_good.Volume) / KiloVolume
	b := (float64(selling_good.PoBGood.MaxStock-selling_good.PoBGood.Quantity) * selling_good.Volume) / KiloVolume
	return math.Min(a, b)
}

func (e *TradePathExporter) GetBaseTradePaths(base *Base) []*ComboTradeRoute {
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
			if buying_good.Nickname == selling_good_at_base.Nickname && buying_good.ShipClass != selling_good_at_base.ShipClass {
				continue
			}

			trade_route := &ComboTradeRoute{
				Transport: NewTradeRoute(e.Transport, buying_good, selling_good_at_base),
				Frigate:   NewTradeRoute(e.Frigate, buying_good, selling_good_at_base),
				Freighter: NewTradeRoute(e.Freighter, buying_good, selling_good_at_base),
			}
			if trade_route.Transport.GetProffitPerTime() <= 0 {
				continue
			}

			// If u need to limit to specific min distance
			// if trade_route.Transport.GetTime() < 60*10*350 {
			// 	continue
			// }

			// fmt.Println("path for", trade_route.Transport.BuyingGood.BaseNickname, trade_route.Transport.SellingGood.BaseNickname)
			// fmt.Println("trade_route.Transport.GetPaths().length", len(trade_route.Transport.GetPaths()))

			TradeRoutes = append(TradeRoutes, trade_route)
		}
	}

	return TradeRoutes
}

type BaseBestPathTimes struct {
	TransportProfitPerTime *float64
	FrigateProfitPerTime   *float64
	FreighterProfitPerTime *float64
}

func (e *TradePathExporter) GetBaseBestPath(base *Base) *BaseBestPathTimes {
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
		for _, selling_good_at_base := range commodity_selling_bases {
			TransportProfitPerTime := GetProffitPerTime(e.Transport, buying_good, selling_good_at_base)
			FrigateProfitPerTime := GetProffitPerTime(e.Frigate, buying_good, selling_good_at_base)
			FreighterProfitPerTime := GetProffitPerTime(e.Freighter, buying_good, selling_good_at_base)
			if TransportProfitPerTime <= 0 {
				continue
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
		}
	}
	return result
}
