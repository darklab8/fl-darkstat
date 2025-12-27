package configs_export

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/discovery/minecontrol_nodes"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type MiningInfo struct {
	DynamicLootMin        int
	DynamicLootMax        int
	DynamicLootDifficulty int
	MinedGood             *MarketGood
	RespawnCooldown       int
	MaxSpawnCount         int
	NuAsteroidsCount      int
	NuAsteroidsTotal      int
}

func (e *Exporter) GetOres(ctx context.Context, Commodities []*Commodity) []*Base {
	ctx, span := traces.Tracer.Start(ctx, "Exporter.GetOres")
	defer span.End()
	var bases []*Base

	var comm_by_nick map[CommodityKey]*Commodity = make(map[CommodityKey]*Commodity)
	for _, comm := range Commodities {
		comm_by_nick[GetCommodityKey(comm.Nickname, comm.ShipClass)] = comm
	}

	for _, system := range e.Mapped.Systems.Systems {

		system_uni, system_uni_ok := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(system.Nickname)]

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

			var base_nickname string
			base_nickname, _ = zone.Nickname.GetValue()
			bases = append(bases, e.NewOreBase(
				NewOreBaseInput{
					base_nickname: base_nickname,
					location:      location,
					commodity:     commodity,
					system_uni:    system_uni,
					system_uni_ok: system_uni_ok,
					system:        system,
					comm_by_nick:  comm_by_nick,
					asteroids:     asteroids,
					zone:          zone,
				},
			))

		}

		if e.Mapped.Discovery != nil {
			if mining_systems, ok := e.Mapped.Discovery.MinecontrolNodes.MiningSystemsBySystemNick[system_uni.Nickname.Get()]; ok {
				for _, mining_system := range mining_systems {

					base_nickname := mining_system.Nickname.Get()

					var unique_commodity_of_mining_system map[string]bool = make(map[string]bool)

					for _, node_archetype := range mining_system.NodeArchetypes {
						node_arch_nickname := node_archetype.Nickname.Get()

						mining_solar, ok := e.Mapped.Discovery.Minecontrol.MiningSolarByAsteroidNick[node_arch_nickname]
						if !ok {
							error_msg := fmt.Sprintln("cant find node_arch_nickname=", node_arch_nickname, " in MiningSolars")
							logus.Log.Error(error_msg)
							continue // TODO remove this hack when disco devs will fix minecontrol.cfg
							logus.Log.Panic(error_msg)
						}

						for ore_nickname, _ := range mining_solar.OreNicknamesCounts {
							commodity := ore_nickname

							_, already_used := unique_commodity_of_mining_system[commodity]
							if already_used {
								continue
							}

							nu_asteroids_count := 0
							nu_asteroids_total := 0
							for _, node_archetype := range mining_system.NodeArchetypes {
								node_arch_nickname := node_archetype.Nickname.Get()
								mining_solar, ok := e.Mapped.Discovery.Minecontrol.MiningSolarByAsteroidNick[node_arch_nickname]
								if !ok {
									error_msg := fmt.Sprintln("cant find node_arch_nickname=", node_arch_nickname, " in MiningSolars")
									logus.Log.Error(error_msg)
									continue // TODO remove this hack when disco devs will fix minecontrol.cfg
									logus.Log.Panic(error_msg)
								}

								nu_asteroids_count += mining_solar.OreNicknamesCounts[commodity]
								nu_asteroids_total += mining_solar.OreCount
							}

							unique_commodity_of_mining_system[commodity] = true
							bases = append(bases, e.NewOreBase(
								NewOreBaseInput{
									base_nickname: fmt.Sprintf("%s (%s)", base_nickname, commodity),
									location:      mining_system.Position,
									commodity:     commodity,
									system_uni:    system_uni,
									system_uni_ok: system_uni_ok,
									system:        system,
									comm_by_nick:  comm_by_nick,
									asteroids:     nil,
									zone:          nil,
									disco_mining_system: &DiscoMiningSystem{
										base_nickname:    mining_system.Nickname.Get(),
										mining_system:    mining_system,
										NuAsteroidsCount: nu_asteroids_count,
										NuAsteroidsTotal: nu_asteroids_total,
									},
								},
							))
						}
					}

				}

			}
		}
		_ = system
	}

	return bases
}

type DiscoMiningSystem struct {
	base_nickname    string
	mining_system    *minecontrol_nodes.MiningSystem
	NuAsteroidsCount int
	NuAsteroidsTotal int
}

type NewOreBaseInput struct {
	base_nickname       string
	location            cfg.Vector
	commodity           string
	system_uni          *universe_mapped.System
	system_uni_ok       bool
	system              *systems_mapped.System
	comm_by_nick        map[CommodityKey]*Commodity
	asteroids           *systems_mapped.Asteroids
	zone                *systems_mapped.Zone
	disco_mining_system *DiscoMiningSystem
}

func (e *Exporter) NewOreBase(input_data NewOreBaseInput) *Base {
	var added_goods map[string]bool = make(map[string]bool)
	base := &Base{
		MiningInfo:         &MiningInfo{},
		Pos:                input_data.location,
		MarketGoodsPerNick: make(map[CommodityKey]*MarketGood),
	}

	base.Nickname = cfg.BaseUniNick(input_data.base_nickname)
	if input_data.asteroids != nil {
		base.DynamicLootMin, _ = input_data.asteroids.LootableZone.DynamicLootMin.GetValue()
		base.DynamicLootMax, _ = input_data.asteroids.LootableZone.DynamicLootMax.GetValue()
		base.DynamicLootDifficulty, _ = input_data.asteroids.LootableZone.DynamicLootDifficulty.GetValue()
	}
	if input_data.zone != nil {
		base.InfocardID, _ = input_data.zone.IDsInfo.GetValue()
		base.StridName, _ = input_data.zone.IdsName.GetValue()
	}
	if input_data.disco_mining_system != nil {
		base.RespawnCooldown = input_data.disco_mining_system.mining_system.RespawnCooldown.Get()
		base.MaxSpawnCount = input_data.disco_mining_system.mining_system.MaxSpawnCount.Get()

		base.NuAsteroidsCount = input_data.disco_mining_system.NuAsteroidsCount
		base.NuAsteroidsTotal = input_data.disco_mining_system.NuAsteroidsTotal
	}

	base.Archetypes = append(base.Archetypes, "mining_operation")
	base.FactionName = "Mining Field"

	base.SystemNickname = input_data.system.Nickname
	if system, ok := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(base.SystemNickname)]; ok {
		base.System = e.GetInfocardName(system.StridName.Get(), base.SystemNickname)
		base.Region = e.GetRegionName(system)
		base.SectorCoord = VectorToSectorCoord(input_data.system_uni, base.Pos)
	}

	logus.Log.Debug("GetOres", typelog.String("commodity=", input_data.commodity))

	equipment := e.Mapped.Equip().CommoditiesMap[input_data.commodity]
	for _, volume_info := range equipment.Volumes {

		market_good := &MarketGood{
			GoodInfo:          e.GetGoodInfo(input_data.commodity),
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
		if input_data.disco_mining_system != nil {
			base.Name += " (Nu)"
		}

		market_good.BaseName = market_good.Name

		market_good_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
		base.MarketGoodsPerNick[market_good_key] = market_good
		base.MinedGood = market_good
		added_goods[market_good.Nickname] = true

		if commodity, ok := input_data.comm_by_nick[market_good_key]; ok {
			commodity.Bases[market_good.BaseNickname] = market_good
		}

	}

	if e.Mapped.Discovery != nil {
		if recipes, ok := e.Mapped.Discovery.BaseRecipeItems.RecipePerConsumed[input_data.commodity]; ok {
			for _, recipe := range recipes {
				recipe_produces_only_commodities := true

				for _, produced := range recipe.ProcucedItem {

					_, is_commodity := e.Mapped.Equip().CommoditiesMap[produced.Get()]
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
						equipment := e.Mapped.Equip().CommoditiesMap[commodity_produced]
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
							if input_data.system_uni_ok {
								market_good.SectorCoord = VectorToSectorCoord(input_data.system_uni, market_good.BasePos)
							}
							market_good_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
							base.MarketGoodsPerNick[market_good_key] = market_good
							if commodity, ok := input_data.comm_by_nick[market_good_key]; ok {
								commodity.Bases[market_good.BaseNickname] = market_good
							}
							added_goods[commodity_produced] = true
						}

					}
				}
			}

		}
	}

	var sb infocarder.InfocardBuilder
	sb.WriteLineStr(base.Name)
	sb.WriteLineStr((`This is is not a base.
It is a mining field with droppable ores`))
	sb.WriteLineStr((""))
	sb.WriteLineStr(("Trade routes shown do not account for a time it takes to mine those ores."))

	if e.Mapped.Discovery != nil {
		sb.WriteLineStr("")
		sb.WriteLine(infocarder.InfocardPhrase{Link: ptr.Ptr("https://discoverygc.com/wiki2/Mining"), Phrase: "Check mining tutorial"}, infocarder.InfocardPhrase{Phrase: " to see how they can be mined"})

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
	sb.WriteLineStr("")

	var infocard_addition infocarder.InfocardBuilder
	if e.Mapped.Discovery != nil {
		if player_bonuses, ok := e.Mapped.Discovery.Minecontrol.PlayerBonusByOreNickname[base.MinedGood.Nickname]; ok {
			infocard_addition.WriteLineStr(`MINING BONUSES (darkstat):`)
			for _, player_bonus := range player_bonuses {
				id_nickname := player_bonus.IDNickname.Get()
				id_name := id_nickname
				if tractor, ok := e.Mapped.Equip().TractorsMap[id_nickname]; ok {
					if name_id, ok := tractor.IdsName.GetValue(); ok {
						id_name = e.GetInfocardName(name_id, string(id_nickname))
					}
				}
				infocard_addition.WriteLineStr(id_name, "= ", strconv.FormatFloat(player_bonus.Bonus.Get(), 'f', 2, 64))
			}
			infocard_addition.WriteLineStr("")
		}
	}

	e.PutInfocard(infocarder.InfocardKey(base.Nickname), append(sb.Lines, infocard_addition.Lines...))

	return base
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
