package configs_export

import (
	"strings"
	"sync"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
)

type InfocardKey string

type InfocardPhrase struct {
	Phrase string  `json:"phrase"  validate:"required"`
	Link   *string `json:"link"`
	Bold   bool    `json:"bold"  validate:"required"`
}

type InfocardLine struct {
	Phrases []InfocardPhrase `json:"phrases"  validate:"required"`
}

func (i InfocardLine) ToStr() string {
	var sb strings.Builder
	for _, phrase := range i.Phrases {
		sb.WriteString(phrase.Phrase)
	}
	return sb.String()
}

func NewInfocardSimpleLine(line string) InfocardLine {
	return InfocardLine{Phrases: []InfocardPhrase{{Phrase: line}}}
}

func NewInfocardBuilder() InfocardBuilder {
	return InfocardBuilder{}
}
func (i *InfocardBuilder) WriteLine(phrases ...InfocardPhrase) {
	i.Lines = append(i.Lines, InfocardLine{Phrases: phrases})
}
func (i *InfocardBuilder) WriteLineStr(phrase_strs ...string) {
	var phrases []InfocardPhrase
	for _, phrase := range phrase_strs {
		phrases = append(phrases, InfocardPhrase{Phrase: phrase})
	}
	i.Lines = append(i.Lines, InfocardLine{Phrases: phrases})
}

type InfocardBuilder struct {
	Lines Infocard
}

type Infocard []InfocardLine

func (i Infocard) StringsJoin(delimiter string) string {
	var sb strings.Builder

	for _, line := range i {
		for _, phrase := range line.Phrases {
			sb.WriteString(phrase.Phrase)
		}
		sb.WriteString(delimiter)
	}
	return sb.String()
}

func (e *Exporter) exportInfocards(nickname InfocardKey, infocard_ids ...int) {
	if _, ok := e.Infocards[InfocardKey(nickname)]; ok {
		return
	}

	for _, info_id := range infocard_ids {
		if value, ok := e.Mapped.Infocards.Infocards[info_id]; ok {
			for _, line := range value.Lines {
				e.Infocards[InfocardKey(nickname)] = append(e.Infocards[InfocardKey(nickname)], NewInfocardSimpleLine(line))
			}

		}
	}

	if len(e.Infocards[InfocardKey(nickname)]) == 0 {
		e.Infocards[InfocardKey(nickname)] = []InfocardLine{NewInfocardSimpleLine("no infocard")}
	}
}

type Infocards map[InfocardKey]Infocard

type ExporterRelay struct {
	Mapped    *configs_mapped.MappedConfigs
	hashes    HashesByCat
	PoBs      []*PoB
	PoBGoods  []*PoBGood
	Infocards Infocards
}

type Exporter struct {
	*ExporterRelay
	Mapped      *configs_mapped.MappedConfigs
	IsDiscovery bool
	Bases       []*Base
	TradeBases  []*Base
	TravelBases []*Base

	MiningOperations     []*Base
	useful_bases_by_nick map[cfg.BaseUniNick]bool

	ship_speeds trades.ShipSpeeds
	Transport   *GraphResults
	Freighter   *GraphResults
	Frigate     *GraphResults

	Factions     []Faction
	Commodities  []*Commodity
	Guns         []Gun
	Missiles     []Gun
	Mines        []Mine
	Shields      []Shield
	Thrusters    []Thruster
	Ships        []Ship
	Tractors     []*Tractor
	TractorsByID map[cfg.TractorID]*Tractor
	Cloaks       []Cloak
	Engines      []Engine
	CMs          []CounterMeasure
	Scanners     []Scanner
	Ammos        []Ammo

	findable_in_loot_cache map[string]bool
	craftable_cached       map[string]bool
	pob_buyable_cache      map[string][]*PobShopItem
}

type OptExport func(e *Exporter)

func NewExporter(mapped *configs_mapped.MappedConfigs, opts ...OptExport) *Exporter {
	e := &Exporter{
		Mapped:      mapped,
		ship_speeds: trades.VanillaSpeeds,
		ExporterRelay: &ExporterRelay{
			Infocards: map[InfocardKey]Infocard{},
			Mapped:    mapped,
			hashes:    NewHashesCategories(mapped),
		},
	}

	for _, opt := range opts {
		opt(e)
	}
	return e
}

type GraphResults struct {
	e       *ExporterRelay
	Graph   *trades.GameGraph
	Time    [][]trades.Intg
	Parents [][]trades.Parent
}

func NewGraphResults(
	e *Exporter,
	avgCruiserSpeed int,
	can_visit_freighter_only_jhs trades.WithFreighterPaths,
	mining_bases_by_system map[string][]trades.ExtraBase,
	graph_options trades.MappingOptions,
) *GraphResults {
	logus.Log.Info("mapping configs to graph")
	graph := trades.MapConfigsToFGraph(
		e.Mapped,
		avgCruiserSpeed,
		can_visit_freighter_only_jhs,
		mining_bases_by_system,
		graph_options,
	)
	logus.Log.Info("new dijkstra apsp from graph")
	dijkstra_apsp := trades.NewDijkstraApspFromGraph(graph)
	logus.Log.Info("calculating dijkstra")
	dists, parents := dijkstra_apsp.DijkstraApsp()

	graph.WipeMatrix()
	return &GraphResults{
		e:       e.ExporterRelay,
		Graph:   graph,
		Time:    dists,
		Parents: parents,
	}
}

type ExportOptions struct {
	trades.MappingOptions
}

func (e *Exporter) Export(options ExportOptions) *Exporter {
	var wg sync.WaitGroup

	if e.Mapped.Discovery != nil {
		e.IsDiscovery = true
	}

	logus.Log.Info("getting bases")
	e.Bases = e.GetBases()
	useful_bases := FilterToUserfulBases(e.Bases)
	e.useful_bases_by_nick = make(map[cfg.BaseUniNick]bool)
	for _, base := range useful_bases {
		e.useful_bases_by_nick[base.Nickname] = true
	}
	e.useful_bases_by_nick[pob_crafts_nickname] = true
	e.useful_bases_by_nick[BaseLootableNickname] = true

	e.Commodities = e.GetCommodities()
	EnhanceBasesWithServerOverrides(e.Bases, e.Commodities)

	e.MiningOperations = e.GetOres(e.Commodities)
	if e.Mapped.Discovery != nil {
		e.PoBs = e.GetPoBs()
		e.PoBGoods = e.GetPoBGoods(e.PoBs)
	}

	extra_graph_bases := make(map[string][]trades.ExtraBase)
	for _, base := range e.MiningOperations {
		extra_graph_bases[base.SystemNickname] = append(extra_graph_bases[base.SystemNickname], trades.ExtraBase{
			Pos:      base.Pos,
			Nickname: base.Nickname,
		})
	}
	for _, base := range e.PoBs {
		if base.SystemNick == nil || base.Pos == nil {
			continue
		}
		extra_graph_bases[*base.SystemNick] = append(extra_graph_bases[*base.SystemNick], trades.ExtraBase{
			Pos:      *StrPosToVectorPos(*base.Pos),
			Nickname: cfg.BaseUniNick(base.Nickname),
		})
	}
	if e.Mapped.Discovery != nil {
		e.ship_speeds = trades.DiscoverySpeeds
	}

	if e.Mapped.FLSR != nil {
		e.ship_speeds = trades.FLSRSpeeds
	}

	{
		wg.Add(1)
		go func() {
			logus.Log.Info("graph launching for tranposrt")
			e.Transport = NewGraphResults(e, e.ship_speeds.AvgTransportCruiseSpeed, trades.WithFreighterPaths(false), extra_graph_bases, options.MappingOptions)
			// e.Freighter = e.Transport
			// e.Frigate = e.Transport
			wg.Done()
			logus.Log.Info("graph finished for tranposrt")
		}()
		wg.Add(1)
		go func() {
			e.Freighter = NewGraphResults(e, e.ship_speeds.AvgFreighterCruiseSpeed, trades.WithFreighterPaths(true), extra_graph_bases, options.MappingOptions)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			e.Frigate = NewGraphResults(e, e.ship_speeds.AvgFrigateCruiseSpeed, trades.WithFreighterPaths(false), extra_graph_bases, options.MappingOptions)
			wg.Done()
		}()
	}

	logus.Log.Info("getting get tractors")

	e.Tractors = e.GetTractors()
	e.TractorsByID = make(map[cfg.TractorID]*Tractor)
	for _, tractor := range e.Tractors {
		e.TractorsByID[tractor.Nickname] = tractor
	}
	e.Factions = e.GetFactions(e.Bases)
	e.Bases = e.GetMissions(e.Bases, e.Factions)

	logus.Log.Info("getting shields")

	e.Shields = e.GetShields(e.Tractors)
	buyable_shield_tech := e.GetBuyableShields(e.Shields)
	e.Guns = e.GetGuns(e.Tractors, buyable_shield_tech)
	e.Missiles = e.GetMissiles(e.Tractors, buyable_shield_tech)
	e.Mines = e.GetMines(e.Tractors)
	e.Thrusters = e.GetThrusters(e.Tractors)
	logus.Log.Info("getting ships")
	e.Ships = e.GetShips(e.Tractors, e.TractorsByID, e.Thrusters)
	e.Engines = e.GetEngines(e.Tractors)
	e.Cloaks = e.GetCloaks(e.Tractors)
	e.CMs = e.GetCounterMeasures(e.Tractors)
	e.Scanners = e.GetScanners(e.Tractors)
	logus.Log.Info("getting ammo")

	e.Ammos = e.GetAmmo(e.Tractors)
	logus.Log.Info("waiting for graph to finish")

	wg.Wait()

	logus.Log.Info("getting pob to bases")
	// TODO refactor. all my ram allocated problems are majorly here.
	BasesFromPobs := e.PoBsToBases(e.PoBs)
	TradeBases := append(e.Bases, BasesFromPobs...)
	e.TradeBases, e.Commodities = e.TradePaths(TradeBases, e.Commodities)
	e.MiningOperations, e.Commodities = e.TradePaths(e.MiningOperations, e.Commodities)
	e.TravelBases = e.TradeBases

	e.EnhanceBasesWithIsTransportReachable(e.Bases, e.Transport, e.Freighter)
	e.Bases = e.EnhanceBasesWithPobCrafts(e.Bases)
	e.Bases = e.EnhanceBasesWithLoot(e.Bases)
	logus.Log.Info("finished exporting")

	return e
}

func (e *Exporter) EnhanceBasesWithIsTransportReachable(
	bases []*Base,
	transports_graph *GraphResults,
	frighter_graph *GraphResults,
) {
	reachable_base_example := "li01_01_base"
	tg := transports_graph
	fg := frighter_graph

	for _, base := range bases {
		base_nickname := base.Nickname.ToStr()
		if trades.GetTimeMs2(tg.Graph, tg.Time, reachable_base_example, base_nickname) >= trades.INFthreshold {
			base.IsTransportUnreachable = true
		}
		if trades.GetTimeMs2(fg.Graph, fg.Time, reachable_base_example, base_nickname) < trades.INFthreshold {
			base.Reachable = true
		}
	}

	enhance_with_transport_unrechability := func(Bases map[cfg.BaseUniNick]*MarketGood) {
		for _, base := range Bases {
			if trades.GetTimeMs2(tg.Graph, tg.Time, reachable_base_example, string(base.BaseNickname)) >= trades.INF/2 {
				base.IsTransportUnreachable = true
			}
		}
	}

	for _, item := range e.Commodities {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Guns {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Missiles {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Mines {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Shields {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Thrusters {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Ships {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Tractors {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Engines {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.CMs {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Scanners {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Ammos {
		enhance_with_transport_unrechability(item.Bases)
	}
	for _, item := range e.Cloaks {
		enhance_with_transport_unrechability(item.Bases)
	}
}

func Export(mapped *configs_mapped.MappedConfigs, options ExportOptions) *Exporter {
	return NewExporter(mapped).Export(options)
}

func Empty(phrase string) bool {
	for _, letter := range phrase {
		if letter != ' ' {
			return false
		}
	}
	return true
}

func (e *Exporter) Buyable(Bases map[cfg.BaseUniNick]*MarketGood) bool {
	for _, base := range Bases {

		if e.useful_bases_by_nick != nil {
			if _, ok := e.useful_bases_by_nick[base.BaseNickname]; ok {
				return true
			}
		}
	}

	return false
}

func Buyable(Bases map[cfg.BaseUniNick]*MarketGood) bool {
	return len(Bases) > 0
}

type DiscoveryTechCompat struct {
	TechcompatByID map[cfg.TractorID]float64 `json:"techchompat_by_id" validate:"required"`
	TechCell       string                    `json:"tech_cell" validate:"required"`
}

func CalculateTechCompat(Discovery *configs_mapped.DiscoveryConfig, ids []*Tractor, nickname string) *DiscoveryTechCompat {
	if Discovery == nil {
		return nil
	}

	techcompat := &DiscoveryTechCompat{
		TechcompatByID: make(map[cfg.TractorID]float64),
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
	return e.Mapped.GetInfocardName(ids_name, nickname)
}
func (e *ExporterRelay) GetInfocardName(ids_name int, nickname string) string {
	return e.Mapped.GetInfocardName(ids_name, nickname)
}
