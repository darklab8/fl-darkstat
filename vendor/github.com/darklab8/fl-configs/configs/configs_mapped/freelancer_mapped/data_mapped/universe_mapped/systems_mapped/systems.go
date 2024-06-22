package systems_mapped

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils/timeit"
)

const (
	KEY_OBJECT   = "[Object]"
	KEY_NICKNAME = "nickname"
	KEY_BASE     = "base"
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
}

type Jumphole struct {
	semantic.Model
	Nickname  *semantic.String
	GotoHole  *semantic.String
	Archetype *semantic.String
	Pos       *semantic.Vect
}

type System struct {
	semantic.ConfigModel
	Nickname        string
	Bases           []*Base
	BasesByNick     map[string]*Base
	BasesByBases    map[string][]*Base
	Jumpholes       []*Jumphole
	Tradelanes      []*TradeLaneRing
	TradelaneByNick map[string]*TradeLaneRing

	MissionZoneVignettes []*MissionVignetteZone

	MissionsSpawnZone           []*MissionPatrolZone
	MissionsSpawnZonesByFaction map[string][]*MissionPatrolZone
}

type Config struct {
	SystemsMap map[string]*System
	Systems    []*System
}

type FileRead struct {
	system_key string
	file       *file.File
	ini        *inireader.INIFile
}

func Read(universe_config *universe_mapped.Config, filesystem *filefind.Filesystem) *Config {
	frelconfig := &Config{}

	var system_files map[string]*file.File = make(map[string]*file.File)
	timeit.NewTimerF(func(m *timeit.Timer) {
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
		timeit.NewTimerF(func(m *timeit.Timer) {
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

	timeit.NewTimerF(func(m *timeit.Timer) {
		frelconfig.SystemsMap = make(map[string]*System)
		frelconfig.Systems = make([]*System, 0)
		for system_key, sysiniconf := range system_iniconfigs {
			system_to_add := &System{
				MissionsSpawnZonesByFaction: make(map[string][]*MissionPatrolZone),
				TradelaneByNick:             make(map[string]*TradeLaneRing),
			}
			system_to_add.Init(sysiniconf.Sections, sysiniconf.Comments, sysiniconf.File.GetFilepath())

			system_to_add.Nickname = system_key
			system_to_add.BasesByNick = make(map[string]*Base)
			system_to_add.BasesByBases = make(map[string][]*Base)
			system_to_add.Bases = make([]*Base, 0)
			frelconfig.SystemsMap[system_key] = system_to_add
			frelconfig.Systems = append(frelconfig.Systems, system_to_add)

			if objects, ok := sysiniconf.SectionMap[KEY_OBJECT]; ok {
				for _, obj := range objects {

					// check if it is base object
					if _, ok := obj.ParamMap[KEY_BASE]; ok {
						base_to_add := &Base{
							Archetype: semantic.NewString(obj, "archetype", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}
						base_to_add.Map(obj)

						base_to_add.Nickname = semantic.NewString(obj, KEY_NICKNAME, semantic.WithLowercaseS(), semantic.WithoutSpacesS())
						base_to_add.Base = semantic.NewString(obj, KEY_BASE, semantic.WithLowercaseS(), semantic.WithoutSpacesS())
						base_to_add.DockWith = semantic.NewString(obj, "dock_with", semantic.OptsS(semantic.Optional()))
						base_to_add.RepNickname = semantic.NewString(obj, "reputation", semantic.OptsS(semantic.Optional()), semantic.WithLowercaseS(), semantic.WithoutSpacesS())

						base_to_add.IDsInfo = semantic.NewInt(obj, "ids_info", semantic.Optional())
						base_to_add.IdsName = semantic.NewInt(obj, "ids_name", semantic.Optional())

						base_to_add.Pos = semantic.NewVector(obj, "pos", semantic.Precision(0))

						system_to_add.BasesByNick[base_to_add.Nickname.Get()] = base_to_add
						system_to_add.BasesByBases[base_to_add.Base.Get()] = append(system_to_add.BasesByBases[base_to_add.Base.Get()], base_to_add)
						system_to_add.Bases = append(system_to_add.Bases, base_to_add)
					}

					if _, ok := obj.ParamMap["jump_effect"]; ok {
						jumphole := &Jumphole{
							Archetype: semantic.NewString(obj, "archetype", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Nickname:  semantic.NewString(obj, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							GotoHole:  semantic.NewString(obj, "goto", semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Order(1))),
							Pos:       semantic.NewVector(obj, "pos", semantic.Precision(0)),
						}

						system_to_add.Jumpholes = append(system_to_add.Jumpholes, jumphole)
					}

					_, is_trade_lane1 := obj.ParamMap["next_ring"]
					_, is_trade_lane2 := obj.ParamMap["prev_ring"]
					if is_trade_lane1 || is_trade_lane2 {
						tradelane := &TradeLaneRing{
							Nickname: semantic.NewString(obj, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Pos:      semantic.NewVector(obj, "pos", semantic.Precision(0)),
							NextRing: semantic.NewString(obj, "next_ring", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							PrevRing: semantic.NewString(obj, "prev_ring", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}

						system_to_add.Tradelanes = append(system_to_add.Tradelanes, tradelane)
						system_to_add.TradelaneByNick[tradelane.Nickname.Get()] = tradelane
					}

				}
			}

			if zones, ok := sysiniconf.SectionMap["[zone]"]; ok {
				for _, zone_info := range zones {

					if vignette_type, ok := zone_info.ParamMap["vignette_type"]; ok && len(vignette_type) > 0 {
						vignette := &MissionVignetteZone{
							Nickname:     semantic.NewString(zone_info, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Size:         semantic.NewInt(zone_info, "size", semantic.Optional()),
							Shape:        semantic.NewString(zone_info, "shape", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Pos:          semantic.NewVector(zone_info, "pos", semantic.Precision(2)),
							VignetteType: semantic.NewString(zone_info, "vignette_type", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							MissionType:  semantic.NewString(zone_info, "mission_type", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}
						vignette.Map(zone_info)
						system_to_add.MissionZoneVignettes = append(system_to_add.MissionZoneVignettes, vignette)
					}

					if identifier, ok := zone_info.ParamMap["faction"]; ok && len(identifier) > 0 {
						spawn_area := &MissionPatrolZone{
							Nickname: semantic.NewString(zone_info, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Size:     semantic.NewVector(zone_info, "size", semantic.Precision(2)),
							Shape:    semantic.NewString(zone_info, "shape", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Pos:      semantic.NewVector(zone_info, "pos", semantic.Precision(2)),
						}
						spawn_area.Map(zone_info)

						if factions, ok := zone_info.ParamMap["faction"]; ok {
							for index := range factions {
								faction := &Patrol{
									FactionNickname: semantic.NewString(zone_info, "faction",
										semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(index), semantic.Order(0))),
									Chance: semantic.NewFloat(zone_info, "faction", semantic.Precision(2), semantic.Index(index), semantic.Order(1)),
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
