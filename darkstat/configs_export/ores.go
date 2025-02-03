package configs_export

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type MiningInfo struct {
	DynamicLootMin        int
	DynamicLootMax        int
	DynamicLootDifficulty int
	MinedGood             *MarketGood
}

func (e *Exporter) GetOres(Commodities []*Commodity) []*Base {
	var bases []*Base

	var comm_by_nick map[CommodityKey]*Commodity = make(map[CommodityKey]*Commodity)
	for _, comm := range Commodities {
		comm_by_nick[GetCommodityKey(comm.Nickname, comm.ShipClass)] = comm
	}

	for _, system := range e.Configs.Systems.Systems {

		system_uni, system_uni_ok := e.Configs.Universe.SystemMap[universe_mapped.SystemNickname(system.Nickname)]

		for _, asteroids := range system.Asteroids {

			asteroid_zone_nick := asteroids.Zone.Get()

			zone, found_zone := system.ZonesByNick[asteroid_zone_nick]
			if !found_zone {
				continue
			}

			if asteroids.LootableZone == nil {
				continue
			}

			commodity, commodity_found := asteroids.LootableZone.AsteroidLootCommodity.GetValue()

			if !commodity_found {
				continue
			}

			location := zone.Pos.Get()
			var added_goods map[string]bool = make(map[string]bool)
			base := &Base{
				MiningInfo:         &MiningInfo{},
				Pos:                location,
				MarketGoodsPerNick: make(map[CommodityKey]*MarketGood),
			}
			base.DynamicLootMin, _ = asteroids.LootableZone.DynamicLootMin.GetValue()
			base.DynamicLootMax, _ = asteroids.LootableZone.DynamicLootMax.GetValue()
			base.DynamicLootDifficulty, _ = asteroids.LootableZone.DynamicLootDifficulty.GetValue()

			var base_nickname string
			base_nickname, _ = zone.Nickname.GetValue()
			base.Nickname = cfg.BaseUniNick(base_nickname)

			base.NicknameHash = flhash.HashNickname(base_nickname)
			e.Hashes[base_nickname] = base.NicknameHash

			base.InfocardID, _ = zone.IDsInfo.GetValue()
			base.StridName, _ = zone.IdsName.GetValue()

			base.InfocardKey = InfocardKey(base.Nickname)

			base.Archetypes = append(base.Archetypes, "mining_operation")
			base.FactionName = "Mining Field"

			base.SystemNickname = system.Nickname
			base.SystemNicknameHash = flhash.HashNickname(base.SystemNickname)
			if system, ok := e.Configs.Universe.SystemMap[universe_mapped.SystemNickname(base.SystemNickname)]; ok {
				base.System = e.GetInfocardName(system.StridName.Get(), base.SystemNickname)
				base.Region = e.GetRegionName(system)
				base.SectorCoord = VectorToSectorCoord(system_uni, base.Pos)
			}

			logus.Log.Debug("GetOres", typelog.String("commodity=", commodity))

			equipment := e.Configs.Equip.CommoditiesMap[commodity]
			for _, volume_info := range equipment.Volumes {

				market_good := &MarketGood{
					GoodInfo:          e.GetGoodInfo(commodity),
					BaseSells:         true,
					PriceBaseSellsFor: 0,
					PriceBaseBuysFor:  nil,
					Volume:            volume_info.Volume.Get(),
					ShipClass:         volume_info.GetShipClass(),
					BaseInfo: BaseInfo{
						BaseNickname: base.Nickname,
						BaseName:     base.Name,
						SystemName:   base.System,
						BasePos:      base.Pos,
						Region:       base.Region,
						FactionName:  "Mining Field",
						SectorCoord:  base.SectorCoord,
					},
				}
				base.Name = market_good.Name
				market_good.BaseName = market_good.Name

				market_good_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
				base.MarketGoodsPerNick[market_good_key] = market_good
				base.MinedGood = market_good
				added_goods[market_good.Nickname] = true

				if commodity, ok := comm_by_nick[market_good_key]; ok {
					commodity.Bases[market_good.BaseNickname] = market_good
				}

			}

			if e.Configs.Discovery != nil {
				if recipes, ok := e.Configs.Discovery.BaseRecipeItems.RecipePerConsumed[commodity]; ok {
					for _, recipe := range recipes {
						recipe_produces_only_commodities := true

						for _, produced := range recipe.ProcucedItem {

							_, is_commodity := e.Configs.Equip.CommoditiesMap[produced.Get()]
							if !is_commodity {
								recipe_produces_only_commodities = false
								break
							}

						}

						if recipe_produces_only_commodities {
							for _, produced := range recipe.ProcucedItem {
								commodity_produced := produced.Get()

								if _, ok := added_goods[commodity_produced]; ok {
									continue
								}
								equipment := e.Configs.Equip.CommoditiesMap[commodity_produced]
								for _, volume_info := range equipment.Volumes {
									market_good := &MarketGood{
										GoodInfo:          e.GetGoodInfo(commodity_produced),
										BaseSells:         true,
										PriceBaseSellsFor: 0,
										PriceBaseBuysFor:  nil,
										Volume:            volume_info.Volume.Get(),
										ShipClass:         volume_info.GetShipClass(),
										BaseInfo: BaseInfo{
											BaseNickname: base.Nickname,
											BaseName:     base.Name,
											SystemName:   base.System,
											BasePos:      base.Pos,
											Region:       base.Region,
											FactionName:  "Mining Field",
										},
									}
									market_good.BaseName = market_good.Name
									if system_uni_ok {
										market_good.SectorCoord = VectorToSectorCoord(system_uni, market_good.BasePos)
									}
									market_good_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
									base.MarketGoodsPerNick[market_good_key] = market_good
									if commodity, ok := comm_by_nick[market_good_key]; ok {
										commodity.Bases[market_good.BaseNickname] = market_good
									}
									added_goods[commodity_produced] = true
								}

							}
						}
					}

				}
			}

			var sb InfocardBuilder
			sb.WriteLineStr(base.Name)
			sb.WriteLineStr((`This is is not a base.
It is a mining field with droppable ores`))
			sb.WriteLineStr((""))
			sb.WriteLineStr(("Trade routes shown do not account for a time it takes to mine those ores."))

			if e.Configs.Discovery != nil {
				sb.WriteLineStr("")
				sb.WriteLine(InfocardPhrase{Link: ptr.Ptr("https://discoverygc.com/wiki2/Mining"), Phrase: "Check mining tutorial"}, InfocardPhrase{Phrase: " to see how they can be mined"})

				sb.WriteLineStr("")
				sb.WriteLineStr(`NOTE:
for Freelancer Discovery we also add possible sub products of refinery at player bases to possible trade routes from mining field.
				`)
			}

			sb.WriteLineStr("")
			sb.WriteLineStr("commodities:")
			for _, good := range base.MarketGoodsPerNick {
				if good.Nickname == base.MinedGood.Nickname {
					sb.WriteLineStr(fmt.Sprintf("Minable: %s (%s)", good.Name, good.Nickname))
				} else {
					sb.WriteLineStr(fmt.Sprintf("Refined at POB: %s (%s)", good.Name, good.Nickname))
				}
			}

			e.Infocards[InfocardKey(base.Nickname)] = sb.Lines

			bases = append(bases, base)

		}
		_ = system
	}

	return bases
}

var not_useful_ores []string = []string{
	"commodity_water",              // sellable
	"commodity_oxygen",             // sellable
	"commodity_scrap_metal",        // sellable
	"commodity_toxic_waste",        // a bit
	"commodity_cerulite_crystals",  // not
	"commodity_alien_organisms",    // sellable
	"commodity_hydrocarbons",       // sellable
	"commodity_inert_artifacts",    // not
	"commodity_organic_capacitors", // not
	"commodity_event_ore_01",       // not
	"commodity_cryo_organisms",     // not
	"commodity_chirodebris",        // not
}

func FitlerToUsefulOres(bases []*Base) []*Base {
	var useful_bases []*Base = make([]*Base, 0, len(bases))
	for _, item := range bases {
		if strings.Contains(item.System, "Bastille") {
			continue
		}

		is_useful := true
		for _, useless_commodity := range not_useful_ores {
			if item.MinedGood.Nickname == useless_commodity {
				is_useful = false
				break
			}

		}
		if !is_useful {
			continue
		}

		useful_bases = append(useful_bases, item)
	}
	return useful_bases
}
