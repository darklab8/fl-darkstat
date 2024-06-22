package configs_export

import (
	"fmt"
	"sync"

	"github.com/darklab8/fl-configs/configs/configs_export/trades"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-configs/configs/configs_settings"
	"github.com/darklab8/fl-configs/configs/conftypes"
)

type InfocardKey string

type Infocard []string

func (e *Exporter) exportInfocards(nickname InfocardKey, infocard_ids ...int) {
	if _, ok := e.Infocards[InfocardKey(nickname)]; ok {
		return
	}

	for _, info_id := range infocard_ids {
		if value, ok := e.configs.Infocards.Infocards[info_id]; ok {
			e.Infocards[InfocardKey(nickname)] = append(e.Infocards[InfocardKey(nickname)], value.Lines...)
		}
	}

	if len(e.Infocards[InfocardKey(nickname)]) == 0 {
		e.Infocards[InfocardKey(nickname)] = []string{"no infocard"}
	}
}

type Infocards map[InfocardKey]Infocard

type Exporter struct {
	configs *configs_mapped.MappedConfigs

	Bases                []*Base
	useful_bases_by_nick map[string]*Base

	transport_graph *trades.GameGraph
	transport_dists [][]int
	freighter_graph *trades.GameGraph
	freighter_dists [][]int

	Factions    []Faction
	Infocards   Infocards
	Commodities []*Commodity
	Guns        []Gun
	Missiles    []Gun
	Mines       []Mine
	Shields     []Shield
	Thrusters   []Thruster
	Ships       []Ship
	Tractors    []Tractor
	Engines     []Engine
	CMs         []CounterMeasure
	Scanners    []Scanner
	Ammos       []Ammo
}

type OptExport func(e *Exporter)

func NewExporter(configs *configs_mapped.MappedConfigs, opts ...OptExport) *Exporter {
	e := &Exporter{
		configs:   configs,
		Infocards: map[InfocardKey]Infocard{},
	}

	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Exporter) Export() *Exporter {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		e.transport_graph = trades.MapConfigsToFloyder(e.configs, trades.WithFreighterPaths(false))
		johnson := trades.NewDijkstraApspFromGraph(e.transport_graph)
		e.transport_dists = johnson.DijkstraApsp()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		e.freighter_graph = trades.MapConfigsToFloyder(e.configs, trades.WithFreighterPaths(true))
		johnson := trades.NewDijkstraApspFromGraph(e.freighter_graph)
		e.freighter_dists = johnson.DijkstraApsp()
		wg.Done()
	}()
	e.Bases = e.GetBases()
	useful_bases := FilterToUserfulBases(e.Bases)
	e.useful_bases_by_nick = make(map[string]*Base)
	for _, base := range useful_bases {
		e.useful_bases_by_nick[base.Nickname] = base
	}

	e.Tractors = e.GetTractors()
	e.Factions = e.GetFactions(e.Bases)
	e.Bases = e.GetMissions(e.Bases, e.Factions)
	e.Commodities = e.GetCommodities()
	e.Guns = e.GetGuns(e.Tractors)
	e.Missiles = e.GetMissiles(e.Tractors)
	e.Mines = e.GetMines(e.Tractors)
	e.Shields = e.GetShields(e.Tractors)
	e.Thrusters = e.GetThrusters(e.Tractors)
	e.Ships = e.GetShips(e.Tractors)
	e.Engines = e.GetEngines(e.Tractors)
	e.CMs = e.GetCounterMeasures(e.Tractors)
	e.Scanners = e.GetScanners(e.Tractors)
	e.Ammos = e.GetAmmo(e.Tractors)
	wg.Wait()
	e.Bases, e.Commodities = e.TradePaths(e.Bases, e.Commodities)
	return e
}

func Export(configs *configs_mapped.MappedConfigs) *Exporter {
	return NewExporter(configs).Export()
}

func Empty(phrase string) bool {
	for _, letter := range phrase {
		if letter != ' ' {
			return false
		}
	}
	return true
}

func (e *Exporter) Buyable(Bases []*GoodAtBase) bool {
	for _, base := range Bases {

		if e.useful_bases_by_nick != nil {
			if _, ok := e.useful_bases_by_nick[base.BaseNickname]; ok {
				return true
			}
		}
	}

	return false
}

func Buyable(Bases []*GoodAtBase) bool {
	return len(Bases) > 0
}

type DiscoveryTechCompat struct {
	TechcompatByID map[conftypes.TractorID]float64
	TechCell       string
}

func CalculateTechCompat(Discovery *configs_mapped.DiscoveryConfig, ids []Tractor, nickname string) *DiscoveryTechCompat {
	if Discovery == nil {
		return nil
	}

	techcompat := &DiscoveryTechCompat{
		TechcompatByID: make(map[conftypes.TractorID]float64),
	}
	techcompat.TechcompatByID[""] = Discovery.Techcompat.GetCompatibilty(nickname, "")

	for _, id := range ids {
		techcompat.TechcompatByID[id.Nickname] = Discovery.Techcompat.GetCompatibilty(nickname, id.Nickname)
	}

	if compat, ok := Discovery.Techcompat.CompatByItem[nickname]; ok {
		techcompat.TechCell = compat.TechCell
	}

	return techcompat
}

func (e *Exporter) GetInfocardName(ids_name int, nickname string) string {
	if configs_settings.Env.FallbackInfonamesToNickname {
		return fmt.Sprintf("[%s]", nickname)
	}

	if infoname, ok := e.configs.Infocards.Infonames[ids_name]; ok {
		return string(infoname)
	} else {
		return fmt.Sprintf("[%s]", nickname)
	}
}
