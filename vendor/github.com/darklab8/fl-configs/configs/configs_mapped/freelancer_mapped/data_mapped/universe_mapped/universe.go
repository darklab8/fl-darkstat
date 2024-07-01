/*
parse universe.ini
*/
package universe_mapped

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
)

// Feel free to map it xD
// terrain_tiny = nonmineable_asteroid90
// terrain_sml = nonmineable_asteroid60
// terrain_mdm = nonmineable_asteroid90
// terrain_lrg = nonmineable_asteroid60
// terrain_dyna_01 = mineable1_asteroid10
// terrain_dyna_02 = mineable1_asteroid10

var KEY_BASE_TERRAINS = [...]string{"terrain_tiny", "terrain_sml", "terrain_mdm", "terrain_lrg", "terrain_dyna_01", "terrain_dyna_02"}

const (
	FILENAME      = "universe.ini"
	KEY_BASE_TAG  = "[Base]"
	KEY_NICKNAME  = "nickname"
	KEY_STRIDNAME = "strid_name"
	KEY_SYSTEM    = "system"
	KEY_FILE      = "file"

	KEY_BASE_BGCS = "BGCS_base_run_by"

	KEY_SYSTEM_TAG           = "[system]"
	KEY_SYSTEM_MSG_ID_PREFIX = "msg_id_prefix"
	KEY_SYSTEM_VISIT         = "visit"
	KEY_SYSTEM_IDS_INFO      = "ids_info"
	KEY_SYSTEM_NAVMAPSCALE   = "NavMapScale"
	KEY_SYSTEM_POS           = "pos"

	KEY_TIME_TAG     = "[Time]"
	KEY_TIME_SECONDS = "seconds_per_day"
)

type Base struct {
	semantic.Model

	Nickname         *semantic.String
	System           *semantic.String
	StridName        *semantic.Int
	File             *semantic.Path
	BGCS_base_run_by *semantic.String
	// Terrains *semantic.StringStringMap
}

type BaseNickname string

type SystemNickname string

type System struct {
	semantic.Model
	Nickname *semantic.String
	// Pos        *semantic.Pos
	Msg_id_prefix *semantic.String
	Visit         *semantic.Int
	Strid_name    *semantic.Int
	Ids_info      *semantic.Int
	File          *semantic.Path
	NavMapScale   *semantic.Float
}

type Config struct {
	File     *iniload.IniLoader
	Bases    []*Base
	BasesMap map[BaseNickname]*Base

	Systems   []*System
	SystemMap map[SystemNickname]*System

	TimeSeconds *semantic.Int
}

func Read(ini *iniload.IniLoader) *Config {
	frelconfig := &Config{File: ini}

	frelconfig.TimeSeconds = semantic.NewInt(ini.SectionMap[KEY_TIME_TAG][0], KEY_TIME_TAG)
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

	if systems, ok := ini.SectionMap[KEY_SYSTEM_TAG]; ok {
		for _, system := range systems {
			system_to_add := System{
				NavMapScale: semantic.NewFloat(system, "NavMapScale", semantic.Precision(2)),
			}
			system_to_add.Map(system)

			system_to_add.Visit = semantic.NewInt(system, KEY_SYSTEM_VISIT, semantic.Optional())
			system_to_add.Strid_name = semantic.NewInt(system, KEY_STRIDNAME, semantic.Optional())
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
