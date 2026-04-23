package configs_export

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/fl-darkstat/darkstat/settings/lootable_shown"
	"github.com/google/uuid"
)

type LootKind int8

const (
	LootUnknown LootKind = iota
	LootWreck
	LootFLSRSolar
	LootFLSRNPC
	LootEncounter
	LootNotFoundNpcArch
	LootPhantom
	// LootPhantomLoop // lets not seek it
)

func (l LootKind) ToStr() string {
	if l == LootWreck {
		return "wreck"
	}
	if l == LootFLSRSolar {
		return "flsr_solar"
	}
	if l == LootFLSRNPC {
		return "flsr_npc"
	}
	if l == LootEncounter {
		return "encounter"
	}

	if l == LootNotFoundNpcArch {
		return "err_enc"
	}
	if l == LootPhantom {
		return "global"
	}

	return "unknown"
}

const (
	LootMaxWrecks     = 4
	LootMaxFLSR       = 2
	LootMaxEncounters = 2

	LootMaxNotFoundNPCDrops = 1
	LootMaxPhantom          = 1
)

func (e *Exporter) SetPermitted(permitted_wrecks map[string]bool, permitted_encounters map[string]bool, loot_info *LootInfo) {
	if e.Mapped.FLSR != nil {
		loot_info.Permitted = true
	} else if loot_info.Kind == LootNotFoundNpcArch {
		loot_info.Permitted = true
	} else if loot_info.Kind == LootEncounter {
		loot_info.Permitted = true

	} else if loot_info.Kind == LootWreck {
		loot_info.Permitted = permitted_wrecks[loot_info.Nickname]
	}

	if !loot_info.Permitted {
		loot_info.SystemName = ""
		loot_info.SectorCoord = ""
		loot_info.Pos = cfg.Vector{}
		loot_info.PlaceNick = ""
		loot_info.LootSource = LootSourceUnknown
	}
}

// if is_lootable_by_plugin is true, it means lootable by disco plugin MarketController
func (e *Exporter) IsLootable(item_nickname string, loot_source LootSource) (is_lootable bool, is_lootable_by_plugin bool) {
	item, found_item := e.Mapped.Equip().ItemsMap[item_nickname]
	if !found_item {
		return false, false
	}

	if loot_source == LootSourceCargo || loot_source == LootSourceAny {
		if _, ok := item.DropChanceNpcUnmounted.GetValue(); ok {
			is_lootable = true
			is_lootable_by_plugin = true
		}
	}
	if loot_source == LootSourceEquip || loot_source == LootSourceAny {
		if _, ok := item.DropChanceNpcMounted.GetValue(); ok {
			is_lootable = true
			is_lootable_by_plugin = true
		}
	}

	if is_lootable_value, ok := item.Lootable.GetValue(); ok {
		if is_lootable_value {
			is_lootable = true
		}
	}
	return is_lootable, is_lootable_by_plugin
}

type Wreck struct {
	LoadoutNickname string
	Archetype       string
	Nickname        string
	Pos             cfg.Vector
	Kind            LootKind
	Event           string
}

func (e *Exporter) ProcessWreck(wreck Wreck, system *systems_mapped.System) ([]*LootInfo, error) {
	var results []*LootInfo
	loadout_nickname := wreck.LoadoutNickname

	if loadout, ok := e.Mapped.Loadouts.LoadoutsByNick[loadout_nickname]; ok {
		type Item struct {
			nickname    string
			loot_source LootSource
		}

		// query here archetype of wreck ? by archetype of wreck, get Solar

		solar := e.Mapped.Solararch.SolarsByNick[wreck.Archetype]
		if is_destructible, _ := solar.Destructible.GetValue(); !is_destructible {
			return nil, errors.New("not destructable")
		}

		allowed_cargo := false
		allowed_hardpoints := make(map[string]bool)
		for _, fuse := range solar.Fuses {
			if fuse, ok := e.Mapped.Fuses.FuseMap[fuse.Get()]; ok {
				if fuse.DoesDropCargo {
					allowed_cargo = true
				}
				for key, _ := range fuse.LootableHardpoints {
					allowed_hardpoints[key] = true
				}
			}
		}
		for _, fuse := range solar.Fuses {
			if fuse, ok := e.Mapped.Fuses.FuseMap[fuse.Get()]; ok {
				for key, _ := range fuse.NotLootableHardpoints {
					delete(allowed_hardpoints, key)
				}
			}
		}

		var item_nicknames []Item
		if allowed_cargo {
			for _, cargo := range loadout.Cargos {
				item_nicknames = append(item_nicknames, Item{
					nickname:    cargo.Nickname.Get(),
					loot_source: LootSourceCargo,
				})
			}
		}

		for _, equip := range loadout.Equips {

			item_nickname := equip.Nickname.Get()

			hardpoint, found_hardpoint := equip.Hardpoint.GetValue()
			if !found_hardpoint {
				continue
			}

			_, permitted_hardpoint := allowed_hardpoints[hardpoint]
			if !permitted_hardpoint {
				continue
			}

			item_nicknames = append(item_nicknames, Item{
				nickname:    item_nickname,
				loot_source: LootSourceEquip,
			})
		}

		for _, item := range item_nicknames {

			is_lootable, is_lootable_by_disco := e.IsLootable(item.nickname, item.loot_source)
			if !is_lootable {
				continue
			}

			loot_info := &LootInfo{
				Nickname:   item.nickname,
				Kind:       wreck.Kind,
				LootSource: LootSourceEquip,
				PlaceNick:  wreck.Nickname,
			}

			system_uni := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(system.Nickname)]
			loot_info.Pos = wreck.Pos
			loot_info.SectorCoord = VectorToSectorCoord(system_uni, loot_info.Pos)
			loot_info.SystemName = e.GetInfocardName(system_uni.StridName.Get(), system.Nickname)
			loot_info.SystemNickname = system.Nickname
			loot_info.ObjNickname = wreck.Nickname

			// TODO [ ] missing to validate dump_cargo and destroy_hp_attachment at Solar archetype fuses
			is_fuse_allowed := true
			if !is_fuse_allowed && !is_lootable_by_disco {
				continue
			}

			results = append(results, loot_info)
		}
	}

	if len(results) == 0 {
		return nil, errors.New("not found loot")
	}

	return results, nil
}

func (e *Exporter) FindableInLoot() (map[string]bool, []*LootInfo) {

	var permitted_wrecks map[string]bool

	if e.Mapped.Discovery != nil {
		permitted_wrecks = lootable_shown.GetDiscoveryWrecksAllowed()
	} else if e.Mapped.FLSR != nil {
		permitted_wrecks = lootable_shown.GetFLSRWrecksAllowed()
	} else {
		permitted_wrecks = lootable_shown.GetVanillaWrecksAlloed()
	}

	var permitted_encounters map[string]bool

	if e.Mapped.Discovery != nil {
		permitted_encounters = lootable_shown.GetDiscoveryEncountersAllowed()
	} else if e.Mapped.FLSR != nil {
		permitted_encounters = lootable_shown.GetFLSEncountersAllowed()
	} else {
		permitted_encounters = lootable_shown.GetVanillaEncountersAlloed()
	}

	var loots []*LootInfo

	if e.findable_in_loot_cache != nil {
		return e.findable_in_loot_cache, e.findable_wrecks
	}

	var flsr_loot []*LootInfo
	var loot_event_sources_by_loot_nickname map[string]map[string]bool = make(map[string]map[string]bool)
	findable_limit_wrecks := make(map[string]int)
	e.findable_in_loot_cache = make(map[string]bool)
	{
		// Purely wrecks stuff
		// get cargo and equip out of "wrecks" ships
		// validate that wrecks have archetype at solar permitting  dump_cargo
		// and validate that fuse destroy_hp_attachment at solar has fate = loot for the equip hptype
		// [x] missing to validate dump_cargo and destroy_hp_attachment at Solar archetype fuses

		unique_wreck_loot := make(map[string]bool)

		process_wreck := func(wreck Wreck, system *systems_mapped.System) {

			wrecks, err := e.ProcessWreck(wreck, system)
			if err != nil {
				return
			}

			for _, loot_info := range wrecks {

				key_uniqueness := loot_info.Nickname + loot_info.SystemName + loot_info.SectorCoord
				if _, ok := unique_wreck_loot[key_uniqueness]; ok {
					continue
				}
				unique_wreck_loot[key_uniqueness] = true

				if loot_info.Kind == LootWreck {
					e.findable_in_loot_cache[loot_info.Nickname] = true
					findable_limit_wrecks[loot_info.Nickname] += 1
					if findable_limit_wrecks[loot_info.Nickname] <= LootMaxWrecks {
						e.SetPermitted(permitted_wrecks, permitted_encounters, loot_info)
						loots = append(loots, loot_info)
					}
				} else if loot_info.Kind == LootFLSRSolar {
					flsr_loot = append(flsr_loot, loot_info)
					if loot_event_sources_by_loot_nickname[loot_info.Nickname] == nil {
						loot_event_sources_by_loot_nickname[loot_info.Nickname] = make(map[string]bool)
					}
					loot_event_sources_by_loot_nickname[loot_info.Nickname][wreck.Event] = true
				} else {
					panic("not identified type of wreeck :)")
				}
			}
		}

		for _, system := range e.Mapped.Systems.Systems {
			for _, wreck := range system.Wrecks {
				process_wreck(Wreck{
					LoadoutNickname: wreck.Loadout.Get(),
					Archetype:       wreck.Archetype.Get(),
					Nickname:        wreck.Nickname.Get(),
					Pos:             wreck.Pos.Get(),
					Kind:            LootWreck,
				}, system)
			}
		}

		fmt.Println("mission parsing started zero")
		if e.Mapped.FLSR != nil {
			fmt.Println("mission parsing started")
			missions := e.Mapped.FLSR.FLSRMissions.Missions
			for _, mission := range missions {
				// is_active, _ := mission.InitState.GetValue()
				// if !is_active {
				// 	continue
				// }
				// fmt.Println("mission=", mission.Nickname.Get())
				for _, wreck := range mission.Solars {

					loadout, found_loadout := wreck.Loadout.GetValue()
					if !found_loadout {
						continue
					}

					system_nickname, found_system := wreck.System.GetValue()
					if !found_system {
						continue
					}

					system, found_sys := e.Mapped.Systems.SystemsMap[system_nickname]
					if !found_sys {
						continue
					}

					process_wreck(Wreck{
						LoadoutNickname: loadout,
						Archetype:       wreck.Archetype.Get(),
						Nickname:        wreck.Nickname.Get(),
						Pos:             wreck.Pos.Get(),
						Kind:            LootFLSRSolar,
						Event:           mission.Nickname.Get(),
					}, system)
				}
			}
			fmt.Println("mission parsing finished")

		}
	}

	// Now this is purely for encounter i think? because only they search in [Ship]
	// By NPC ship loadout, we find ship class architecture
	// then we need to find ship class
	// through which we validate it is in encounters
	// throughthat we find zones in which encounter is, and we get where it is dropped
	// to validate if "equip" in mLootProps

	// the loot will drop anyway regardless if it is in mlootprops
	// only checking forbidding fuses is important?
	mlootprops_allowed := make(map[string]bool)
	for _, MLootProp := range e.Mapped.LookProps.LootProps {
		equipment_nickname := MLootProp.Nickname.Get()
		is_lootable, _ := e.IsLootable(equipment_nickname, LootSourceAny)
		if !is_lootable {
			continue
		}

		mlootprops_allowed[equipment_nickname] = true
	}

	type NpcLoot struct {
		*LootInfo
		LoadoutNickname string
	}

	unique_loadouts := make(map[string]bool)
	GetNpcLoots := func(loadout_nickname string, loadout_archetype string) []*NpcLoot {
		var results []*NpcLoot
		if _, ok := unique_loadouts[loadout_nickname]; ok {
			return results
		}
		unique_loadouts[loadout_nickname] = true

		if loadout, ok := e.Mapped.Loadouts.LoadoutsByNick[loadout_nickname]; ok {
			shiparch := e.Mapped.Shiparch.ShipsMap[loadout_archetype]
			forbidden_hardpoints := make(map[string]bool)
			fuse_drops_equips := make(map[string]bool)

			for _, fuse := range shiparch.Fuses {
				if fuse, ok := e.Mapped.Fuses.FuseMap[fuse.Get()]; ok {
					for key, _ := range fuse.NotLootableHardpoints {
						forbidden_hardpoints[key] = true
					}

					for key, _ := range fuse.LootableHardpoints {
						fuse_drops_equips[key] = true
					}
				}
			}

			for _, cargo := range loadout.Cargos {
				item_nickname := cargo.Nickname.Get()

				_, is_plugin_lootable := e.IsLootable(item_nickname, LootSourceCargo)

				loot_info := &LootInfo{
					Kind:           LootNotFoundNpcArch,
					Nickname:       item_nickname,
					LootSource:     LootSourceCargo,
					DiscoEncounter: is_plugin_lootable,
				}
				loot_npc_drop := &NpcLoot{
					LootInfo:        loot_info,
					LoadoutNickname: loadout_nickname,
				}

				// loot_droppable := false
				// if _, ok := mlootprops_allowed[loot_info.Nickname]; ok {
				// 	loot_droppable = true
				// }
				// _ = loot_droppable

				// if !is_lootable && !loot_droppable && !fuse_cargo_drop && !is_plugin_lootable {
				// 	continue
				// }

				results = append(results, loot_npc_drop)
			}
			for _, equip := range loadout.Equips {
				item_nickname := equip.Nickname.Get()
				is_lootable, is_plugin_lootable := e.IsLootable(item_nickname, LootSourceEquip)

				if !is_lootable {
					continue
				}

				// skipping dissapearing hardpoint equips
				hardpoint, found_hardpoint := equip.Hardpoint.GetValue()
				if !found_hardpoint && !is_plugin_lootable {
					continue
				}
				_, is_forbidden_hardpoint := forbidden_hardpoints[hardpoint]
				if is_forbidden_hardpoint && !is_plugin_lootable {
					continue
				}

				loot_droppable := false
				if _, ok := mlootprops_allowed[item_nickname]; ok {
					loot_droppable = true
				}
				fuse_droppable := false
				if _, ok := fuse_drops_equips[hardpoint]; ok {
					fuse_droppable = true
				}

				if !loot_droppable && !fuse_droppable && !is_plugin_lootable {
					continue
				}

				loot_info := &LootInfo{
					Kind:           LootNotFoundNpcArch,
					Nickname:       item_nickname,
					LootSource:     LootSourceEquip,
					DiscoEncounter: is_plugin_lootable,
				}
				loot_npc_drop := &NpcLoot{
					LootInfo:        loot_info,
					LoadoutNickname: loadout_nickname,
				}

				results = append(results, loot_npc_drop)
			}
		}
		return results
	}

	// May time to make it iterator :)
	IteratorNpcShips := func(send_npc_loot func(npc_loot *NpcLoot)) {
		for _, npc_arch := range e.Mapped.NpcShips.NpcShips {
			loadout_nickname := npc_arch.Loadout.Get()

			loadout_archetype := npc_arch.ShipArchetype.Get()

			for _, loot_npc_drop := range GetNpcLoots(loadout_nickname, loadout_archetype) {
				send_npc_loot(loot_npc_drop)
			}
		}
	}

	// [ ] missing to validate dump_cargo and destroy_hp_attachment at Ship fuses
	// [ ] FIX to finding by faction=fc_c_grp, and there using encounter to validate only ship_class that will spawn
	// encounter = deep_bossencounter01, 19, 1 , patrol
	// faction = fc_c_grp, 1
	unique_encounter_loot := make(map[string]bool)
	max_encounter_by_nick := make(map[string]int)
	max_notfound_npc_drop := make(map[string]int)
	found_npc_loot_in_encounters := make(map[string]bool)
	var not_matched_loot []*NpcLoot
	IteratorNpcShips(func(npc_loot *NpcLoot) {
		if max_encounter_by_nick[npc_loot.LootInfo.Nickname] > LootMaxEncounters {
			return
		}
		// :x: by cargo find loadout name in [Loadout]
		// by louadout nickname, find in [NPCShipArch] ship class of it like `class_deep_cbcr` (npcships.ini)

		var ship_class_members []string
		var affiliations []string // like fc_n_grp, which we will be able to use valid Zones
		var ship_arch_nickname string
		if npcshiparch, ok := e.Mapped.NpcShips.NpcShipsByLoadout[npc_loot.LoadoutNickname]; ok {
			for _, ship_cl := range npcshiparch.ShipClass {
				ship_class_members = append(ship_class_members, ship_cl.Get())
			}

			ship_arch_nickname = npcshiparch.Nickname.Get()

			if faction_props, ok := e.Mapped.FactionProps.FactionPropMapByNpcShip[ship_arch_nickname]; ok {
				for _, faction_prop := range faction_props {
					affiliations = append(affiliations, faction_prop.Affiliation.Get())
				}
			}
		}

		var ship_class_nicknames map[string]bool = make(map[string]bool)
		ship_class_by_member := e.Mapped.ShipClasses.ShipClassByMember
		for _, ship_class_member := range ship_class_members {
			if ship_classes, ok := ship_class_by_member[ship_class_member]; ok {
				for _, shipclass := range ship_classes {
					ship_class_nicknames[shipclass.Nickname.Get()] = true
				}
			}
		}

		// by ship class nickname, find encounter formation filename in [EncounterFormation]
		// by encounter filename, find nickname of it in [EncounterParameters]
		// by encounter nickname, find in system [Object], relevant positions/system names and sector coords

		// Change to seek valid encounters by faction ?
		// find zones by encounter nickname
		found_zones := make(map[string]*systems_mapped.EncounterZoneInSystem)
		for _, afilliation := range affiliations {
			if encounters, ok := e.Mapped.Systems.EncounterByAfilliation[afilliation]; ok {
				for _, encounter := range encounters {
					// TODO Validate here that Encounter the zone has, has matching ShipClasses
					matched_ship_class_in_formation := false
					for _, encounter_name := range encounter.Zone.Encounters {
						encounter_param := encounter.System.EncounterParametersByName[encounter_name.Get()]
						encoutners_forms_by_filepath := e.Mapped.Systems.EncounterFormationByFilepath
						encounter_formation := encoutners_forms_by_filepath[encounter_param.Filename.Get()]

						for _, shipclass := range encounter_formation.ShipClasses {
							if _, ok := ship_class_nicknames[shipclass.Get()]; ok {
								matched_ship_class_in_formation = true
								break
							}
						}
						for _, ship_arch := range encounter_formation.ShipsByNpcArch {
							if ship_arch_nickname == ship_arch.Get() {
								matched_ship_class_in_formation = true
								break
							}
						}
					}
					if !matched_ship_class_in_formation {
						continue
					}

					found_zones[encounter.Zone.Nickname.Get()] = encounter
				}
			}
		}

		for _, zone := range found_zones {
			// YOU FOUND IT. ADD POS/SYSTEM NAME/Sector CORD

			density, found_density := zone.Zone.Density.GetValue()
			if !found_density {
				continue
			}
			if density == 0 {
				continue
			}

			item_nickname := npc_loot.LootInfo.Nickname
			_, is_plugin_lootable := e.IsLootable(item_nickname, npc_loot.LootSource)
			// if !is_lootable {
			// 	continue
			// }
			// _ = is_plugin_lootable
			loot_info := &LootInfo{
				Nickname:       item_nickname,
				Kind:           LootEncounter,
				LootSource:     npc_loot.LootSource,
				PlaceNick:      zone.Zone.Nickname.Get(),
				DiscoEncounter: is_plugin_lootable,
			}

			if _, ok := zone.Zone.Pos.GetValue(); ok {
				loot_info.Pos = zone.Zone.Pos.Get()
			}

			system_uni := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(zone.System.Nickname)]
			loot_info.SectorCoord = VectorToSectorCoord(system_uni, loot_info.Pos)
			loot_info.SystemName = e.GetInfocardName(system_uni.StridName.Get(), zone.System.Nickname)
			loot_info.SystemNickname = zone.System.Nickname
			loot_info.Kind = LootEncounter
			loot_info.ObjNickname = loot_info.PlaceNick

			key_uniqueness := loot_info.Nickname + zone.System.Nickname + loot_info.SectorCoord
			if _, ok := unique_encounter_loot[key_uniqueness]; ok {
				continue
			}
			unique_encounter_loot[key_uniqueness] = true

			if max_encounter_by_nick[loot_info.Nickname] > LootMaxEncounters && loot_info.Nickname != "commodity_sciencedata" {
				continue
			}
			max_encounter_by_nick[loot_info.Nickname] += 1

			e.SetPermitted(permitted_wrecks, permitted_encounters, loot_info)

			found_npc_loot_in_encounters[loot_info.Nickname] = true
			e.findable_in_loot_cache[loot_info.Nickname] = true

			loots = append(loots, loot_info)

		}

		// Fallback, add to not found then
		if _, ok := found_npc_loot_in_encounters[npc_loot.Nickname]; !ok {
			not_matched_loot = append(not_matched_loot, npc_loot)
		}

	})

	for _, npc_loot := range not_matched_loot {
		// Fallback, add to not found then
		if _, ok := found_npc_loot_in_encounters[npc_loot.Nickname]; !ok {
			is_lootable, is_plugin_lootable := e.IsLootable(npc_loot.Nickname, npc_loot.LootSource)
			if !is_lootable {
				continue
			}
			_ = is_plugin_lootable

			e.SetPermitted(permitted_wrecks, permitted_encounters, npc_loot.LootInfo)

			if max_notfound_npc_drop[npc_loot.LootInfo.Nickname] > LootMaxNotFoundNPCDrops {
				continue
			}
			max_notfound_npc_drop[npc_loot.LootInfo.Nickname] += 1

			loots = append(loots, npc_loot.LootInfo)
		}
	}

	unique_phantom := make(map[string]bool)
	for _, phantom_loot := range e.Mapped.LookProps.PhantomLoots {
		item_nickname := phantom_loot.Nickname.Get()
		if _, ok := unique_phantom[item_nickname]; ok {
			continue
		}
		unique_phantom[item_nickname] = true
		e.findable_in_loot_cache[item_nickname] = true
		loots = append(loots, &LootInfo{
			Nickname:   item_nickname,
			Kind:       LootPhantom,
			SystemName: "Sirius Sector",
			Permitted:  true,
			PlaceNick:  "phantom_loot",
		})
	}

	if e.Mapped.FLSR != nil {
		missions := e.Mapped.FLSR.FLSRMissions.Missions

		for _, mission := range missions {
			for _, msn_npc := range mission.MsnNpc {

				system_nickname, found_system := msn_npc.System.GetValue()
				if !found_system {
					continue
				}

				system, found_sys := e.Mapped.Systems.SystemsMap[system_nickname]
				if !found_sys {
					continue
				}

				npc, found_npc := mission.NpcByNick[msn_npc.Npc.Get()]
				if !found_npc {
					continue
				}

				loadout, found_loadout := npc.Loadout.GetValue()
				if !found_loadout {
					continue
				}

				for _, npc_loot := range GetNpcLoots(loadout, npc.Archetype.Get()) {
					loot_info := npc_loot.LootInfo
					flsr_loot = append(flsr_loot, loot_info)
					if loot_event_sources_by_loot_nickname[loot_info.Nickname] == nil {
						loot_event_sources_by_loot_nickname[loot_info.Nickname] = make(map[string]bool)
					}
					loot_event_sources_by_loot_nickname[loot_info.Nickname][mission.Nickname.Get()] = true

					loot_info.Pos = msn_npc.Pos.Get()
					loot_info.Kind = LootFLSRNPC

					system_uni := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(system_nickname)]
					loot_info.SectorCoord = VectorToSectorCoord(system_uni, loot_info.Pos)
					loot_info.SystemName = e.GetInfocardName(system_uni.StridName.Get(), system.Nickname)
					loot_info.SystemNickname = system.Nickname
					loot_info.ObjNickname = msn_npc.Nickname.Get()

					flsr_loot = append(flsr_loot, loot_info)
				}

			}
		}

	}

	unique_npc_loot := make(map[string]bool)
	if e.Mapped.FLSR != nil {
		// add all flsr_wrecks that were not previously found in regular wrecks
		for _, item := range flsr_loot {
			e.findable_in_loot_cache[item.Nickname] = true
			findable_limit_wrecks[item.Nickname] += 1
			if findable_limit_wrecks[item.Nickname] > LootMaxFLSR {
				continue
			}

			if item.Kind == LootFLSRNPC {
				key := item.Nickname + item.SystemName + item.SectorCoord
				if _, ok := unique_npc_loot[key]; ok {
					continue
				}
				unique_npc_loot[key] = true
			}

			e.SetPermitted(permitted_wrecks, permitted_encounters, item)

			if events, ok := loot_event_sources_by_loot_nickname[item.Nickname]; ok {
				var events_sorted []string
				for event, _ := range events {
					events_sorted = append(events_sorted, event)
				}
				sort.Strings(events_sorted)
				var buf strings.Builder
				for index, event := range events_sorted {
					buf.WriteString(event)

					if index != len(events_sorted)-1 {
						buf.WriteString(",")
					}
				}
				item.PlaceNick = buf.String()
			}

			loots = append(loots, item)

		}

	}

	e.findable_wrecks = loots
	return e.findable_in_loot_cache, e.findable_wrecks
}

/*
It fixes issue of Guns obtainable only via wrecks being invisible
*/
const (
	BaseLootableName     = "Lootable"
	BaseLootableFaction  = "Wrecks and Encounters"
	BaseLootableNickname = "base_loots"
)

type LootSource int8

const (
	LootSourceUnknown LootSource = iota
	LootSourceEquip
	LootSourceCargo
	LootSourceAny
)

func (l LootSource) ToStr() string {
	if l == LootSourceEquip {
		return "equip"
	}
	if l == LootSourceCargo {
		return "cargo"
	}
	return "unknown"
}

type LootInfo struct {
	Nickname       string
	Kind           LootKind
	Pos            cfg.Vector
	SectorCoord    string
	SystemName     string
	SystemNickname string
	ObjNickname    string
	Permitted      bool
	LootSource     LootSource
	PlaceNick      string

	DiscoEncounter bool
}

func (e *Exporter) EnhanceBasesWithLoot(bases []*Base) []*Base {

	_, wrecks := e.FindableInLoot()

	base := &Base{
		Name:               "Lootable",
		MarketGoodsPerNick: make(map[CommodityKey]*MarketGood),
		Nickname:           cfg.BaseUniNick(BaseLootableNickname),
		SystemNickname:     "readme",
		System:             "README",
		Region:             "README",
		FactionName:        BaseLootableFaction,
	}

	base.Archetypes = append(base.Archetypes, BaseLootableNickname)

	for _, loot_info := range wrecks {
		market_good := &MarketGood{
			GoodInfo:  e.GetGoodInfo(loot_info.Nickname),
			BaseSells: true,
			LootInfo:  loot_info,
		}

		if market_good.Category == "" && market_good.Name == "" {
			continue
		}

		market_good_key, _ := uuid.NewUUID()
		base.MarketGoodsPerNick[CommodityKey(market_good_key.String())] = market_good
	}

	var sb []infocarder.InfocardLine
	sb = append(sb, infocarder.NewInfocardSimpleLine(base.Name))
	sb = append(sb, infocarder.NewInfocardSimpleLine(`This is only pseudo base to show availability of lootable content`))
	sb = append(sb, infocarder.NewInfocardSimpleLine(`The content is findable in wrecks or drops from ships at missions`))
	sb = append(sb, infocarder.NewInfocardSimpleLine(``))
	sb = append(sb, infocarder.NewInfocardSimpleLine(`Go to tab "BASES" and find this base there to see FULL LIST of possible lootable content`))

	e.PutInfocard(infocarder.InfocardKey(base.Nickname), sb)

	bases = append(bases, base)
	e.LootableBase = base
	return bases
}
