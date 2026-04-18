package configs_export

import (
	"context"
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/solararch_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (e *Exporter) TraderExists(base_nickname string) bool {
	universe_base, ok := e.Mapped.Universe.BasesMap[universe_mapped.BaseNickname(base_nickname)]
	if !ok {
		return true
	}
	return universe_base.TraderExists
}

func VectorToSectorCoord(system *universe_mapped.System, pos cfg.Vector) string {
	var scale float64 = 1.0
	if value, ok := system.NavMapScale.GetValue(); ok {
		scale = value
	}

	var fGridSize float64 = 34000.0 / scale // 34000 suspiciously looks like math.MaxInt16
	var gridRefX = int((pos.X+(fGridSize*5))/fGridSize) - 1
	var gridRefZ = int((pos.Z+(fGridSize*5))/fGridSize) - 1
	gridRefX = min(max(gridRefX, 0), 7)
	scXPos := rune('A' + gridRefX)
	gridRefZ = min(max(gridRefZ, 0), 7)
	scZPos := rune('1' + gridRefZ)
	return fmt.Sprintf("%c-%c", scXPos, scZPos)

}

func (e *Exporter) GetBases(ctx context.Context) []*Base {
	ctx, span := traces.Tracer.Start(ctx, "Exporter.GetBases")
	defer span.End()

	results := make([]*Base, 0, len(e.Mapped.Universe.Bases))

	commodities_per_base := e.getMarketGoods()

	for _, base := range e.Mapped.Universe.Bases {
		var name string = e.GetInfocardName(base.StridName.Get(), base.Nickname.Get())

		var system_name infocard.Infoname
		var Region string
		system, found_system := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(base.System.Get())]

		if found_system {

			system_name = infocard.Infoname(e.GetInfocardName(system.StridName.Get(), system.Nickname.Get()))

			Region = e.GetRegionName(system)
		}

		var infocard_id int
		var reputation_nickname string
		var pos cfg.Vector

		var archetypes []string
		var obj_nickname string

		if system, ok := e.Mapped.Systems.SystemsMap[base.System.Get()]; ok {
			if system_base, ok := system.BasesByBases[base.Nickname.Get()]; ok {
				infocard_id = system_base.IDsInfo.Get()
				reputation_nickname = system_base.RepNickname.Get()
				obj_nickname = system_base.Nickname.Get()
			}

			if system_bases, ok := system.AllBasesByDockWith[base.Nickname.Get()]; ok {
				for _, system_base := range system_bases {
					pos, _ = system_base.Pos.GetValue()
					archetype, _ := system_base.Archetype.GetValue()
					archetypes = append(archetypes, archetype)
				}
			}
		}

		var infocard_ids []int = make([]int, 0)
		infocard_ids = append(infocard_ids, infocard_id)
		if infocard_middle_id, exists := e.Mapped.InfocardmapINI.InfocardMapTable.Map[infocard_id]; exists {
			infocard_ids = append(infocard_ids, infocard_middle_id)
		}

		var factionName string
		if group, exists := e.Mapped.InitialWorld.GroupsMap[reputation_nickname]; exists {
			infocard_ids = append(infocard_ids, group.IdsInfo.Get())
			factionName = e.GetInfocardName(group.IdsName.Get(), reputation_nickname)
		}

		var market_goods_per_good_nick map[CommodityKey]*MarketGood = make(map[CommodityKey]*MarketGood)

		if found_commodities, ok := commodities_per_base[cfg.BaseUniNick(base.Nickname.Get())]; ok {
			market_goods_per_good_nick = found_commodities
		}

		var nickname cfg.BaseUniNick = cfg.BaseUniNick(base.Nickname.Get())

		e.ExportInfocards(infocarder.InfocardKey(nickname), infocard_ids...)
		e.WriteConfigToInfocard(&base.Model, string(nickname))

		if system, ok := e.Mapped.Systems.SystemsMap[base.System.Get()]; ok {

			if system_bases, ok := system.AllBasesByDockWith[base.Nickname.Get()]; ok {
				for _, system_base := range system_bases {
					e.WriteConfigToInfocard(&system_base.Model, string(nickname))

					archetype := system_base.Archetype.Get()
					solar := e.Mapped.Solararch.SolarsByNick[archetype]
					e.WriteConfigToInfocard(&solar.Model, string(nickname))

				}
			}
		}

		base := &Base{
			Missions:           &BaseMissions{},
			Name:               name,
			Nickname:           nickname,
			ObjNickname:        obj_nickname,
			FactionName:        factionName,
			System:             string(system_name),
			SystemNickname:     base.System.Get(),
			StridName:          base.StridName.Get(),
			InfocardID:         infocard_id,
			File:               utils_types.FilePath(base.File.Get()),
			BGCS_base_run_by:   base.BGCS_base_run_by.Get(),
			MarketGoodsPerNick: market_goods_per_good_nick,
			Pos:                pos,
			Archetypes:         archetypes,
			Region:             Region,
		}

		if found_system {
			base.SectorCoord = VectorToSectorCoord(system, base.Pos)
		}

		results = append(results, base)
	}

	return results
}

func EnhanceBasesWithServerOverrides(bases []*Base, commodities []*Commodity) {
	var base_per_nick map[cfg.BaseUniNick]*Base = make(map[cfg.BaseUniNick]*Base)
	for _, base := range bases {
		base_per_nick[base.Nickname] = base
	}

	for _, commodity := range commodities {
		for _, market_good := range commodity.Bases {
			commodity_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
			if base, ok := base_per_nick[market_good.BaseNickname]; ok {
				base.MarketGoodsPerNick[commodity_key] = market_good
			}
		}
	}
}

// you can use for_front_view=true, only after bases were already enhanced with IsTransportReachable and etc values
// otheriwise u have to use for_front_view=false
func (e *Exporter) FilterToUserfulBases(bases []*Base, for_front_view bool) []*Base {
	var useful_bases []*Base = make([]*Base, 0, len(bases))
	for _, item := range bases {

		if item.IsPob {
			useful_bases = append(useful_bases, item)
			continue
		}

		if len(item.MarketGoodsPerNick) == 0 {
			continue
		}

		if for_front_view {
			if !item.IsFreighterReachable && !item.IsFrigateReachable && !item.IsTransportReachable {
				continue
			}
		}

		if strings.Contains(item.System, "Bastille") {
			continue
		}

		is_invisible := true
		for _, archetype := range item.Archetypes {
			if archetype != systems_mapped.BaseArchetypeInvisible {
				is_invisible = false
			}
		}
		if is_invisible {
			continue
		}

		is_dockable := false
		for _, archetype := range item.Archetypes {

			if solar, ok := e.Mapped.Solararch.SolarsByNick[archetype]; ok {
				dockable := solar.IsDockable(solararch_mapped.DockableOptions{
					IsDisco:                  e.Mapped.Discovery != nil,
					PlayersCanDockBerth:      true,
					PlayersCanDockMoorMedium: true,
					PlayersCanDockMoorLarge:  true,
				})
				if dockable.IsDockable {
					is_dockable = true
				}
			}
		}
		if !is_dockable {
			continue
		}

		useful_bases = append(useful_bases, item)
	}
	return useful_bases
}

type Base struct {
	Name               string               `json:"name"  validate:"required"`       // Infocard Name
	Archetypes         []string             `json:"archetypes"  validate:"required"` // Base Archetypes
	Nickname           cfg.BaseUniNick      `json:"nickname"  validate:"required"`
	ObjNickname        string               `json:"obj_nickname"`
	FactionName        string               `json:"faction_name"  validate:"required"`
	System             string               `json:"system_name"  validate:"required"`
	SystemNickname     string               `json:"system_nickname"  validate:"required"`
	Region             string               `json:"region_name"  validate:"required"`
	StridName          int                  `json:"strid_name"  validate:"required"`
	InfocardID         int                  `json:"infocard_id"  validate:"required"`
	File               utils_types.FilePath `json:"file" validate:"required"`
	BGCS_base_run_by   string
	MarketGoodsPerNick map[CommodityKey]*MarketGood `json:"-" swaggerignore:"true"`
	Pos                cfg.Vector                   `json:"pos" validate:"required"`
	SectorCoord        string                       `json:"sector_coord" validate:"required"`

	cfg.Reachability

	Missions    *BaseMissions `json:"-" swaggerignore:"true"`
	*MiningInfo `json:"mining_info,omitempty"`
	IsPob       bool `validate:"required"`
}

func (b Base) GetNickname() string { return string(b.Nickname) }

type CommodityKey string

func (c CommodityKey) ToStr() string { return string(c) }

func GetCommodityKey(nickname string, ship_class *cfg.ShipClass) CommodityKey {
	return CommodityKey(fmt.Sprintf("%s_%s", nickname, cfg.ShipClassToKey(ship_class)))
}
