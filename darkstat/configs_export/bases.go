package configs_export

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (e *Exporter) TraderExists(base_nickname string) bool {
	universe_base, ok := e.Configs.Universe.BasesMap[universe_mapped.BaseNickname(base_nickname)]
	if !ok {
		return true
	}
	return universe_base.TraderExists
}

func VectorToSectorCoord(system *universe_mapped.System, pos cfgtype.Vector) string {
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

func (e *Exporter) GetBases() []*Base {
	results := make([]*Base, 0, len(e.Configs.Universe.Bases))

	commodities_per_base := e.getMarketGoods()

	for _, base := range e.Configs.Universe.Bases {
		var name string = e.GetInfocardName(base.StridName.Get(), base.Nickname.Get())

		var system_name infocard.Infoname
		var Region string
		system, found_system := e.Configs.Universe.SystemMap[universe_mapped.SystemNickname(base.System.Get())]

		if found_system {

			system_name = infocard.Infoname(e.GetInfocardName(system.StridName.Get(), system.Nickname.Get()))

			Region = e.GetRegionName(system)
		}

		var infocard_id int
		var reputation_nickname string
		var pos cfgtype.Vector

		var archetypes []string

		if system, ok := e.Configs.Systems.SystemsMap[base.System.Get()]; ok {
			if system_base, ok := system.BasesByBases[base.Nickname.Get()]; ok {
				infocard_id = system_base.IDsInfo.Get()
				reputation_nickname = system_base.RepNickname.Get()
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

		if infocard_middle_id, exists := e.Configs.InfocardmapINI.InfocardMapTable.Map[infocard_id]; exists {
			infocard_ids = append(infocard_ids, infocard_middle_id)
		}

		var factionName string
		if group, exists := e.Configs.InitialWorld.GroupsMap[reputation_nickname]; exists {
			infocard_ids = append(infocard_ids, group.IdsInfo.Get())
			factionName = e.GetInfocardName(group.IdsName.Get(), reputation_nickname)
		}

		var market_goods_per_good_nick map[CommodityKey]MarketGood = make(map[CommodityKey]MarketGood)

		if found_commodities, ok := commodities_per_base[cfgtype.BaseUniNick(base.Nickname.Get())]; ok {
			market_goods_per_good_nick = found_commodities
		}

		var nickname cfgtype.BaseUniNick = cfgtype.BaseUniNick(base.Nickname.Get())

		e.exportInfocards(InfocardKey(nickname), infocard_ids...)

		base := &Base{
			Missions:           &BaseMissions{},
			Name:               name,
			Nickname:           nickname,
			NicknameHash:       flhash.HashNickname(nickname.ToStr()),
			FactionName:        factionName,
			System:             string(system_name),
			SystemNickname:     base.System.Get(),
			SystemNicknameHash: flhash.HashNickname(base.System.Get()),
			StridName:          base.StridName.Get(),
			InfocardID:         infocard_id,
			InfocardKey:        InfocardKey(nickname),
			File:               utils_types.FilePath(base.File.Get()),
			BGCS_base_run_by:   base.BGCS_base_run_by.Get(),
			MarketGoodsPerNick: market_goods_per_good_nick,
			Pos:                pos,
			Archetypes:         archetypes,
			Region:             Region,
		}

		e.Hashes[string(base.Nickname)] = base.NicknameHash
		e.Hashes[base.SystemNickname] = base.SystemNicknameHash

		if found_system {
			base.SectorCoord = VectorToSectorCoord(system, base.Pos)
		}

		base.Infocard = e.Infocards[InfocardKey(base.Nickname)]

		results = append(results, base)
	}

	return results
}

func EnhanceBasesWithServerOverrides(bases []*Base, commodities []*Commodity) {
	var base_per_nick map[cfgtype.BaseUniNick]*Base = make(map[cfgtype.BaseUniNick]*Base)
	for _, base := range bases {
		base_per_nick[base.Nickname] = base
	}

	for _, commodity := range commodities {
		for _, base_location := range commodity.Bases {

			// no idea why :) but crashes otherwise
			if base_location.BaseNickname == pob_crafts_nickname || base_location.BaseNickname == BaseLootableNickname {
				continue
			}
			if base_location.IsServerSideOverride {
				continue
			}

			var market_good MarketGood
			commodity_key := GetCommodityKey(commodity.Nickname, commodity.ShipClass)

			if base, ok := base_per_nick[base_location.BaseNickname]; ok {
				if good, ok := base.MarketGoodsPerNick[commodity_key]; ok {
					market_good = good
				}
			}

			market_good.Name = commodity.Name
			market_good.Nickname = commodity.Nickname
			market_good.NicknameHash = commodity.NicknameHash
			market_good.Category = "commodity"

			market_good.LevelRequired = base_location.LevelRequired
			market_good.RepRequired = base_location.RepRequired
			market_good.BaseSells = base_location.BaseSells

			market_good.PriceToBuy = base_location.PriceBaseSellsFor
			market_good.PriceToSell = base_location.PriceBaseBuysFor
			market_good.Infocard = commodity.Infocard
			market_good.Volume = commodity.Volume
			market_good.IsServerSideOverride = base_location.IsServerSideOverride

			base_per_nick[base_location.BaseNickname].MarketGoodsPerNick[commodity_key] = market_good
		}
	}
}

func FilterToUserfulBases(bases []*Base) []*Base {
	var useful_bases []*Base = make([]*Base, 0, len(bases))
	for _, item := range bases {
		if item.Reachable {
			useful_bases = append(useful_bases, item)
			continue
		}

		if (item.Name == "Object Unknown" || item.Name == "") && len(item.MarketGoodsPerNick) == 0 {
			continue
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
		useful_bases = append(useful_bases, item)
	}
	return useful_bases
}

type Base struct {
	Name               string              `json:"name"`       // Infocard Name
	Archetypes         []string            `json:"archetypes"` // Base Archetypes
	Nickname           cfgtype.BaseUniNick `json:"nickname"`
	NicknameHash       flhash.HashCode     `json:"nickname_hash" format:"int64"` // Flhash of nickname
	FactionName        string              `json:"faction_nickname"`
	System             string              `json:"system_name"`
	SystemNickname     string              `json:"system_nickname"`
	SystemNicknameHash flhash.HashCode     `json:"system_nickname_hash" format:"int64"`
	Region             string              `json:"region_name"`
	StridName          int                 `json:"strid_name"`
	InfocardID         int                 `json:"infocard_id"`
	Infocard           Infocard            `json:"infocard"`
	InfocardKey        InfocardKey
	File               utils_types.FilePath `json:"file"`
	BGCS_base_run_by   string
	MarketGoodsPerNick map[CommodityKey]MarketGood `json:"-"`
	Pos                cfgtype.Vector              `json:"pos"`
	SectorCoord        string                      `json:"sector_coord"`

	IsTransportUnreachable bool `json:"is_transport_unreachable"` // Check if base is NOT reachable from manhattan by Transport through Graph method (at Discovery base has to have Transport dockable spheres)

	Missions           *BaseMissions `json:"-"`
	baseAllTradeRoutes `json:"-"`
	baseAllRoutes      `json:"-"`
	*MiningInfo        `json:"mining_info,omitempty"`

	Reachable bool `json:"is_reachhable"` // is base IS Rechable by frighter from Manhattan
}

type CommodityKey string

func GetCommodityKey(nickname string, ship_class cfgtype.ShipClass) CommodityKey {
	return CommodityKey(fmt.Sprintf("%s_%d", nickname, ship_class))
}
