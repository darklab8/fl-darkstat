package configs_export

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/solararch_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
	"github.com/darklab8/go-utils/utils/ptr"
)

func TestGetTrades(t *testing.T) {

	ctx := context.Background()
	configs := configs_mapped.TestFixtureConfigs()
	e := NewExporter(configs)
	e.ship_speeds = trades.DiscoverySpeeds

	e.Commodities = e.GetCommodities(ctx)

	mining_bases := e.GetOres(ctx, e.Commodities)
	mining_bases_by_system := make(map[string][]trades.ExtraBase)
	for _, base := range mining_bases {
		mining_bases_by_system[base.SystemNickname] = append(mining_bases_by_system[base.SystemNickname], trades.ExtraBase{
			Pos:      base.Pos,
			Nickname: base.Nickname,
		})
	}

	var wg sync.WaitGroup
	graph_options := trades.MappingOptions{TradeRoutesDetailedTradeLane: ptr.Ptr(true)}

	var DockOpts solararch_mapped.DockableOptions

	wg.Add(1)
	go func() {
		e.Transport = NewGraphResults(e, e.ship_speeds.AvgTransportCruiseSpeed, trades.MapConfigOptions{
			WithFreighterPaths: trades.WithFreighterPaths(false),
			DockOpts:           DockOpts,
		}, mining_bases_by_system, graph_options)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		e.Frigate = NewGraphResults(e, e.ship_speeds.AvgFrigateCruiseSpeed, trades.MapConfigOptions{
			WithFreighterPaths: trades.WithFreighterPaths(false),
			DockOpts:           DockOpts,
		}, mining_bases_by_system, graph_options)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		e.Freighter = NewGraphResults(e, e.ship_speeds.AvgFreighterCruiseSpeed,
			trades.MapConfigOptions{
				WithFreighterPaths: trades.WithFreighterPaths(true),
				DockOpts:           DockOpts,
			},
			mining_bases_by_system, graph_options)
		wg.Done()
	}()

	e.Bases = e.GetBases(ctx)
	e.Bases = append(e.Bases, mining_bases...)
	if e.Mapped.Discovery != nil {
		e.Bases = append(e.Bases, e.PoBsToBases(e.GetPoBs())...)
	}
	wg.Wait()

	trade_path_exporter := NewTradePathExporter(
		e,
		e.Bases,
		[]*Base{},
	)

	time_start := time.Now()

	// for profiling only stuff.
	// go test -run TestGetTrades ./...
	// go tool pprof main.go darkstat/configs_export/best_trades.prof
	// >>> web
	f, err := os.Create("best_trades.prof")
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	logus.Log.CheckError(err, "failed to start pprof")

	_ = trade_path_exporter.GetBestTradeDeals(ctx, e.Bases)

	pprof.StopCPUProfile()

	fmt.Println("best trade deals in ", time.Now().Sub(time_start).Seconds(), " seconds")

	for _, base := range e.Bases {
		// if base.Nickname != "zone_br05_gold_dummy_field" {
		// 	continue
		// }
		for _, trade_route := range trade_path_exporter.GetBaseTradePathsFrom(ctx, base) {
			trade_route.Transport.Route.GetPaths()
			trade_route.Frigate.Route.GetTimeMs()
			KiloVolumesDeliverable(trade_route.Transport.BuyingGood, trade_route.Transport.SellingGood)
		}
		break
	}

	e.EnhanceBasesWithIsTransportReachable(e.Bases, e.Transport, e.Freighter)

	fmt.Println()
}
