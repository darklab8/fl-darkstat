package configs_export

import (
	"sort"
	"strconv"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
)

type EnemyFaction struct {
	Faction
	NpcExist bool
}

/*
Calculates for enemy faction percentage of ships defined in faction_props/npcships.ini
If they aren't defined, Freelancer will be showing corrupted no missions when they encounter.
*/
func (e *Exporter) NewEnemyFaction(faction Faction, npc_ranks []int) EnemyFaction {
	var npc_ranks_need map[int]bool = make(map[int]bool)
	for _, rank := range npc_ranks {
		npc_ranks_need[rank] = true
	}

	var npc_ranks_exist map[int]bool = make(map[int]bool)

	result := EnemyFaction{
		Faction: faction,
	}

	faction_prop, prop_exists := e.Configs.FactionProps.FactionPropMapByNickname[faction.Nickname]

	if !prop_exists {
		return result
	}

	for _, npc_ship := range faction_prop.NpcShips {
		npc_ship_nickname := npc_ship.Get()
		if npc_shiparch, ok := e.Configs.NpcShips.NpcShipsByNickname[npc_ship_nickname]; ok {

			has_class_fighter := false
			for _, npc_class := range npc_shiparch.NpcClass {
				if npc_class.Get() == "class_fighter" {
					has_class_fighter = true
					break
				}
			}
			if !has_class_fighter {
				continue
			}
			str_level := npc_shiparch.Level.Get()
			if level, err := strconv.Atoi(str_level[1:]); err == nil {

				if _, ok := npc_ranks_need[level]; ok {
					npc_ranks_exist[level] = true
				}
			}
		}
	}

	result.NpcExist = len(npc_ranks_exist) > 0
	return result
}

type MissioNFaction struct {
	FactionName     string
	FactionNickname string
	MinDifficulty   float64
	MaxDifficulty   float64
	Weight          int
	Infocard        InfocardKey

	MinAward int
	MaxAward int
	NpcRanks []int
	Enemies  []EnemyFaction
	Err      cfg.Err
}

type BaseMissions struct {
	MinOffers         int
	MaxOffers         int
	Factions          []MissioNFaction
	NpcRanksAtBaseMap map[int]bool
	NpcRanksAtBase    []int

	EnemiesAtBaseMap map[string]EnemyFaction

	MinMoneyAward int
	MaxMoneyAward int
	Vignettes     int
	Err           cfg.Err
}

type DiffToMoney struct {
	MinLevel   float64
	MoneyAward int
}

func (e *Exporter) GetMissions(bases []*Base, factions []Faction) []*Base {

	var factions_map map[string]Faction = make(map[string]Faction)
	for _, faction := range factions {
		factions_map[faction.Nickname] = faction
	}

	var diffs_to_money []DiffToMoney
	for _, diff_to_money := range e.Configs.DiffToMoney.DiffToMoney {
		diffs_to_money = append(diffs_to_money, DiffToMoney{
			MinLevel:   diff_to_money.MinLevel.Get(),
			MoneyAward: diff_to_money.MoneyAward.Get(),
		})
	}
	sort.Slice(diffs_to_money, func(i, j int) bool {
		return diffs_to_money[i].MinLevel > diffs_to_money[j].MinLevel
	})

	for base_index, base := range bases {
		base.Missions.NpcRanksAtBaseMap = make(map[int]bool)
		base.Missions.EnemiesAtBaseMap = make(map[string]EnemyFaction)

		base_info, ok := e.Configs.MBases.BaseMap[base.Nickname]
		if !ok {
			base.Missions.Err = cfg.NewErr("base is not defined in mbases")
			bases[base_index] = base
			continue
		}

		if universe_base, ok := e.Configs.Universe.BasesMap[universe_mapped.BaseNickname(base.Nickname)]; ok {

			_, bar_exists := universe_base.ConfigBase.RoomMapByRoomNickname["bar"]
			if !bar_exists {
				base.Missions.Err = cfg.NewErr("bar is not defined for the base")
				bases[base_index] = base
				continue
			}
		}

		// Firstly finding SystemBase coresponding to base
		system, system_exists := e.Configs.Systems.SystemsMap[base.SystemNickname]
		if !system_exists {
			base.Missions.Err = cfg.NewErr("system is not found for base")
			bases[base_index] = base
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
			base.Missions.Err = cfg.NewErr("base is not found in system")
			bases[base_index] = base
			continue
		}

		// Verify that base vignette fields exist in 30k around of it, otherwise base is not able to start missions

		vignette_valid_base_mission_range := float64(30000)
		for _, vignette := range system.MissionZoneVignettes {
			distance, dist_err := DistanceForVecs(system_base.Pos, vignette.Pos)
			if dist_err != nil {
				continue
			}

			if distance < vignette_valid_base_mission_range+float64(vignette.Size.Get()) {
				base.Missions.Vignettes += 1

			}
		}

		if base.Missions.Vignettes == 0 {
			base.Missions.Err = cfg.NewErr("base has no vignette zones")
			bases[base_index] = base
			continue
		}

		if base_info.MVendor == nil {
			base.Missions.Err = cfg.NewErr("no mvendor in mbase")
			bases[base_index] = base
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

			faction_export_info, faction_exists := factions_map[faction.FactionNickname]
			if !faction_exists {
				faction.Err = cfg.NewErr("mission faction does not eixst")
				base.Missions.Factions = append(base.Missions.Factions, faction)
				continue
			}

			faction.Infocard = faction_export_info.InfocardKey
			faction.FactionName = faction_export_info.Name

			_, gives_missions := faction_info.MissionType.MinDifficulty.GetValue()
			if !gives_missions {
				faction.Err = cfg.NewErr("mission_type is not in mbase")
				base.Missions.Factions = append(base.Missions.Factions, faction)
				continue
			}

			for money_index, diff_to_money := range diffs_to_money {

				if money_index == 0 {
					continue
				}

				if faction.MinDifficulty >= diff_to_money.MinLevel {
					diff_range := diffs_to_money[money_index-1].MinLevel - diff_to_money.MinLevel
					bonus_range := faction.MinDifficulty - diff_to_money.MinLevel
					bonus_money_percentage := bonus_range / diff_range
					bonus_money := int(float64(diffs_to_money[money_index-1].MoneyAward-diff_to_money.MoneyAward) * bonus_money_percentage)
					faction.MinAward = diff_to_money.MoneyAward + bonus_money
				}

				if faction.MaxDifficulty >= diff_to_money.MinLevel && faction.MaxAward == 0 {
					diff_range := diffs_to_money[money_index-1].MinLevel - diff_to_money.MinLevel
					bonus_range := faction.MaxDifficulty - diff_to_money.MinLevel
					bonus_money_percentage := bonus_range / diff_range
					bonus_money := int(float64(diffs_to_money[money_index-1].MoneyAward-diff_to_money.MoneyAward) * bonus_money_percentage)
					faction.MaxAward = diff_to_money.MoneyAward + bonus_money
				}
			}

			// NpcRank appropriate to current mission difficulty based on set range.
			for _, rank_to_diff := range e.Configs.NpcRankToDiff.NPCRankToDifficulties {

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

			// Find if enemy npc spawn zones are intersecting with Vignettes.
			// They will be all the enemies for the faction.
			var target_reputation_by_faction map[string]Reputation = make(map[string]Reputation)
			for _, reputation := range faction_export_info.Reputations {
				target_reputation_by_faction[reputation.Nickname] = reputation
			}
			var base_enemies map[string]Faction = make(map[string]Faction)
			for _, npc_spawn_zone := range system.MissionsSpawnZone {

				var enemies []*systems_mapped.Patrol = make([]*systems_mapped.Patrol, 0, len(npc_spawn_zone.Factions))
				for _, potential_enemy := range npc_spawn_zone.Factions {
					potential_enemy_nickname, _ := potential_enemy.FactionNickname.GetValue()
					potential_enemy_rep, rep_exists := target_reputation_by_faction[potential_enemy_nickname]
					if !rep_exists {
						continue
					}

					if potential_enemy_rep.Rep <= -(0.3 - 0.001) {
						enemies = append(enemies, potential_enemy)
					}
				}

				if len(enemies) == 0 {
					continue
				}

				// EXPERIMENTAL: Turned off check for vignette check to be within npc spawning zones.
				// looks to be not necessary
				// if !IsAnyVignetteWithinNPCSpawnRange(system, npc_spawn_zone) {
				// 	continue
				// }

				for _, enemy := range enemies {
					faction_enemy, faction_found := factions_map[enemy.FactionNickname.Get()]
					if !faction_found {
						continue
					}
					copy_enemy := faction_enemy
					base_enemies[faction_enemy.Nickname] = copy_enemy
				}
			}

			if len(base_enemies) == 0 {
				faction.Err = cfg.NewErr("no npc spawn zones with enemies")
				base.Missions.Factions = append(base.Missions.Factions, faction)
				continue
			}
			for _, enemy_faction := range base_enemies {
				faction.Enemies = append(faction.Enemies, e.NewEnemyFaction(enemy_faction, faction.NpcRanks))
			}

			base.Missions.Factions = append(base.Missions.Factions, faction)
		}

		// Make sanity check that Factions were added to base
		// If not then don't add to it mission existence.
		if len(base.Missions.Factions) == 0 {
			base.Missions.Err = cfg.NewErr("no msn giving factions found")
			bases[base_index] = base
			continue
		}

		npc_exist := false
		for _, faction := range base.Missions.Factions {
			if faction.Err != nil {
				continue
			}
			for _, enemy_faction := range faction.Enemies {
				if enemy_faction.NpcExist {
					npc_exist = true
				}
			}
		}
		if !npc_exist {
			base.Missions.Err = cfg.NewErr("npcs do not exist")
			bases[base_index] = base
			continue
		}

		if base_info.MVendor != nil {
			base.Missions.MinOffers, _ = base_info.MVendor.MinOffers.GetValue()
			base.Missions.MaxOffers, _ = base_info.MVendor.MaxOffers.GetValue()

			if base.Missions.Vignettes < base.Missions.MinOffers {
				base.Missions.MinOffers = base.Missions.Vignettes
			}
			if base.Missions.Vignettes < base.Missions.MaxOffers {
				base.Missions.MaxOffers = base.Missions.Vignettes
			}
		}

		// summarization for base
		for fc_index, faction := range base.Missions.Factions {
			if faction.Err != nil {
				faction.MaxAward = 0
				faction.MinAward = 0
				base.Missions.Factions[fc_index] = faction
				continue
			}

			for _, npc_rank := range faction.NpcRanks {
				base.Missions.NpcRanksAtBaseMap[npc_rank] = true
			}

			for _, enemy_faction := range faction.Enemies {
				base.Missions.EnemiesAtBaseMap[enemy_faction.Nickname] = enemy_faction
			}

			if faction.MinAward < base.Missions.MinMoneyAward || base.Missions.MinMoneyAward == 0 {
				base.Missions.MinMoneyAward = faction.MinAward
			}

			if faction.MaxAward > base.Missions.MaxMoneyAward {
				base.Missions.MaxMoneyAward = faction.MaxAward
			}
		}

		// add unique found ship categories from factions to Missions overview
		for key := range base.Missions.NpcRanksAtBaseMap {
			base.Missions.NpcRanksAtBase = append(base.Missions.NpcRanksAtBase, key)
		}
		sort.Ints(base.Missions.NpcRanksAtBase)

		bases[base_index] = base
	}

	return bases
}
