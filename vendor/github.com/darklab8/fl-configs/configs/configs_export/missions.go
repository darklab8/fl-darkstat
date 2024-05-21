package configs_export

import (
	"fmt"
	"math"
	"sort"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
)

type MissioNFaction struct {
	FactionName     string
	FactionNickname string
	MinDifficulty   float64
	MaxDifficulty   float64
	Weight          int

	Infocard InfocardKey

	MoneyAward int
	NpcRanks   []int
}

type BaseMissions struct {
	MinOffers    int
	MaxOffers    int
	Factions     []MissioNFaction
	ShipRanksMap map[int]bool
	ShipRanks    []int

	AvgMoneyAward int
	MaxMoneyAward int
}

type DiffToMoney struct {
	MinLevel   float64
	MoneyAward int
}

func (e *Exporter) GetMissions(bases []Base, factions []Faction) []Base {

	var factions_map map[string]Faction = make(map[string]Faction)
	for _, faction := range factions {
		factions_map[faction.Nickname] = faction
	}

	var diffs_to_money []DiffToMoney
	for _, diff_to_money := range e.configs.DiffToMoney.DiffToMoney {
		diffs_to_money = append(diffs_to_money, DiffToMoney{
			MinLevel:   diff_to_money.MinLevel.Get(),
			MoneyAward: diff_to_money.MoneyAward.Get(),
		})
	}
	sort.Slice(diffs_to_money, func(i, j int) bool {
		return diffs_to_money[i].MinLevel > diffs_to_money[j].MinLevel
	})

	for base_index, base := range bases {
		base.Missions.ShipRanksMap = make(map[int]bool)

		base_info, ok := e.configs.MBases.BaseMap[base.Nickname]
		if !ok {
			continue
		}

		// Firstly finding SystemBase coresponding to base
		system, system_exists := e.configs.Systems.SystemsMap[base.SystemNickname]
		if !system_exists {
			continue
		}

		var system_base *systems_mapped.Base
		for _, sys_base := range system.Bases {
			if sys_base.IdsName.Get() == base.StridName {
				system_base = sys_base
				break
			}
		}
		if system_base == nil {
			continue
		}

		// Verify that base vignette fields exist in 30k around of it, otherwise base is not able to start missions
		base_has_vignettes := false
		vignette_valid_base_mission_range := float64(30000)
		for _, vignette := range system.MissionZoneVignettes {

			if _, ok := system_base.Pos.Y.GetValue(); !ok {
				fmt.Println()
			}
			if _, ok := vignette.Pos.Y.GetValue(); !ok {
				fmt.Println()
				if _, ok := vignette.Pos.Y.GetValue(); !ok {
				}
			}

			x_dist := math.Pow((system_base.Pos.X.Get() - vignette.Pos.X.Get()), 2)
			y_dist := math.Pow((system_base.Pos.Y.Get() - vignette.Pos.Y.Get()), 2)
			z_dist := math.Pow((system_base.Pos.Z.Get() - vignette.Pos.Z.Get()), 2)
			distance := math.Pow((x_dist + y_dist + z_dist), 0.5)

			if distance < vignette_valid_base_mission_range-float64(vignette.Size.Get()) {
				base_has_vignettes = true
				break
			}
		}

		if !base_has_vignettes {
			continue
		}

		for _, faction_info := range base_info.BaseFactions {
			faction := MissioNFaction{
				FactionNickname: faction_info.Faction.Get(),
			}
			faction.MinDifficulty, _ = faction_info.MissionType.MinDifficulty.GetValue()
			faction.MaxDifficulty, _ = faction_info.MissionType.MaxDifficulty.GetValue()

			if value, ok := faction_info.Weight.GetValue(); ok {
				faction.Weight = value
			} else {
				faction.Weight = 100
			}

			// Verify that faction has Spawn zones before adding
			// Otherwise skip
			spawn_zones, found_zones := system.MissionsSpawnZonesByFaction[faction.FactionNickname]
			if !found_zones {
				continue
			}
			_ = spawn_zones

			for _, diff_to_money := range diffs_to_money {
				if faction.MinDifficulty > diff_to_money.MinLevel {
					faction.MoneyAward = diff_to_money.MoneyAward
					break
				}
			}

			// NpcRank appropriate to current mission difficulty based on set range.
			for _, rank_to_diff := range e.configs.NpcRankToDiff.NPCRankToDifficulties {

				min_diff := rank_to_diff.Difficulties[0].Get()
				max_diff := rank_to_diff.Difficulties[len(rank_to_diff.Difficulties)-1].Get()

				if faction.MinDifficulty >= min_diff && faction.MinDifficulty <= max_diff {
					faction.NpcRanks = append(faction.NpcRanks, rank_to_diff.Rank.Get())
					continue
				}
				if faction.MaxDifficulty >= min_diff && faction.MaxDifficulty <= max_diff {
					faction.NpcRanks = append(faction.NpcRanks, rank_to_diff.Rank.Get())
					continue
				}

			}

			for _, npc_rank := range faction.NpcRanks {
				base.Missions.ShipRanksMap[npc_rank] = true
			}

			if faction_info, ok := factions_map[faction.FactionNickname]; ok {
				faction.Infocard = faction_info.Infocard
				faction.FactionName = faction_info.Name

			}

			// Find if the factions has defined NPCs in faction_prop.ini
			// Which have the necessary NPCRank
			// if not defined, then skip this faction

			// Add class of ships that will be encountered during missions of this faction

			base.Missions.Factions = append(base.Missions.Factions, faction)
		}

		// Make sanity check that Factions were added to base
		// If not then don't add to it mission existence.
		if len(base.Missions.Factions) == 0 {
			continue
		}

		if base_info.MVendor != nil {
			base.Missions.MinOffers, _ = base_info.MVendor.MinOffers.GetValue()
			base.Missions.MaxOffers, _ = base_info.MVendor.MaxOffers.GetValue()
		}

		// add unique found ship categories from factions to Missions overview
		for key, _ := range base.Missions.ShipRanksMap {
			base.Missions.ShipRanks = append(base.Missions.ShipRanks, key)
		}
		sort.Ints(base.Missions.ShipRanks)

		// Calculated base awardness
		for _, faction := range base.Missions.Factions {
			base.Missions.AvgMoneyAward += faction.MoneyAward

			if faction.MoneyAward > base.Missions.MaxMoneyAward {
				base.Missions.MaxMoneyAward = faction.MoneyAward
			}
		}
		base.Missions.AvgMoneyAward = base.Missions.AvgMoneyAward / len(base.Missions.Factions)
		bases[base_index] = base
	}

	return bases
}
