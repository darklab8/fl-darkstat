package configs_export

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/go-utils/utils/ptr"
)

type MiningInfo struct {
	DynamicLootMin        int
	DynamicLootMax        int
	DynamicLootDifficulty int
	MinedGood             MarketGood
}

func (e *Exporter) GetOres(Commodities []*Commodity) []*Base {
	var bases []*Base

	var comm_by_nick map[string]*Commodity = make(map[string]*Commodity)
	for _, comm := range Commodities {
		comm_by_nick[comm.Nickname] = comm
	}

	for _, system := range e.configs.Systems.Systems {

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
				Pos:        location,
				MiningInfo: MiningInfo{},
			}
			base.DynamicLootMin, _ = asteroids.LootableZone.DynamicLootMin.GetValue()
			base.DynamicLootMax, _ = asteroids.LootableZone.DynamicLootMax.GetValue()
			base.DynamicLootDifficulty, _ = asteroids.LootableZone.DynamicLootDifficulty.GetValue()

			base.Nickname, _ = zone.Nickname.GetValue()
			base.InfocardID, _ = zone.IDsInfo.GetValue()
			base.StridName, _ = zone.IdsName.GetValue()

			base.Infocard = InfocardKey(base.Nickname)

			base.Archetypes = append(base.Archetypes, "mining_operation")
			base.FactionName = "neutral"

			base.SystemNickname = system.Nickname
			if system, ok := e.configs.Universe_config.SystemMap[universe_mapped.SystemNickname(base.SystemNickname)]; ok {
				base.System = e.GetInfocardName(system.Strid_name.Get(), base.SystemNickname)
				base.Region = e.GetRegionName(system)
			}

			fmt.Println("commodity=", commodity)

			market_good := MarketGood{
				Nickname:      commodity,
				BaseSells:     true,
				PriceModifier: 0,
				PriceBase:     0,
				PriceToBuy:    0,
				PriceToSell:   ptr.Ptr(0),
				Type:          "commodity",
			}

			if equipment, ok := e.configs.Equip.CommoditiesMap[commodity]; ok {
				market_good.Name = e.GetInfocardName(equipment.IdsName.Get(), market_good.Nickname)
			}
			base.Name = market_good.Name
			base.MarketGoods = append(base.MarketGoods, market_good)
			base.MinedGood = market_good

			added_goods[market_good.Nickname] = true

			if commodity, ok := comm_by_nick[market_good.Nickname]; ok {
				good_at_base := &GoodAtBase{
					BaseNickname:      base.Nickname,
					BaseName:          base.Name,
					BaseSells:         true,
					PriceBaseBuysFor:  0,
					PriceBaseSellsFor: 0,
					Volume:            commodity.Volume,

					SystemName: base.System,
					BasePos:    base.Pos,
					Region:     base.Region,
				}
				commodity.Bases = append(commodity.Bases, good_at_base)
			}

			if e.configs.Discovery != nil {

				if recipes, ok := e.configs.Discovery.BaseRecipeItems.RecipePerConsumed[market_good.Nickname]; ok {

					for _, recipe := range recipes {
						recipe_produces_only_commodities := true

						for _, produced := range recipe.ProcucedItem {

							_, is_commodity := e.configs.Equip.CommoditiesMap[produced.Get()]
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
								market_good := MarketGood{
									Nickname:      commodity_produced,
									BaseSells:     true,
									PriceModifier: 0,
									PriceBase:     0,
									PriceToBuy:    0,
									PriceToSell:   ptr.Ptr(0),
									Type:          "commodity",
								}
								if equipment, ok := e.configs.Equip.CommoditiesMap[commodity]; ok {
									market_good.Name = e.GetInfocardName(equipment.IdsName.Get(), market_good.Nickname)
								}
								base.MarketGoods = append(base.MarketGoods, market_good)

								if commodity, ok := comm_by_nick[market_good.Nickname]; ok {
									good_at_base := &GoodAtBase{
										BaseNickname:      base.Nickname,
										BaseName:          base.Name,
										BaseSells:         true,
										PriceBaseBuysFor:  0,
										PriceBaseSellsFor: 0,
										Volume:            commodity.Volume,

										SystemName: base.System,
										BasePos:    base.Pos,
										Region:     base.Region,
									}
									commodity.Bases = append(commodity.Bases, good_at_base)
								}
								added_goods[commodity_produced] = true
							}
						}
					}

				}
			}

			var sb []string
			sb = append(sb, base.Name)
			sb = append(sb, `This is is not a base.
It is a mining field with droppable ores`)
			sb = append(sb, "")
			sb = append(sb, "Trade routes shown do not account for a time it takes to mine those ores.")

			if e.configs.Discovery != nil {
				sb = append(sb, "")
				sb = append(sb, `<a href="https://discoverygc.com/wiki2/Mining">Check mining tutorial</a> to see how they can be mined`)

				sb = append(sb, "")
				sb = append(sb, `NOTE:
for Freelancer Discovery we also add possible sub products of refinery at player bases to possible trade routes from mining field.
				`)
			}

			sb = append(sb, "")
			sb = append(sb, "commodities:")
			for _, good := range base.MarketGoods {
				if good.Nickname == base.MinedGood.Nickname {
					sb = append(sb, fmt.Sprintf("Minable: %s (%s)", good.Name, good.Nickname))
				} else {
					sb = append(sb, fmt.Sprintf("Refined at POB: %s (%s)", good.Name, good.Nickname))
				}
			}

			e.Infocards[InfocardKey(base.Nickname)] = sb

			if base.Nickname == "zone_br05_gold_dummy_field" {
				fmt.Println()
			}

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
