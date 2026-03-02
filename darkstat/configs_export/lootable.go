package configs_export

import (
	"fmt"

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
	LootEncounter
	LootNotFoundNpcArch
	// LootPhantomLoop // lets not seek it
)

func (l LootKind) ToStr() string {
	if l == LootWreck {
		return "wreck"
	}
	if l == LootEncounter {
		return "encounter"
	}

	if l == LootNotFoundNpcArch {
		return "no_match_npc"
	}

	return "unknown"
}

const (
	LootMaxWrecks     = 4
	LootMaxEncounters = 2

	LootMaxNotFoundNPCDrops = 1
	LootMaxPhantom          = 1
)

func SetPermitted(permitted_wrecks map[string]bool, permitted_encounters map[string]bool, loot_info *LootInfo) {
	if loot_info.Kind == LootNotFoundNpcArch {
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

func (e *Exporter) IsLootable(item_nickname string) bool {
	item, found_item := e.Mapped.Equip().ItemsMap[item_nickname]
	if !found_item {
		return false
	}

	is_lootable, _ := item.Lootable.GetValue()
	return is_lootable
}

func (e *Exporter) findable_in_loot() (map[string]bool, []*LootInfo) {

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

	e.findable_in_loot_cache = make(map[string]bool)
	{
		// Purely wrecks stuff
		// get cargo and equip out of "wrecks" ships
		// validate that wrecks have archetype at solar permitting  dump_cargo
		// and validate that fuse destroy_hp_attachment at solar has fate = loot for the equip hptype
		// [x] missing to validate dump_cargo and destroy_hp_attachment at Solar archetype fuses
		findable_limit_wrecks := make(map[string]int)
		unique_wreck_loot := make(map[string]bool)
		for _, system := range e.Mapped.Systems.Systems {
			for _, wreck := range system.Wrecks {
				loadout_nickname := wreck.Loadout.Get()

				if loadout, ok := e.Mapped.Loadouts.LoadoutsByNick[loadout_nickname]; ok {

					type Item struct {
						nickname    string
						loot_source LootSource
					}

					// query here archetype of wreck ? by archetype of wreck, get Solar

					solar := e.Mapped.Solararch.SolarsByNick[wreck.Archetype.Get()]
					if is_destructible, _ := solar.Destructible.GetValue(); !is_destructible {
						continue
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
						if !e.IsLootable(item.nickname) {
							continue
						}
						loot_info := &LootInfo{
							Nickname:   item.nickname,
							Kind:       LootWreck,
							LootSource: LootSourceEquip,
							PlaceNick:  wreck.Nickname.Get(),
						}

						system_uni := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(system.Nickname)]
						loot_info.Pos = wreck.Pos.Get()
						loot_info.SectorCoord = VectorToSectorCoord(system_uni, loot_info.Pos)
						loot_info.SystemName = e.GetInfocardName(system_uni.StridName.Get(), system.Nickname)

						key_uniqueness := loot_info.Nickname + loot_info.SystemName + loot_info.SectorCoord
						if _, ok := unique_wreck_loot[key_uniqueness]; ok {
							continue
						}
						unique_wreck_loot[key_uniqueness] = true

						// TODO [ ] missing to validate dump_cargo and destroy_hp_attachment at Solar archetype fuses
						is_fuse_allowed := true
						if !is_fuse_allowed {
							continue
						}

						e.findable_in_loot_cache[item.nickname] = true
						findable_limit_wrecks[item.nickname] += 1
						if findable_limit_wrecks[item.nickname] <= LootMaxWrecks {
							SetPermitted(permitted_wrecks, permitted_encounters, loot_info)
							loots = append(loots, loot_info)
						}
					}
				}
			}
		}
	}
	{

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
			if !e.IsLootable(equipment_nickname) {
				continue
			}

			mlootprops_allowed[equipment_nickname] = true
		}

		type NpcLoot struct {
			*LootInfo
			LoadoutNickname string
		}

		unique_loadouts := make(map[string]bool)
		// May time to make it iterator :)
		IteratorNpcShips := func(send_npc_loot func(npc_loot *NpcLoot)) {
			for _, npc_arch := range e.Mapped.NpcShips.NpcShips {
				loadout_nickname := npc_arch.Loadout.Get()

				if _, ok := unique_loadouts[loadout_nickname]; ok {
					continue
				}
				unique_loadouts[loadout_nickname] = true

				if loadout, ok := e.Mapped.Loadouts.LoadoutsByNick[loadout_nickname]; ok {

					loadout_archetype := npc_arch.ShipArchetype.Get()

					shiparch := e.Mapped.Shiparch.ShipsMap[loadout_archetype]
					forbidden_hardpoints := make(map[string]bool)
					fuse_drops_equips := make(map[string]bool)
					fuse_cargo_drop := false

					for _, fuse := range shiparch.Fuses {
						if fuse, ok := e.Mapped.Fuses.FuseMap[fuse.Get()]; ok {
							for key, _ := range fuse.NotLootableHardpoints {
								forbidden_hardpoints[key] = true
							}
							if fuse.DoesDropCargo {
								fuse_cargo_drop = true
							}

							for key, _ := range fuse.LootableHardpoints {
								fuse_drops_equips[key] = true
							}
						}
					}

					for _, cargo := range loadout.Cargos {
						item_nickname := cargo.Nickname.Get()
						if !e.IsLootable(item_nickname) {
							continue
						}
						loot_info := &LootInfo{
							Kind:       LootNotFoundNpcArch,
							Nickname:   item_nickname,
							LootSource: LootSourceCargo,
						}
						loot_npc_drop := &NpcLoot{
							LootInfo:        loot_info,
							LoadoutNickname: loadout_nickname,
						}

						loot_droppable := false
						if _, ok := mlootprops_allowed[loot_info.Nickname]; ok {
							loot_droppable = true
						}

						if !loot_droppable && !fuse_cargo_drop {
							continue
						}

						send_npc_loot(loot_npc_drop)
					}
					for _, equip := range loadout.Equips {
						item_nickname := equip.Nickname.Get()
						if !e.IsLootable(item_nickname) {
							continue
						}

						// skipping dissapearing hardpoint equips
						hardpoint, found_hardpoint := equip.Hardpoint.GetValue()
						if !found_hardpoint {
							continue
						}
						_, is_forbidden_hardpoint := forbidden_hardpoints[hardpoint]
						if is_forbidden_hardpoint {
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

						if !loot_droppable && !fuse_droppable {
							continue
						}

						loot_info := &LootInfo{
							Kind:       LootNotFoundNpcArch,
							Nickname:   item_nickname,
							LootSource: LootSourceEquip,
						}
						loot_npc_drop := &NpcLoot{
							LootInfo:        loot_info,
							LoadoutNickname: loadout_nickname,
						}

						send_npc_loot(loot_npc_drop)
					}
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

			if npc_loot.LootInfo.Nickname == "no2_gun_medium03" {
				fmt.Println()
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
				if !e.IsLootable(item_nickname) {
					continue
				}
				loot_info := &LootInfo{
					Nickname:   item_nickname,
					Kind:       LootEncounter,
					LootSource: npc_loot.LootSource,
					PlaceNick:  zone.Zone.Nickname.Get(),
				}

				if _, ok := zone.Zone.Pos.GetValue(); ok {
					loot_info.Pos = zone.Zone.Pos.Get()
				}

				system_uni := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(zone.System.Nickname)]
				loot_info.SectorCoord = VectorToSectorCoord(system_uni, loot_info.Pos)
				loot_info.SystemName = e.GetInfocardName(system_uni.StridName.Get(), zone.System.Nickname)
				loot_info.Kind = LootEncounter

				key_uniqueness := loot_info.Nickname + zone.System.Nickname + loot_info.SectorCoord
				if _, ok := unique_encounter_loot[key_uniqueness]; ok {
					continue
				}
				unique_encounter_loot[key_uniqueness] = true

				if max_encounter_by_nick[loot_info.Nickname] > LootMaxEncounters && loot_info.Nickname != "commodity_sciencedata" {
					continue
				}
				max_encounter_by_nick[loot_info.Nickname] += 1

				SetPermitted(permitted_wrecks, permitted_encounters, loot_info)

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
				if !e.IsLootable(npc_loot.Nickname) {
					continue
				}
				SetPermitted(permitted_wrecks, permitted_encounters, npc_loot.LootInfo)

				if max_notfound_npc_drop[npc_loot.LootInfo.Nickname] > LootMaxNotFoundNPCDrops {
					continue
				}
				max_notfound_npc_drop[npc_loot.LootInfo.Nickname] += 1

				loots = append(loots, npc_loot.LootInfo)
			}
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
	BaseLootableFaction  = "Wrecks and Missions"
	BaseLootableNickname = "base_loots"
)

type LootSource int8

const (
	LootSourceUnknown LootSource = iota
	LootSourceEquip
	LootSourceCargo
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
	Nickname    string
	Kind        LootKind
	Pos         cfg.Vector
	SectorCoord string
	SystemName  string
	Permitted   bool
	LootSource  LootSource
	PlaceNick   string
}

func (e *Exporter) EnhanceBasesWithLoot(bases []*Base) []*Base {

	_, wrecks := e.findable_in_loot()

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
