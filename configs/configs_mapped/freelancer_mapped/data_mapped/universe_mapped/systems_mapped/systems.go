package systems_mapped

import (
	"fmt"
	"strings"
	"sync"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_types"
)

const (
	KEY_OBJECT = "[object]"
)

var (
	KEY_NICKNAME = cfg.Key("nickname")
	KEY_BASE     = cfg.Key("base")
)

type MissionVignetteZone struct {
	// [zone]
	// nickname = Zone_BR07_destroy_vignette_02
	// pos = -39714, 0, -20328
	// shape = SPHERE
	// size = 10000
	// mission_type = lawful, unlawful
	// sort = 99.500000
	// vignette_type = field

	// vignettes
	semantic.Model
	Nickname     *semantic.String
	Size         *semantic.Int
	Shape        *semantic.String
	Pos          *semantic.Vect
	VignetteType *semantic.String
	MissionType  *semantic.String

	// it has mission_type = lawful, unlawful this.
	// who is lawful and unlawful? :)

	// if has vignette_type = field
	// Then it is Vignette
}

type Patrol struct {
	semantic.Model
	FactionNickname *semantic.String
	Chance          *semantic.Float
}

type MissionPatrolZone struct {
	semantic.Model
	Nickname *semantic.String
	Size     *semantic.Vect
	Shape    *semantic.String
	Pos      *semantic.Vect

	Factions []*Patrol
	// [zone]
	// nickname = Path_outcasts1_2
	// pos = -314, 0, -1553.2
	// rotate = 90, -75.2, 180
	// shape = CYLINDER
	// size = 750, 50000
	// sort = 99
	// toughness = 14
	// density = 5
	// repop_time = 30
	// max_battle_size = 4
	// pop_type = attack_patrol
	// relief_time = 20
	// path_label = BR07_outcasts1, 2
	// usage = patrol
	// mission_eligible = True
	// encounter = patrolp_assault, 14, 0.4
	// faction = fc_m_grp, 1.0
}

// TODO finish coding stuff with them
type PvEEncounter struct {
	semantic.Model
	Nickname *semantic.String
	Pos      *semantic.Vect
	Size     *semantic.Vect
	Shape    *semantic.String

	Density       *semantic.Int // Max enemies spawned
	MaxBattleSize *semantic.Int // Max enemies spawn in battle
	RepopTime     *semantic.Int // respawn speed
	ReliefTime    *semantic.Int

	Encounter []*Encounter
}

// TODO finish coding stuff with them
type Encounter struct {
	semantic.Model
	Nickname      *semantic.String
	Difficulty    *semantic.Float
	ChanceToSpawn *semantic.Float
}

type TradeLaneRing struct {
	// [Object]
	// nickname = BR07_Trade_Lane_Ring_3_1
	// ids_name = 260659
	// pos = -20293, 0, 21375
	// rotate = 0, 5, 0
	// archetype = Trade_Lane_Ring
	// next_ring = BR07_Trade_Lane_Ring_3_2
	// ids_info = 66170
	// reputation = br_n_grp
	// tradelane_space_name = 501168
	// difficulty_level = 1
	// loadout = trade_lane_ring_br_01
	// pilot = pilot_solar_easiest
	semantic.Model
	Nickname *semantic.String
	Pos      *semantic.Vect
	NextRing *semantic.String
	PrevRing *semantic.String
	// has next_ring, then it is tradelane
	// or if has Trade_Lane_Ring, then trade lane too.
}

type Base struct {
	semantic.Model
	Nickname  *semantic.String
	Base      *semantic.String // base.nickname in universe.ini
	DockWith  *semantic.String
	Archetype *semantic.String

	IDsInfo     *semantic.Int
	IdsName     *semantic.Int
	RepNickname *semantic.String
	Pos         *semantic.Vect
	System      *System
	Parent      *semantic.String
}

const (
	BaseArchetypeInvisible = "invisible_base"
)

type Jumphole struct {
	semantic.Model
	Nickname  *semantic.String
	GotoHole  *semantic.String
	Archetype *semantic.String
	Pos       *semantic.Vect
	IdsName   *semantic.Int

	System *System
}

type Object struct {
	semantic.Model
	Nickname *semantic.String
}

type Wreck struct {
	semantic.Model
	Nickname *semantic.String
	Loadout  *semantic.String
}

type Asteroids struct {
	semantic.Model
	File         *semantic.Path
	Zone         *semantic.String
	LootableZone *LootableZone
}

type LootableZone struct {
	semantic.Model
	AsteroidLootCommodity  *semantic.String
	AsteroidLootMin        *semantic.Int
	AsteroidLootMax        *semantic.Int
	DynamicLootMin         *semantic.Int
	DynamicLootMax         *semantic.Int
	AsteroidLootDifficulty *semantic.Int
	DynamicLootDifficulty  *semantic.Int
}

type Zone struct {
	semantic.Model
	Nickname *semantic.String
	Pos      *semantic.Vect
	IDsInfo  *semantic.Int
	IdsName  *semantic.Int
}

type System struct {
	semantic.ConfigModel
	Nickname           string
	Bases              []*Base
	BasesByNick        map[string]*Base
	BasesByBases       map[string]*Base
	BasesByDockWith    map[string]*Base
	AllBasesByBases    map[string][]*Base
	AllBasesByDockWith map[string][]*Base
	Jumpholes          []*Jumphole
	Tradelanes         []*TradeLaneRing
	TradelaneByNick    map[string]*TradeLaneRing

	MissionZoneVignettes []*MissionVignetteZone

	MissionsSpawnZone           []*MissionPatrolZone
	MissionsSpawnZonesByFaction map[string][]*MissionPatrolZone

	Asteroids   []*Asteroids
	ZonesByNick map[string]*Zone
	Objects     []*Object
	Wrecks      []*Wreck
}

type Config struct {
	SystemsMap map[string]*System
	Systems    []*System

	// it can contain more than one base meeting condition.
	BasesByBases    map[string]*Base
	BasesByDockWith map[string]*Base
	BasesByNick     map[string]*Base
	JumpholesByNick map[string]*Jumphole
}

type FileRead struct {
	system_key string
	file       *file.File
	ini        *inireader.INIFile
}

func Read(universe_config *universe_mapped.Config, filesystem *filefind.Filesystem) *Config {
	frelconfig := &Config{
		BasesByBases:    make(map[string]*Base),
		BasesByDockWith: make(map[string]*Base),

		BasesByNick:     make(map[string]*Base),
		JumpholesByNick: make(map[string]*Jumphole),
	}
	var wg sync.WaitGroup

	var system_files map[string]*file.File = make(map[string]*file.File)

	timeit.NewTimerF(func() {
		for _, base := range universe_config.Bases {
			base_system := base.System.Get()
			universe_system := universe_config.SystemMap[universe_mapped.SystemNickname(base_system)]
			filename := universe_system.File.FileName()
			path := filesystem.GetFile(filename)
			system_files[base.System.Get()] = file.NewFile(path.GetFilepath())
		}
	}, timeit.WithMsg("systems prepared files"))

	var system_iniconfigs map[string]*inireader.INIFile = make(map[string]*inireader.INIFile)

	func() {
		timeit.NewTimerF(func() {
			// Read system files with parallelism ^_^
			iniconfigs_channel := make(chan *FileRead)
			read_file := func(data *FileRead) {
				data.ini = inireader.Read(data.file)
				iniconfigs_channel <- data
			}
			for system_key, file := range system_files {
				go read_file(&FileRead{
					system_key: system_key,
					file:       file,
				})
			}
			for range system_files {
				result := <-iniconfigs_channel
				system_iniconfigs[result.system_key] = result.ini
			}
		}, timeit.WithMsg("Read system files with parallelism ^_^"))
	}()

	timeit.NewTimerF(func() {
		frelconfig.SystemsMap = make(map[string]*System)
		frelconfig.Systems = make([]*System, 0)
		for system_key, sysiniconf := range system_iniconfigs {
			system_to_add := &System{
				MissionsSpawnZonesByFaction: make(map[string][]*MissionPatrolZone),
				TradelaneByNick:             make(map[string]*TradeLaneRing),
				ZonesByNick:                 make(map[string]*Zone),
				AllBasesByBases:             make(map[string][]*Base),
				AllBasesByDockWith:          make(map[string][]*Base),

				BasesByNick:     make(map[string]*Base),
				BasesByBases:    make(map[string]*Base),
				BasesByDockWith: make(map[string]*Base),
			}
			system_to_add.Init(sysiniconf.Sections, sysiniconf.Comments, sysiniconf.File.GetFilepath())

			system_to_add.Nickname = system_key

			system_to_add.Bases = make([]*Base, 0)
			frelconfig.SystemsMap[system_key] = system_to_add
			frelconfig.Systems = append(frelconfig.Systems, system_to_add)

			if asteroids, ok := sysiniconf.SectionMap["[asteroids]"]; ok {
				for _, obj := range asteroids {
					asteroids_to_add := &Asteroids{
						File: semantic.NewPath(obj, cfg.Key("file")),
						Zone: semantic.NewString(obj, cfg.Key("zone"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					}
					asteroids_to_add.Map(obj)

					system_to_add.Asteroids = append(system_to_add.Asteroids, asteroids_to_add)

					wg.Add(1)
					go func(ast *Asteroids) {
						filename := ast.File.FileName()

						file_to_read := filesystem.GetFile(utils_types.FilePath(strings.ToLower(filename.ToString())))
						if file_to_read == nil {
							fmt.Println("not able to find mining file for asteroids zone", filename)
							wg.Done()
							return
						}
						config := inireader.Read(file_to_read)

						if lootable_zones, ok := config.SectionMap["[lootablezone]"]; ok {
							obj := lootable_zones[0]
							lootable_zone := &LootableZone{
								AsteroidLootCommodity: semantic.NewString(obj, cfg.Key("asteroid_loot_commodity"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),

								AsteroidLootMin:        semantic.NewInt(obj, cfg.Key("asteroid_loot_count"), semantic.Optional(), semantic.Order(0)),
								AsteroidLootMax:        semantic.NewInt(obj, cfg.Key("asteroid_loot_count"), semantic.Optional(), semantic.Order(1)),
								DynamicLootMin:         semantic.NewInt(obj, cfg.Key("dynamic_loot_count"), semantic.Optional(), semantic.Order(0)),
								DynamicLootMax:         semantic.NewInt(obj, cfg.Key("dynamic_loot_count"), semantic.Optional(), semantic.Order(1)),
								AsteroidLootDifficulty: semantic.NewInt(obj, cfg.Key("asteroid_loot_difficulty"), semantic.Optional()),
								DynamicLootDifficulty:  semantic.NewInt(obj, cfg.Key("dynamic_loot_difficulty"), semantic.Optional()),
							}
							lootable_zone.Map(obj)
							ast.LootableZone = lootable_zone
						}
						wg.Done()

					}(asteroids_to_add)
				}
			}
			if objects, ok := sysiniconf.SectionMap[KEY_OBJECT]; ok {
				for _, obj := range objects {

					object_to_add := &Object{
						Nickname: semantic.NewString(obj, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					}
					object_to_add.Map(obj)
					system_to_add.Objects = append(system_to_add.Objects, object_to_add)

					// check if it is base object
					_, has_base := obj.ParamMap[KEY_BASE]
					_, has_dock_with := obj.ParamMap[cfg.Key("dock_with")]
					_ = has_dock_with
					if has_base || has_dock_with { // || has_dock_with
						base_to_add := &Base{
							Archetype:   semantic.NewString(obj, cfg.Key("archetype"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Parent:      semantic.NewString(obj, cfg.Key("parent"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Nickname:    semantic.NewString(obj, KEY_NICKNAME, semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Base:        semantic.NewString(obj, KEY_BASE, semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							DockWith:    semantic.NewString(obj, cfg.Key("dock_with"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							RepNickname: semantic.NewString(obj, cfg.Key("reputation"), semantic.OptsS(semantic.Optional()), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							IDsInfo:     semantic.NewInt(obj, cfg.Key("ids_info"), semantic.Optional()),
							IdsName:     semantic.NewInt(obj, cfg.Key("ids_name"), semantic.Optional()),
							Pos:         semantic.NewVector(obj, cfg.Key("pos"), semantic.Precision(0)),
							System:      system_to_add,
						}
						base_to_add.Map(obj)

						system_to_add.BasesByNick[base_to_add.Nickname.Get()] = base_to_add

						if base, ok := base_to_add.Base.GetValue(); ok {
							if _, ok := system_to_add.BasesByBases[base]; !ok {
								system_to_add.BasesByBases[base] = base_to_add
							}
							system_to_add.AllBasesByBases[base] = append(system_to_add.AllBasesByBases[base], base_to_add)

							if _, ok := frelconfig.BasesByBases[base]; !ok {
								frelconfig.BasesByBases[base] = base_to_add
							}

						}

						if dock_with_base, ok := base_to_add.DockWith.GetValue(); ok {
							if _, ok := system_to_add.AllBasesByDockWith[dock_with_base]; !ok {
								system_to_add.AllBasesByDockWith[dock_with_base] = append(system_to_add.AllBasesByDockWith[dock_with_base], base_to_add)
							}
							system_to_add.AllBasesByDockWith[dock_with_base] = append(system_to_add.AllBasesByDockWith[dock_with_base], base_to_add)

							if _, ok := frelconfig.BasesByDockWith[dock_with_base]; !ok {
								frelconfig.BasesByDockWith[dock_with_base] = base_to_add
							}

						}

						system_to_add.Bases = append(system_to_add.Bases, base_to_add)

						if base_nickname, ok := base_to_add.Nickname.GetValue(); ok {
							frelconfig.BasesByNick[base_nickname] = base_to_add
						}

					}

					if _, ok := obj.ParamMap[cfg.Key("goto")]; ok {
						jumphole := &Jumphole{
							Archetype: semantic.NewString(obj, cfg.Key("archetype"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Nickname:  semantic.NewString(obj, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							GotoHole:  semantic.NewString(obj, cfg.Key("goto"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Order(1))),
							Pos:       semantic.NewVector(obj, cfg.Key("pos"), semantic.Precision(0)),
							IdsName:   semantic.NewInt(obj, cfg.Key("ids_name"), semantic.Optional()),
							System:    system_to_add,
						}

						system_to_add.Jumpholes = append(system_to_add.Jumpholes, jumphole)
						frelconfig.JumpholesByNick[jumphole.Nickname.Get()] = jumphole
					}

					if _, ok := obj.ParamMap[cfg.Key("loadout")]; ok {
						wreck := &Wreck{
							Nickname: semantic.NewString(obj, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Loadout:  semantic.NewString(obj, cfg.Key("loadout"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}

						system_to_add.Wrecks = append(system_to_add.Wrecks, wreck)
					}

					_, is_trade_lane1 := obj.ParamMap[cfg.Key("next_ring")]
					_, is_trade_lane2 := obj.ParamMap[cfg.Key("prev_ring")]
					if is_trade_lane1 || is_trade_lane2 {
						tradelane := &TradeLaneRing{
							Nickname: semantic.NewString(obj, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Pos:      semantic.NewVector(obj, cfg.Key("pos"), semantic.Precision(0)),
							NextRing: semantic.NewString(obj, cfg.Key("next_ring"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							PrevRing: semantic.NewString(obj, cfg.Key("prev_ring"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}

						system_to_add.Tradelanes = append(system_to_add.Tradelanes, tradelane)
						system_to_add.TradelaneByNick[tradelane.Nickname.Get()] = tradelane
					}

				}
			}

			if zones, ok := sysiniconf.SectionMap["[zone]"]; ok {
				for _, zone_info := range zones {

					zone_to_add := &Zone{
						Nickname: semantic.NewString(zone_info, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						Pos:      semantic.NewVector(zone_info, cfg.Key("pos"), semantic.Precision(0)),
						IdsName:  semantic.NewInt(zone_info, cfg.Key("ids_name"), semantic.Optional()),
						IDsInfo:  semantic.NewInt(zone_info, cfg.Key("ids_info"), semantic.Optional()),
					}
					system_to_add.ZonesByNick[zone_to_add.Nickname.Get()] = zone_to_add

					if vignette_type, ok := zone_info.ParamMap[cfg.Key("vignette_type")]; ok && len(vignette_type) > 0 {
						vignette := &MissionVignetteZone{
							Nickname:     semantic.NewString(zone_info, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Size:         semantic.NewInt(zone_info, cfg.Key("size"), semantic.Optional()),
							Shape:        semantic.NewString(zone_info, cfg.Key("shape"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Pos:          semantic.NewVector(zone_info, cfg.Key("pos"), semantic.Precision(2)),
							VignetteType: semantic.NewString(zone_info, cfg.Key("vignette_type"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							MissionType:  semantic.NewString(zone_info, cfg.Key("mission_type"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}
						vignette.Map(zone_info)
						system_to_add.MissionZoneVignettes = append(system_to_add.MissionZoneVignettes, vignette)
					}

					if identifier, ok := zone_info.ParamMap[cfg.Key("faction")]; ok && len(identifier) > 0 {
						spawn_area := &MissionPatrolZone{
							Nickname: semantic.NewString(zone_info, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Size:     semantic.NewVector(zone_info, cfg.Key("size"), semantic.Precision(2)),
							Shape:    semantic.NewString(zone_info, cfg.Key("shape"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Pos:      semantic.NewVector(zone_info, cfg.Key("pos"), semantic.Precision(2)),
						}
						spawn_area.Map(zone_info)

						if factions, ok := zone_info.ParamMap[cfg.Key("faction")]; ok {
							for index := range factions {
								faction := &Patrol{
									FactionNickname: semantic.NewString(zone_info, cfg.Key("faction"),
										semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(index), semantic.Order(0))),
									Chance: semantic.NewFloat(zone_info, cfg.Key("faction"), semantic.Precision(2), semantic.OptsF(semantic.Index(index), semantic.Order(1))),
								}
								faction.Map(zone_info)
								spawn_area.Factions = append(spawn_area.Factions, faction)
							}
						}

						system_to_add.MissionsSpawnZone = append(system_to_add.MissionsSpawnZone, spawn_area)

						for _, faction := range spawn_area.Factions {
							faction_nickname := faction.FactionNickname.Get()
							system_to_add.MissionsSpawnZonesByFaction[faction_nickname] = append(system_to_add.MissionsSpawnZonesByFaction[faction_nickname], spawn_area)
						}
					}
				}
			}
		}
	}, timeit.WithMsg("Map universe itself"))

	wg.Wait()

	// Making sure we selected Parent Bases to return
	for _, base := range frelconfig.BasesByBases {
		if parent_nickname, ok := base.Parent.GetValue(); ok {
			if main_base, ok := frelconfig.BasesByNick[parent_nickname]; ok {
				frelconfig.BasesByBases[main_base.Base.Get()] = main_base
			}
		}
	}
	for _, system := range frelconfig.Systems {
		for _, base := range system.BasesByBases {
			if parent_nickname, ok := base.Parent.GetValue(); ok {
				if main_base, ok := system.BasesByNick[parent_nickname]; ok {
					system.BasesByBases[main_base.Base.Get()] = main_base
				}
			}
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() []*file.File {
	var files []*file.File = make([]*file.File, 0)
	for _, system := range frelconfig.Systems {
		inifile := system.Render()
		files = append(files, inifile.Write(inifile.File))
	}
	return files
}
