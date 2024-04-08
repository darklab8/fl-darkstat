package configs_export

import (
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
)

type Exporter struct {
	configs            *configs_mapped.MappedConfigs
	show_empty_records bool

	Bases            []Base
	Factions         []Faction
	Infocards        map[InfocardKey]*Infocard
	Commodities      []Commodity
	Guns             []Gun
	Missiles         []Gun
	Mines            []Mine
	Shields          []Shield
	Thrusters        []Thruster
	Ships            []Ship
	Tractors         []Tractor
	Engines          []Engine
	CMs              []CounterMeasure
	infocards_parser *InfocardsParser
}

type OptExport func(e *Exporter)

func WithEmptyRecords() OptExport {
	return func(e *Exporter) { e.show_empty_records = true }
}

func NewExporter(configs *configs_mapped.MappedConfigs, opts ...OptExport) *Exporter {
	e := &Exporter{
		configs:            configs,
		show_empty_records: false,
		infocards_parser:   NewInfocardsParser(configs.Infocards),
	}

	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Exporter) Export() *Exporter {
	e.Bases = e.GetBases()
	e.Factions = e.GetFactions(e.Bases)
	e.Commodities = e.GetCommodities()
	e.Guns = e.GetGuns()
	e.Missiles = e.GetMissiles()
	e.Mines = e.GetMines()
	e.Shields = e.GetShields()
	e.Thrusters = e.GetThrusters()
	e.Ships = e.GetShips()
	e.Tractors = e.GetTractors()
	e.Engines = e.GetEngines()
	e.CMs = e.GetCounterMeasures()
	e.Infocards = e.infocards_parser.Get()
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

func Buyable(Bases []GoodAtBase) bool {
	for _, base := range Bases {
		if !strings.Contains(base.SystemName, "Bastille") {
			return true
		}
	}

	return false
}
