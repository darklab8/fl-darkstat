package configs_export

import (
	"sort"

	"github.com/darklab8/fl-darkstat/configs/cfg"
)

type TradeRoute struct {
	Route       *Route
	Commodity   *Commodity
	BuyingGood  *MarketGood
	SellingGood *MarketGood
}

func NewTradeRoute(g *GraphResults, buying_good *MarketGood, selling_good *MarketGood, commodity *Commodity) *TradeRoute {
	if g == nil {
		return &TradeRoute{Route: &Route{is_disabled: true}}
	}

	route := &TradeRoute{
		Route:       NewRoute(g, buying_good.BaseNickname.ToStr(), selling_good.BaseNickname.ToStr()),
		BuyingGood:  buying_good,
		SellingGood: selling_good,
		Commodity:   commodity,
	}

	return route
}

func (t *TradeRoute) GetProffitPerV() float64 {
	if t.Route.is_disabled {
		return 0
	}

	if t.SellingGood.GetPriceBaseBuysFor()-t.BuyingGood.PriceBaseSellsFor == 0 {
		return 0
	}

	return float64(t.SellingGood.GetPriceBaseBuysFor()-t.BuyingGood.PriceBaseSellsFor) / float64(t.Commodity.Volume)
}

func (t *TradeRoute) GetProffitPerTime() float64 {
	return t.GetProffitPerV() / t.Route.GetTimeS()
}

type ComboTradeRoute struct {
	Transport *TradeRoute
	Frigate   *TradeRoute
	Freighter *TradeRoute
}

type BaseAllTradeRoutes struct {
	TradeRoutes        []*ComboTradeRoute
	BestTransportRoute *TradeRoute
	BestFrigateRoute   *TradeRoute
	BestFreighterRoute *TradeRoute
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

func (e *TradePathExporter) GetBaseTradePathsFiltered(base *Base) *BaseAllTradeRoutes {
	result := e.GetBaseTradePaths(base)
	sort.Slice(result.TradeRoutes, func(i, j int) bool {
		return result.TradeRoutes[i].Transport.GetProffitPerTime() > result.TradeRoutes[j].Transport.GetProffitPerTime()
	})
	return result
}

func (e *TradePathExporter) GetBaseTradePaths(base *Base) *BaseAllTradeRoutes {
	result := &BaseAllTradeRoutes{}
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
			trade_route := &ComboTradeRoute{
				Transport: NewTradeRoute(e.Transport, buying_good, selling_good_at_base, commodity),
				Frigate:   NewTradeRoute(e.Frigate, buying_good, selling_good_at_base, commodity),
				Freighter: NewTradeRoute(e.Freighter, buying_good, selling_good_at_base, commodity),
			}
			if trade_route.Transport.GetProffitPerV() <= 0 {
				continue
			}

			// If u need to limit to specific min distance
			// if trade_route.Transport.GetTime() < 60*10*350 {
			// 	continue
			// }

			// fmt.Println("path for", trade_route.Transport.BuyingGood.BaseNickname, trade_route.Transport.SellingGood.BaseNickname)
			// fmt.Println("trade_route.Transport.GetPaths().length", len(trade_route.Transport.GetPaths()))

			result.TradeRoutes = append(result.TradeRoutes, trade_route)
			// commodity.TradeRoutes = append(commodity.TradeRoutes, trade_route)
		}
	}

	for _, trade_route := range result.TradeRoutes {
		if result.BestTransportRoute == nil {
			result.BestTransportRoute = trade_route.Transport
		} else if trade_route.Transport.GetProffitPerTime() > result.BestTransportRoute.GetProffitPerTime() {
			result.BestTransportRoute = trade_route.Transport
		}

		if result.BestFreighterRoute == nil {
			result.BestFreighterRoute = trade_route.Freighter
		} else if trade_route.Freighter.GetProffitPerTime() > result.BestFreighterRoute.GetProffitPerTime() {
			result.BestFreighterRoute = trade_route.Freighter
		}

		if result.BestFrigateRoute == nil {
			result.BestFrigateRoute = trade_route.Frigate
		} else if trade_route.Frigate.GetProffitPerTime() > result.BestFrigateRoute.GetProffitPerTime() {
			result.BestFrigateRoute = trade_route.Frigate
		}
	}
	return result
}
