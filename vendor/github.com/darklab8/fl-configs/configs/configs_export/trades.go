package configs_export

import (
	"github.com/darklab8/fl-configs/configs/configs_export/trades"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
)

type Trades struct {
	TradeRoutes        []*ComboTradeRoute
	BestTransportRoute *TradeRoute
	BestFrigateRoute   *TradeRoute
	BestFreighterRoute *TradeRoute
}

type TradeRoute struct {
	g           *GraphResults
	Commodity   *Commodity
	BuyingGood  *GoodAtBase
	SellingGood *GoodAtBase
	is_disabled bool
}

type ComboTradeRoute struct {
	Transport *TradeRoute
	Frigate   *TradeRoute
	Freighter *TradeRoute
}

func (c *ComboTradeRoute) GetID() string {
	if c.Transport.is_disabled {
		return ""
	}
	return c.Transport.Commodity.Nickname + c.Transport.BuyingGood.BaseNickname + c.Transport.SellingGood.BaseNickname
}

func NewTradeRoute(g *GraphResults, buying_good *GoodAtBase, selling_good *GoodAtBase, commodity *Commodity) *TradeRoute {
	if g == nil {
		return &TradeRoute{is_disabled: true}
	}

	route := &TradeRoute{
		g:           g,
		BuyingGood:  buying_good,
		SellingGood: selling_good,
		Commodity:   commodity,
	}

	return route
}

func (t *TradeRoute) GetCruiseSpeed() int {
	if t.is_disabled {
		return 0
	}
	return t.g.graph.AvgCruiseSpeed
}

func (t *TradeRoute) GetCanVisitFreighterOnlyJH() bool {
	if t.is_disabled {
		return false
	}
	return bool(t.g.graph.CanVisitFreightersOnlyJHs)
}

func (t *TradeRoute) GetProffitPerV() float64 {
	if t.is_disabled {
		return 0
	}
	return float64(t.SellingGood.PriceBaseBuysFor-t.BuyingGood.PriceBaseSellsFor) / float64(t.Commodity.Volume)
}

type PathWithNavmap struct {
	trades.DetailedPath
	SectorCoord string
}

func (t *TradeRoute) GetPaths() []PathWithNavmap {
	var results []PathWithNavmap
	paths := t.g.graph.GetPaths(t.g.parents, t.g.dists, t.BuyingGood.BaseNickname, t.SellingGood.BaseNickname)

	for _, path := range paths {
		// path.NextName // nickname of object

		augmented_path := PathWithNavmap{
			DetailedPath: path,
		}

		if jh, ok := t.g.e.configs.Systems.JumpholesByNick[path.NextName]; ok {
			pos := jh.Pos.Get()

			system_uni := t.g.e.configs.Universe_config.SystemMap[universe_mapped.SystemNickname(jh.System.Nickname)]
			augmented_path.SectorCoord = VectorToSectorCoord(system_uni, pos)
		}
		if base, ok := t.g.e.configs.Systems.BasesByBases[path.NextName]; ok {
			pos := base.Pos.Get()

			system_uni := t.g.e.configs.Universe_config.SystemMap[universe_mapped.SystemNickname(base.System.Nickname)]
			augmented_path.SectorCoord = VectorToSectorCoord(system_uni, pos)
		}

		results = append(results, augmented_path)
	}
	return results
}

func (t *TradeRoute) GetNameByIdsName(ids_name int) string {
	return string(t.g.e.configs.Infocards.Infonames[ids_name])
}

func (t *TradeRoute) GetDist() int {
	return trades.GetDist(t.g.graph, t.g.dists, t.BuyingGood.BaseNickname, t.SellingGood.BaseNickname)
}

func (t *TradeRoute) GetProffitPerTime() float64 {
	return t.GetProffitPerV() / t.GetTime()
}

func (t *TradeRoute) GetTime() float64 {
	return float64(t.GetDist())/float64(t.g.graph.AvgCruiseSpeed) + float64(trades.BaseDockingDelay)
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
				trade_route := &ComboTradeRoute{
					Transport: NewTradeRoute(e.transport, buying_good, selling_good_at_base, commodity),
					Frigate:   NewTradeRoute(e.frigate, buying_good, selling_good_at_base, commodity),
					Freighter: NewTradeRoute(e.freighter, buying_good, selling_good_at_base, commodity),
				}

				if trade_route.Transport.GetProffitPerV() <= 0 {
					continue
				}

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
