package configs_export

import (
	"sort"

	"github.com/darklab8/fl-darkstat/configs/cfg"
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
	commodity_by_nick          map[CommodityKey]*Commodity
	commodity_by_good_and_base map[CommodityKey]map[cfg.BaseUniNick]*MarketGood
	commodity_bases            []*Base
}

func newTradePathExporter(
	e *Exporter,
	commodities []*Commodity,
	commodity_bases []*Base,
) *TradePathExporter {
	var commodity_by_nick map[CommodityKey]*Commodity = make(map[CommodityKey]*Commodity)
	var commodity_by_good_and_base map[CommodityKey]map[cfg.BaseUniNick]*MarketGood = make(map[CommodityKey]map[cfg.BaseUniNick]*MarketGood)

	for _, commodity := range commodities {
		commodity_key := GetCommodityKey(commodity.Nickname, commodity.ShipClass)
		commodity_by_nick[commodity_key] = commodity
		if _, ok := commodity_by_good_and_base[commodity_key]; !ok {
			commodity_by_good_and_base[commodity_key] = make(map[cfg.BaseUniNick]*MarketGood)
		}
		for _, good_at_base := range commodity.Bases {
			commodity_by_good_and_base[commodity_key][good_at_base.BaseNickname] = good_at_base
		}
	}

	tp := &TradePathExporter{
		Exporter:                   e,
		commodity_by_nick:          commodity_by_nick,
		commodity_by_good_and_base: commodity_by_good_and_base,
		commodity_bases:            commodity_bases,
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


func (e *TradePathExporter) GetBaseTradePaths(base *Base) []*ComboTradeRoute {
	var TradeRoutes []*ComboTradeRoute

	for _, good := range base.MarketGoodsPerNick {
		if good.Category != "commodity" {
			continue
		}
		if !good.BaseSells {
			continue
		}
		commodity_key := GetCommodityKey(good.Nickname, good.ShipClass)
		commodity := e.commodity_by_nick[commodity_key]
		buying_good := e.commodity_by_good_and_base[commodity_key][base.Nickname]

		if buying_good == nil {
			continue
		}
		for _, selling_good_at_base := range commodity.Bases {
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
	for _, good := range base.MarketGoodsPerNick {
		if good.Category != "commodity" {
			continue
		}
		if !good.BaseSells {
			continue
		}
		commodity_key := GetCommodityKey(good.Nickname, good.ShipClass)
		commodity := e.commodity_by_nick[commodity_key]
		buying_good := e.commodity_by_good_and_base[commodity_key][base.Nickname]

		if buying_good == nil {
			continue
		}
		for _, selling_good_at_base := range commodity.Bases {
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
