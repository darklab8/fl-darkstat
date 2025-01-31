/*
parse universe.ini
*/
package universe_mapped

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_types"
)

// Feel free to map it xD
// terrain_tiny = nonmineable_asteroid90
// terrain_sml = nonmineable_asteroid60
// terrain_mdm = nonmineable_asteroid90
// terrain_lrg = nonmineable_asteroid60
// terrain_dyna_01 = mineable1_asteroid10
// terrain_dyna_02 = mineable1_asteroid10

var KEY_BASE_TERRAINS = [...]string{"terrain_tiny", "terrain_sml", "terrain_mdm", "terrain_lrg", "terrain_dyna_01", "terrain_dyna_02"}

var (
	KEY_NICKNAME  = cfg.Key("nickname")
	KEY_STRIDNAME = cfg.Key("strid_name")
	KEY_SYSTEM    = cfg.Key("system")
	KEY_FILE      = cfg.Key("file")

	KEY_BASE_BGCS = cfg.Key("BGCS_base_run_by")

	KEY_SYSTEM_MSG_ID_PREFIX = cfg.Key("msg_id_prefix")
	KEY_SYSTEM_VISIT         = cfg.Key("visit")
	KEY_SYSTEM_IDS_INFO      = cfg.Key("ids_info")
	KEY_SYSTEM_NAVMAPSCALE   = cfg.Key("NavMapScale")
	KEY_SYSTEM_POS           = cfg.Key("pos")

	KEY_TIME_SECONDS = cfg.Key("seconds_per_day")
)

const (
	FILENAME       = "universe.ini"
	KEY_BASE_TAG   = "[base]"
	KEY_TIME_TAG   = "[time]"
	KEY_SYSTEM_TAG = "[system]"
)

type Base struct {
	semantic.Model

	Nickname         *semantic.String
	System           *semantic.String
	StridName        *semantic.Int
	File             *semantic.Path
	BGCS_base_run_by *semantic.String
	// Terrains *semantic.StringStringMap

	ConfigBase *ConfigBase

	TraderExists bool
}

type BaseNickname string

type SystemNickname string

type System struct {
	semantic.Model
	Nickname *semantic.String
	// Pos        *semantic.Pos
	Msg_id_prefix *semantic.String
	Visit         *semantic.Int
	StridName     *semantic.Int
	Ids_info      *semantic.Int
	File          *semantic.Path
	NavMapScale   *semantic.Float

	PosX *semantic.Float
	PosY *semantic.Float
}

type Config struct {
	File     *iniload.IniLoader
	Bases    []*Base
	BasesMap map[BaseNickname]*Base

	Systems   []*System
	SystemMap map[SystemNickname]*System

	TimeSeconds *semantic.Int
}

type FileRead struct {
	base_nickname string
	file          *file.File
	ini           *inireader.INIFile
}

type RoomInfoRead struct {
	Base *Base
	Room *Room
	file *file.File
	ini  *inireader.INIFile
}

type Room struct {
	semantic.Model
	Nickname *semantic.String
	File     *semantic.Path

	RoomInfo
}

type RoomInfo struct {
	HotSpots []*HotSpot
}
type HotSpot struct {
	semantic.Model
	Name       *semantic.String
	RoomSwitch *semantic.String
}

type ConfigBase struct {
	File                  *inireader.INIFile
	Rooms                 []*Room
	RoomMapByRoomNickname map[string]*Room
}

func Read(ini *iniload.IniLoader, filesystem *filefind.Filesystem) *Config {
	frelconfig := &Config{File: ini}

	frelconfig.TimeSeconds = semantic.NewInt(ini.SectionMap[KEY_TIME_TAG][0], KEY_TIME_SECONDS)
	frelconfig.BasesMap = make(map[BaseNickname]*Base)
	frelconfig.Bases = make([]*Base, 0)
	frelconfig.SystemMap = make(map[SystemNickname]*System)
	frelconfig.Systems = make([]*System, 0)

	if bases, ok := ini.SectionMap[KEY_BASE_TAG]; ok {
		for _, base := range bases {
			base_to_add := &Base{}
			base_to_add.Map(base)
			base_to_add.Nickname = semantic.NewString(base, KEY_NICKNAME, semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			base_to_add.StridName = semantic.NewInt(base, KEY_STRIDNAME)
			base_to_add.BGCS_base_run_by = semantic.NewString(base, KEY_BASE_BGCS, semantic.OptsS(semantic.Optional()))
			base_to_add.System = semantic.NewString(base, KEY_SYSTEM, semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			base_to_add.File = semantic.NewPath(base, KEY_FILE, semantic.WithLowercaseP())

			frelconfig.Bases = append(frelconfig.Bases, base_to_add)
			frelconfig.BasesMap[BaseNickname(base_to_add.Nickname.Get())] = base_to_add
		}
	}

	// Reading Base Files
	var base_files map[string]*file.File = make(map[string]*file.File)
	timeit.NewTimerF(func() {
		for _, base := range frelconfig.Bases {
			filename := base.File.FileName()
			path := filesystem.GetFile(filename)
			base_files[base.Nickname.Get()] = file.NewFile(path.GetFilepath())
		}
	}, timeit.WithMsg("systems prepared files"))
	var base_fileconfigs map[string]*inireader.INIFile = make(map[string]*inireader.INIFile)
	func() {
		timeit.NewTimerF(func() {
			// Read system files with parallelism ^_^
			iniconfigs_channel := make(chan *FileRead)
			read_file := func(data *FileRead) {
				data.ini = inireader.Read(data.file)
				iniconfigs_channel <- data
			}
			for base_nickname, file := range base_files {
				go read_file(&FileRead{
					base_nickname: base_nickname,
					file:          file,
				})
			}
			for range base_files {
				result := <-iniconfigs_channel
				base_fileconfigs[result.base_nickname] = result.ini
			}
		}, timeit.WithMsg("Read system files with parallelism ^_^"))
	}()
	// Enhancing bases with Base File info about Rooms inside
	for _, base := range frelconfig.Bases {
		if base_file, ok := base_fileconfigs[base.Nickname.Get()]; ok {
			base.ConfigBase = &ConfigBase{
				File:                  base_file,
				RoomMapByRoomNickname: make(map[string]*Room),
			}

			if rooms, ok := base_file.SectionMap["[room]"]; ok {
				for _, room_info := range rooms {
					room := &Room{
						Nickname: semantic.NewString(room_info, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						File:     semantic.NewPath(room_info, cfg.Key("file")),
					}
					room.Map(room_info)
					base.ConfigBase.Rooms = append(base.ConfigBase.Rooms, room)
					base.ConfigBase.RoomMapByRoomNickname[room.Nickname.Get()] = room
				}
			}
		}
	}

	//Read Room Infos with parallelism
	iniconfigs_channel := make(chan *RoomInfoRead)
	readd_room_info := func(data *RoomInfoRead) {
		data.ini = inireader.Read(data.file)
		iniconfigs_channel <- data
	}
	for _, base := range frelconfig.Bases {
		for _, room := range base.ConfigBase.Rooms {
			filename := room.File.FileName()
			path := filesystem.GetFile(utils_types.FilePath(strings.ToLower(string(filename))))
			go readd_room_info(&RoomInfoRead{
				Room: room,
				file: file.NewFile(path.GetFilepath()),
				Base: base,
			})
		}

	}
	// Awaiting room info reading
	for _, base := range frelconfig.Bases {
		for _, _ = range base.ConfigBase.Rooms {
			room_info := <-iniconfigs_channel
			if sections, ok := room_info.ini.SectionMap["[hotspot]"]; ok {
				for _, hot_spot_section := range sections {
					hot_spot := &HotSpot{
						Name:       semantic.NewString(hot_spot_section, cfg.Key("name"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						RoomSwitch: semantic.NewString(hot_spot_section, cfg.Key("room_switch")),
					}
					hot_spot.Map(hot_spot_section)
					room_info.Room.HotSpots = append(room_info.Room.HotSpots, hot_spot)

					if room_switch, ok := hot_spot.RoomSwitch.GetValue(); ok {
						if strings.ToLower(room_switch) == "trader" {
							room_info.Base.TraderExists = true
							// fmt.Println("found, trader exists for base_nickname=", base.Nickname.Get())
						}
					}
				}
			}
		}
	}

	if systems, ok := ini.SectionMap[KEY_SYSTEM_TAG]; ok {
		for _, system := range systems {
			system_to_add := System{
				NavMapScale: semantic.NewFloat(system, cfg.Key("NavMapScale"), semantic.Precision(2)),
				PosX:        semantic.NewFloat(system, cfg.Key("pos"), semantic.Precision(2), semantic.OptsF(semantic.Order(0))),
				PosY:        semantic.NewFloat(system, cfg.Key("pos"), semantic.Precision(2), semantic.OptsF(semantic.Order(1))),
			}
			system_to_add.Map(system)

			system_to_add.Visit = semantic.NewInt(system, KEY_SYSTEM_VISIT, semantic.Optional())
			system_to_add.StridName = semantic.NewInt(system, KEY_STRIDNAME, semantic.Optional())
			system_to_add.Ids_info = semantic.NewInt(system, KEY_SYSTEM_IDS_INFO, semantic.Optional())
			system_to_add.Nickname = semantic.NewString(system, KEY_NICKNAME, semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			system_to_add.File = semantic.NewPath(system, KEY_FILE, semantic.WithLowercaseP())
			system_to_add.Msg_id_prefix = semantic.NewString(system, KEY_SYSTEM_MSG_ID_PREFIX, semantic.OptsS(semantic.Optional()))

			frelconfig.Systems = append(frelconfig.Systems, &system_to_add)
			frelconfig.SystemMap[SystemNickname(system_to_add.Nickname.Get())] = &system_to_add
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.File.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
