package configs_export

import (
	"github.com/darklab8/fl-configs/configs/configs_export/trades"
)

type Trades struct {
	TradeRoutes        []*TradeRoute
	BestTransportRoute *TradeRoute
	BestFreighterRoute *TradeRoute
}

type TradeRoute struct {
	Commodity               *Commodity
	BuyingGood              *GoodAtBase
	SellingGood             *GoodAtBase
	TransportDistance       int
	FreighterDistance       int
	TransportProffitPerTime float64
	FreighterProffitPerTime float64
}

func (t TradeRoute) GetProffitPerV() float64 {
	return float64(t.SellingGood.PriceBaseBuysFor-t.BuyingGood.PriceBaseSellsFor) / float64(t.Commodity.Volume)
}

func (e *Exporter) GetTransportRouteDist(t *TradeRoute) int {
	return trades.GetDist(e.transport_graph, e.transport_dists, t.BuyingGood.BaseNickname, t.SellingGood.BaseNickname)
}

func (e *Exporter) GetFreighterRouteDist(t *TradeRoute) int {
	return trades.GetDist(e.freighter_graph, e.freighter_dists, t.BuyingGood.BaseNickname, t.SellingGood.BaseNickname)
}

const TransportSpeed = 350

func (e *Exporter) GetTrProffitPerTime(t *TradeRoute) float64 {
	return t.GetProffitPerV() / (float64(e.GetTransportRouteDist(t)) / float64(TransportSpeed))
}

func (e *Exporter) GetFrProffitPerTime(t *TradeRoute) float64 {
	return t.GetProffitPerV() / (float64(e.GetFreighterRouteDist(t)) / float64(TransportSpeed))
}

func (e *Exporter) TradePaths(
	bases []*Base,
	commodities []*Commodity,
) ([]*Base, []*Commodity) {

	var commodity_by_nick map[string]*Commodity = make(map[string]*Commodity)
	var commodity_by_good_and_base map[string]map[string]*GoodAtBase = make(map[string]map[string]*GoodAtBase)
	for _, commodity := range commodities {
		commodity_by_nick[commodity.Nickname] = commodity
		if _, ok := commodity_by_good_and_base[commodity.Nickname]; !ok {
			commodity_by_good_and_base[commodity.Nickname] = make(map[string]*GoodAtBase)
		}
		for _, good_at_base := range commodity.Bases {
			commodity_by_good_and_base[commodity.Nickname][good_at_base.BaseNickname] = good_at_base
		}
	}

	for _, base := range bases {
		for _, good := range base.MarketGoods {
			if good.Type != "commodity" {
				continue
			}

			if !good.BaseSells {
				continue
			}

			commodity := commodity_by_nick[good.Nickname]
			buying_good := commodity_by_good_and_base[good.Nickname][base.Nickname]

			if buying_good == nil {
				continue
			}

			for _, selling_good_at_base := range commodity.Bases {
				trade_route := &TradeRoute{
					BuyingGood: buying_good,
					Commodity:  commodity,
				}
				trade_route.SellingGood = selling_good_at_base

				if trade_route.GetProffitPerV() <= 0 {
					continue
				}

				trade_route.TransportDistance = e.GetTransportRouteDist(trade_route)
				trade_route.FreighterDistance = e.GetFreighterRouteDist(trade_route)
				trade_route.TransportProffitPerTime = e.GetTrProffitPerTime(trade_route)
				trade_route.FreighterProffitPerTime = e.GetFrProffitPerTime(trade_route)

				base.TradeRoutes = append(base.TradeRoutes, trade_route)
				commodity.TradeRoutes = append(commodity.TradeRoutes, trade_route)
			}
		}
		// bases[base_i]
	}

	for _, commodity := range commodities {
		for _, trade_route := range commodity.TradeRoutes {
			if commodity.BestTransportRoute == nil {
				commodity.BestTransportRoute = trade_route
				commodity.BestTransportRoute = trade_route
			} else if trade_route.TransportProffitPerTime > commodity.BestTransportRoute.TransportProffitPerTime {
				commodity.BestTransportRoute = trade_route
			}

			if commodity.BestFreighterRoute == nil {
				commodity.BestFreighterRoute = trade_route
			} else if trade_route.FreighterProffitPerTime > commodity.BestFreighterRoute.FreighterProffitPerTime {
				commodity.BestFreighterRoute = trade_route
			}
		}
	}

	for _, base := range bases {
		for _, trade_route := range base.TradeRoutes {
			if base.BestTransportRoute == nil {
				base.BestTransportRoute = trade_route
				base.BestTransportRoute = trade_route
			} else if trade_route.TransportProffitPerTime > base.BestTransportRoute.TransportProffitPerTime {
				base.BestTransportRoute = trade_route
			}

			if base.BestFreighterRoute == nil {
				base.BestFreighterRoute = trade_route
			} else if trade_route.FreighterProffitPerTime > base.BestFreighterRoute.FreighterProffitPerTime {
				base.BestFreighterRoute = trade_route
			}
		}
	}

	return bases, commodities
}
