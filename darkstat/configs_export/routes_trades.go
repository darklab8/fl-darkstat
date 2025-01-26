package configs_export

import "github.com/darklab8/fl-darkstat/configs/cfg"

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

type baseAllTradeRoutes struct {
	TradeRoutes        []*ComboTradeRoute
	BestTransportRoute *TradeRoute
	BestFrigateRoute   *TradeRoute
	BestFreighterRoute *TradeRoute
}

type ComboTradeRoute struct {
	Transport *TradeRoute
	Frigate   *TradeRoute
	Freighter *TradeRoute
}

func (e *Exporter) TradePaths(
	bases []*Base,
	commodities []*Commodity,
) ([]*Base, []*Commodity) {

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

	for _, base := range bases {
		for _, good := range base.MarketGoodsPerNick {
			if good.Category != "commodity" {
				continue
			}

			if !good.BaseSells {
				continue
			}

			commodity_key := GetCommodityKey(good.Nickname, good.ShipClass)
			commodity := commodity_by_nick[commodity_key]
			buying_good := commodity_by_good_and_base[commodity_key][base.Nickname]

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

				base.TradeRoutes = append(base.TradeRoutes, trade_route)
				commodity.TradeRoutes = append(commodity.TradeRoutes, trade_route)
			}
		}
	}

	for _, commodity := range commodities {
		for _, trade_route := range commodity.TradeRoutes {
			if commodity.BestTransportRoute == nil {
				commodity.BestTransportRoute = trade_route.Transport
			} else if trade_route.Transport.GetProffitPerTime() > commodity.BestTransportRoute.GetProffitPerTime() {
				commodity.BestTransportRoute = trade_route.Transport
			}

			if commodity.BestFreighterRoute == nil {
				commodity.BestFreighterRoute = trade_route.Freighter
			} else if trade_route.Freighter.GetProffitPerTime() > commodity.BestFreighterRoute.GetProffitPerTime() {
				commodity.BestFreighterRoute = trade_route.Freighter
			}

			if commodity.BestFrigateRoute == nil {
				commodity.BestFrigateRoute = trade_route.Frigate
			} else if trade_route.Frigate.GetProffitPerTime() > commodity.BestFrigateRoute.GetProffitPerTime() {
				commodity.BestFrigateRoute = trade_route.Frigate
			}
		}
	}

	for _, base := range bases {
		for _, trade_route := range base.TradeRoutes {
			if base.BestTransportRoute == nil {
				base.BestTransportRoute = trade_route.Transport
			} else if trade_route.Transport.GetProffitPerTime() > base.BestTransportRoute.GetProffitPerTime() {
				base.BestTransportRoute = trade_route.Transport
			}

			if base.BestFreighterRoute == nil {
				base.BestFreighterRoute = trade_route.Freighter
			} else if trade_route.Freighter.GetProffitPerTime() > base.BestFreighterRoute.GetProffitPerTime() {
				base.BestFreighterRoute = trade_route.Freighter
			}

			if base.BestFrigateRoute == nil {
				base.BestFrigateRoute = trade_route.Frigate
			} else if trade_route.Frigate.GetProffitPerTime() > base.BestFrigateRoute.GetProffitPerTime() {
				base.BestFrigateRoute = trade_route.Frigate
			}
		}
	}

	return bases, commodities
}
