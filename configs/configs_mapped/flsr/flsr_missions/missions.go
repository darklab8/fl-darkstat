package flsr_missions

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type Mission struct {
	semantic.Model
	Nickname  *semantic.String
	InitState *semantic.Bool

	Solars    []*MsnSolar
	MsnNpc    []*MsnNpc
	NpcByNick map[string]*Npc
}

type MsnSolar struct {
	semantic.Model
	Nickname  *semantic.String
	Archetype *semantic.String
	System    *semantic.String
	Pos       *semantic.Vect
	Loadout   *semantic.String
}

type Npc struct {
	semantic.Model
	Nickname  *semantic.String
	Archetype *semantic.String
	Loadout   *semantic.String
}
type MsnNpc struct {
	semantic.Model
	Nickname *semantic.String
	Npc      *semantic.String
	System   *semantic.String
	Pos      *semantic.Vect
}

type Config struct {
	Files []*iniload.IniLoader

	Missions []*Mission
}

func Read(configs []*iniload.IniLoader) *Config {
	frelconfig := &Config{
		Files: configs,
	}

	for _, input_file := range configs {

		for i := 0; i < len(input_file.Sections); i++ {

			if input_file.Sections[i].Type == "[mission]" {

				fuse_section := input_file.Sections[i]
				msn := &Mission{
					Nickname:  semantic.NewString(fuse_section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					InitState: semantic.NewBool(fuse_section, cfg.Key("initstate"), semantic.FLSRActiveBool),
					NpcByNick: make(map[string]*Npc),
				}
				msn.Map(fuse_section)

				for j := i + 1; j < len(input_file.Sections) && input_file.Sections[j].Type != "[mission]"; j++ {
					section := input_file.Sections[j]

					switch section.Type {
					case "[msnsolar]":
						solar := &MsnSolar{
							Nickname:  semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Archetype: semantic.NewString(section, cfg.Key("archetype"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							System:    semantic.NewString(section, cfg.Key("system"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Pos:       semantic.NewVector(section, cfg.Key("position"), semantic.Precision(0)),
							Loadout:   semantic.NewString(section, cfg.Key("loadout"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}
						solar.Map(section)

						msn.Solars = append(msn.Solars, solar)
					case "[npc]":
						obj := &Npc{
							Nickname:  semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Archetype: semantic.NewString(section, cfg.Key("archetype"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Loadout:   semantic.NewString(section, cfg.Key("loadout"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}
						obj.Map(section)

						msn.NpcByNick[obj.Nickname.Get()] = obj
					case "[msnnpc]":
						obj := &MsnNpc{
							Nickname: semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							System:   semantic.NewString(section, cfg.Key("system"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
							Pos:      semantic.NewVector(section, cfg.Key("position"), semantic.Precision(0)),
							Npc:      semantic.NewString(section, cfg.Key("npc"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}
						obj.Map(section)

						msn.MsnNpc = append(msn.MsnNpc, obj)
					}

				}

				frelconfig.Missions = append(frelconfig.Missions, msn)
			}

		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() []*file.File {
	var files []*file.File
	for _, file := range frelconfig.Files {
		inifile := file.Render()
		inifile.Write(inifile.File)
		files = append(files, inifile.File)
	}
	return files
}
