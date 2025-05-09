package configs_export

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
)

func (e *Exporter) findable_in_loot() map[string]bool {
	if e.findable_in_loot_cache != nil {
		return e.findable_in_loot_cache
	}

	e.findable_in_loot_cache = make(map[string]bool)

	for _, system := range e.Mapped.Systems.Systems {
		for _, wreck := range system.Wrecks {
			louadout_nickname := wreck.Loadout.Get()
			if loadout, ok := e.Mapped.Loadouts.LoadoutsByNick[louadout_nickname]; ok {
				for _, cargo := range loadout.Cargos {
					e.findable_in_loot_cache[cargo.Nickname.Get()] = true
				}
			}
		}
	}

	for _, npc_arch := range e.Mapped.NpcShips.NpcShips {
		loadout_nickname := npc_arch.Loadout.Get()
		if loadout, ok := e.Mapped.Loadouts.LoadoutsByNick[loadout_nickname]; ok {
			for _, cargo := range loadout.Cargos {
				e.findable_in_loot_cache[cargo.Nickname.Get()] = true
			}
		}
	}
	return e.findable_in_loot_cache
}

/*
It fixes issue of Guns obtainable only via wrecks being invisible
*/
const (
	BaseLootableName     = "Lootable"
	BaseLootableFaction  = "Wrecks and Missions"
	BaseLootableNickname = "base_loots"
)

func (e *Exporter) EnhanceBasesWithLoot(bases []*Base) []*Base {

	in_wrecks := e.findable_in_loot()

	base := &Base{
		Name:               "Lootable",
		MarketGoodsPerNick: make(map[CommodityKey]*MarketGood),
		Nickname:           cfg.BaseUniNick(BaseLootableNickname),
		SystemNickname:     "neverwhere",
		System:             "Neverwhere",
		Region:             "NEVERWHERE",
		FactionName:        BaseLootableFaction,
	}

	base.Archetypes = append(base.Archetypes, BaseLootableNickname)

	for wreck, _ := range in_wrecks {
		market_good := &MarketGood{
			GoodInfo:             e.GetGoodInfo(wreck),
			BaseSells:            true,
			IsServerSideOverride: true,
		}

		market_good_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
		base.MarketGoodsPerNick[market_good_key] = market_good
	}

	var sb []infocarder.InfocardLine
	sb = append(sb, infocarder.NewInfocardSimpleLine(base.Name))
	sb = append(sb, infocarder.NewInfocardSimpleLine(`This is only pseudo base to show availability of lootable content`))
	sb = append(sb, infocarder.NewInfocardSimpleLine(`The content is findable in wrecks or drops from ships at missions`))

	e.PutInfocard(infocarder.InfocardKey(base.Nickname), sb)

	bases = append(bases, base)
	return bases
}
