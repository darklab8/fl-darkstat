package configs_export

import (
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
		return "notfound_npc"
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
		loot_info.Permitted = permitted_encounters[loot_info.Nickname]

	} else if loot_info.Kind == LootWreck {
		loot_info.Permitted = permitted_wrecks[loot_info.Nickname]
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
	if false {

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
						send_npc_loot(loot_npc_drop)
					}
					for _, equip := range loadout.Equips {
						item_nickname := equip.Nickname.Get()
						if !e.IsLootable(item_nickname) {
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

		// Now this is purely for encounter i think? because only they search in [Ship]
		// By NPC ship loadout, we find ship class architecture
		// then we need to find ship class
		// through which we validate it is in encounters
		// throughthat we find zones in which encounter is, and we get where it is dropped
		// to validate if "equip" in mLootProps
		mlootprops_allowed := make(map[string]bool)
		for _, MLootProp := range e.Mapped.LookProps.LootProps {
			equipment_nickname := MLootProp.Nickname.Get()
			if !e.IsLootable(equipment_nickname) {
				continue
			}

			mlootprops_allowed[equipment_nickname] = true
		}
		// [ ] missing to validate dump_cargo and destroy_hp_attachment at Ship fuses
		// [ ] FIX to finding by faction=fc_c_grp, and there using encounter to validate only ship_class that will spawn
		// encounter = deep_bossencounter01, 19, 1 , patrol
		// faction = fc_c_grp, 1
		unique_encounter_loot := make(map[string]bool)
		max_encounter_by_nick := make(map[string]int)
		max_notfound_npc_drop := make(map[string]int)
		found_npc_loot_in_encounters := make(map[string]bool)
		IteratorNpcShips(func(npc_loot *NpcLoot) {
			if max_encounter_by_nick[npc_loot.LootInfo.Nickname] > LootMaxEncounters {
				return
			}
			// :x: by cargo find loadout name in [Loadout]
			// by louadout nickname, find in [NPCShipArch] ship class of it like `class_deep_cbcr` (npcships.ini)

			var ship_class_members []string
			if npcshiparch, ok := e.Mapped.NpcShips.NpcShipsByLoadout[npc_loot.LoadoutNickname]; ok {
				for _, ship_cl := range npcshiparch.ShipClass {
					ship_class_members = append(ship_class_members, ship_cl.Get())
				}
			}

			var ship_class_nicknames []string
			for _, ship_class_member := range ship_class_members {
				if ship_class, ok := e.Mapped.ShipClasses.ShipClassByMember[ship_class_member]; ok {
					ship_class_nicknames = append(ship_class_nicknames, ship_class.Nickname.Get())
				}
			}

			// by ship class nickname, find encounter formation filename in [EncounterFormation]
			// by encounter filename, find nickname of it in [EncounterParameters]
			// by encounter nickname, find in system [Object], relevant positions/system names and sector coords

			var valid_encounters []string
			for _, ship_class_nickname := range ship_class_nicknames {
				if encounters, ok := e.Mapped.Systems.EncounterByShipClass[ship_class_nickname]; ok {
					for _, encounter := range encounters {
						valid_encounters = append(valid_encounters, encounter.Nickname.Get())
					}
				}
			}

			// find zones by encounter nickname
			found_zones := make(map[string]*systems_mapped.EncounterZoneInSystem)

			for _, encounter_nick := range valid_encounters {
				if encounter_zones, ok := e.Mapped.Systems.EncounterZonesByNickname[encounter_nick]; ok {

					for _, zone := range encounter_zones {

						found_zones[zone.Zone.Nickname.Get()] = zone
					}
				}
			}

			for _, zone := range found_zones {

				// YOU FOUND IT. ADD POS/SYSTEM NAME/Sector CORD
				item_nickname := npc_loot.LootInfo.Nickname
				if !e.IsLootable(item_nickname) {
					continue
				}
				loot_info := &LootInfo{
					Nickname: item_nickname,
					Kind:     LootEncounter,
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

				if max_encounter_by_nick[loot_info.Nickname] > LootMaxEncounters {
					continue
				}
				max_encounter_by_nick[loot_info.Nickname] += 1

				found_npc_loot_in_encounters[loot_info.Nickname] = true
				SetPermitted(permitted_wrecks, permitted_encounters, loot_info)

				_, is_mloot_allowed := mlootprops_allowed[npc_loot.LootInfo.Nickname]

				is_fuse_allowed := true // TODO [ ] missing to validate dump_cargo and destroy_hp_attachment at Ship fuses

				if is_mloot_allowed || is_fuse_allowed {
					loots = append(loots, loot_info)
				}

			}

			// Fallback, add to not found then
			if _, ok := found_npc_loot_in_encounters[npc_loot.Nickname]; !ok {
				if !e.IsLootable(npc_loot.Nickname) {
					return
				}
				SetPermitted(permitted_wrecks, permitted_encounters, npc_loot.LootInfo)

				if max_notfound_npc_drop[npc_loot.LootInfo.Nickname] > LootMaxNotFoundNPCDrops {
					return
				}
				max_notfound_npc_drop[npc_loot.LootInfo.Nickname] += 1

				loots = append(loots, npc_loot.LootInfo)
			}
		})
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
	return bases
}
