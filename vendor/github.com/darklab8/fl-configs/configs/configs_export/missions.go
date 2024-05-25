package configs_export

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/darklab8/fl-configs/configs/config_consts"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-configs/configs/settings/logus"
)

type MissioNFaction struct {
	FactionName     string
	FactionNickname string
	MinDifficulty   float64
	MaxDifficulty   float64
	Weight          int

	Infocard InfocardKey

	MinAward int
	MaxAward int
	NpcRanks []int
	Enemies  []Faction
	Err      error
}

type BaseMissions struct {
	MinOffers         int
	MaxOffers         int
	Factions          []MissioNFaction
	NpcRanksAtBaseMap map[int]bool
	NpcRanksAtBase    []int

	EnemiesAtBaseMap map[string]Faction

	MinMoneyAward int
	MaxMoneyAward int
	Err           error
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
		base.Missions.NpcRanksAtBaseMap = make(map[int]bool)
		base.Missions.EnemiesAtBaseMap = make(map[string]Faction)

		base_info, ok := e.configs.MBases.BaseMap[base.Nickname]
		if !ok {
			base.Missions.Err = errors.New("base is not defined in mbases")
			bases[base_index] = base
			continue
		}

		// Firstly finding SystemBase coresponding to base
		system, system_exists := e.configs.Systems.SystemsMap[base.SystemNickname]
		if !system_exists {
			base.Missions.Err = errors.New("system is not found for base")
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
			base.Missions.Err = errors.New("base is not found in system")
			bases[base_index] = base
			continue
		}

		// Verify that base vignette fields exist in 30k around of it, otherwise base is not able to start missions
		base_has_vignettes := false
		vignette_valid_base_mission_range := float64(30000)
		for _, vignette := range system.MissionZoneVignettes {
			distance, dist_err := DistanceForVecs(system_base.Pos, vignette.Pos)
			if dist_err != nil {
				continue
			}

			if distance < vignette_valid_base_mission_range-float64(vignette.Size.Get()) {
				base_has_vignettes = true
				break
			}
		}

		if !base_has_vignettes {
			base.Missions.Err = errors.New("base has no vignette zones")
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
				faction.Err = errors.New("mission faction does not eixst")
				base.Missions.Factions = append(base.Missions.Factions, faction)
				continue
			}

			faction.Infocard = faction_export_info.Infocard
			faction.FactionName = faction_export_info.Name

			_, gives_missions := faction_info.MissionType.MinDifficulty.GetValue()
			if !gives_missions {
				faction.Err = errors.New("mission_type is not in mbase")
				base.Missions.Factions = append(base.Missions.Factions, faction)
				continue
			}

			// Verify that faction has Spawn zones before adding
			// Otherwise skip
			spawn_zones, found_zones := system.MissionsSpawnZonesByFaction[faction.FactionNickname]
			if !found_zones {
				faction.Err = errors.New("no mission faction npc spawning zones")
				base.Missions.Factions = append(base.Missions.Factions, faction)
				continue
			}
			_ = spawn_zones

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
					relationship_status := config_consts.GetRelationshipStatus(potential_enemy_rep.Rep)
					if relationship_status == config_consts.RepEnemy {
						enemies = append(enemies, potential_enemy)
					}
				}

				if len(enemies) == 0 {
					continue
				}

				matched_vignette := false
				for _, vignette := range system.MissionZoneVignettes {

					distance, dist_err := DistanceForVecs(system_base.Pos, npc_spawn_zone.Pos)
					if dist_err != nil {
						continue
					}

					max_spwn_zone_size, err_max_size := GetMaxRadius(npc_spawn_zone.Size)
					logus.Log.CheckWarn(err_max_size, "expected finding max size, but object does not have it")

					if distance < float64(vignette.Size.Get())+max_spwn_zone_size {
						matched_vignette = true
						break
					}
				}

				if !matched_vignette {
					continue
				}

				for _, enemy := range enemies {
					faction_enemy, faction_found := factions_map[enemy.FactionNickname.Get()]
					if !faction_found {
						continue
					}
					copy_enemy := faction_enemy
					base_enemies[faction_enemy.Nickname] = copy_enemy
				}
			}

			// Find if the factions has defined NPCs in faction_prop.ini
			// Which have the necessary NPCRank
			// if not defined, then skip this faction
			// Add class of ships that will be encountered during missions of this faction
			// for _, enemy_faction := range base_enemies {
			// }

			if len(base_enemies) == 0 {
				faction.Err = errors.New("no enemy npc spawns near vignettes")
				base.Missions.Factions = append(base.Missions.Factions, faction)
				continue
			}
			for _, enemy_faction := range base_enemies {
				faction.Enemies = append(faction.Enemies, enemy_faction)
			}

			base.Missions.Factions = append(base.Missions.Factions, faction)
		}

		// Make sanity check that Factions were added to base
		// If not then don't add to it mission existence.
		if len(base.Missions.Factions) == 0 {
			base.Missions.Err = errors.New("no msn giving factions found")
			bases[base_index] = base
			continue
		}

		if base_info.MVendor != nil {
			base.Missions.MinOffers, _ = base_info.MVendor.MinOffers.GetValue()
			base.Missions.MaxOffers, _ = base_info.MVendor.MaxOffers.GetValue()
		}

		if strings.Contains(base.Name, "Essex") {
			fmt.Println()
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
		for key, _ := range base.Missions.NpcRanksAtBaseMap {
			base.Missions.NpcRanksAtBase = append(base.Missions.NpcRanksAtBase, key)
		}
		sort.Ints(base.Missions.NpcRanksAtBase)

		bases[base_index] = base
	}

	return bases
}

func DistanceForVecs(Pos1 *semantic.Vect, Pos2 *semantic.Vect) (float64, error) {
	if _, ok := Pos1.X.GetValue(); !ok {
		return 0, errors.New("no x")
	}
	if _, ok := Pos2.X.GetValue(); !ok {
		return 0, errors.New("no x")
	}

	x_dist := math.Pow((Pos1.X.Get() - Pos2.X.Get()), 2)
	y_dist := math.Pow((Pos1.Y.Get() - Pos2.Y.Get()), 2)
	z_dist := math.Pow((Pos1.Z.Get() - Pos2.Z.Get()), 2)
	distance := math.Pow((x_dist + y_dist + z_dist), 0.5)
	return distance, nil
}

func GetMaxRadius(Size *semantic.Vect) (float64, error) {
	max_size := 0.0
	if value, ok := Size.X.GetValue(); ok {
		if value > max_size {
			max_size = value
		}
	}
	if value, ok := Size.Y.GetValue(); ok {
		if value > max_size {
			max_size = value
		}
	}
	if value, ok := Size.Z.GetValue(); ok {
		if value > max_size {
			max_size = value
		}
	}
	if max_size == 0 {
		return 0, errors.New("not found size")
	}

	return max_size, nil
}
