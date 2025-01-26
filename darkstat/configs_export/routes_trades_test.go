package configs_export

import (
	"fmt"
	"sync"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
	"github.com/darklab8/go-utils/utils/ptr"
)

func TestGetTrades(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	e := NewExporter(configs)
	e.ship_speeds = trades.DiscoverySpeeds

	e.Commodities = e.GetCommodities()

	mining_bases := e.GetOres(e.Commodities)
	mining_bases_by_system := make(map[string][]trades.ExtraBase)
	for _, base := range mining_bases {
		mining_bases_by_system[base.SystemNickname] = append(mining_bases_by_system[base.SystemNickname], trades.ExtraBase{
			Pos:      base.Pos,
			Nickname: base.Nickname,
		})
	}

	var wg sync.WaitGroup
	graph_options := trades.MappingOptions{TradeRoutesDetailedTradeLane: ptr.Ptr(true)}

	wg.Add(1)
	go func() {
		e.Transport = NewGraphResults(e, e.ship_speeds.AvgTransportCruiseSpeed, trades.WithFreighterPaths(false), mining_bases_by_system, graph_options)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		e.Frigate = NewGraphResults(e, e.ship_speeds.AvgFrigateCruiseSpeed, trades.WithFreighterPaths(false), mining_bases_by_system, graph_options)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		e.Freighter = NewGraphResults(e, e.ship_speeds.AvgFreighterCruiseSpeed, trades.WithFreighterPaths(true), mining_bases_by_system, graph_options)
		wg.Done()
	}()
	e.Bases = e.GetBases()
	e.Bases = append(e.Bases, mining_bases...)
	wg.Wait()
	e.Bases, e.Commodities = e.TradePaths(e.Bases, e.Commodities)

	for _, base := range e.Bases {
		if base.Nickname != "zone_br05_gold_dummy_field" {
			continue
		}
		for _, trade_route := range base.TradeRoutes {
			trade_route.Transport.Route.GetPaths()
			trade_route.Frigate.Route.GetTimeMs()
		}
	}

	e.EnhanceBasesWithIsTransportReachable(e.Bases, e.Transport, e.Freighter)

	fmt.Println()
}
