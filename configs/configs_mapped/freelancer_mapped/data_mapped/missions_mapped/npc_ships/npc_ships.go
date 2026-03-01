package npc_ships

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

const (
	FILENAME = "npcships.ini"
)

type NPCShipArch struct {
	semantic.Model
	Nickname  *semantic.String
	Level     *semantic.String
	NpcClass  []*semantic.String
	Loadout   *semantic.String
	ShipClass []*semantic.String
}

type Config struct {
	*iniload.IniLoader

	NpcShips           []*NPCShipArch
	NpcShipsByNickname map[string]*NPCShipArch
	NpcShipsByLoadout  map[string]*NPCShipArch
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader:          input_file,
		NpcShipsByNickname: make(map[string]*NPCShipArch),
		NpcShipsByLoadout:  make(map[string]*NPCShipArch),
	}
	if sections, ok := frelconfig.SectionMap["[npcshiparch]"]; ok {
		for _, section := range sections {
			npc_ship_arch := &NPCShipArch{
				Loadout: semantic.NewString(section, cfg.Key("loadout"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
			}
			npc_ship_arch.Map(section)
			npc_ship_arch.Nickname = semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			npc_ship_arch.Level = semantic.NewString(section, cfg.Key("level"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())

			if npc_class_param, ok := section.ParamMap[cfg.Key("npc_class")]; ok {
				for param_index, _ := range npc_class_param {
					npc_ship_arch.NpcClass = append(npc_ship_arch.NpcClass,
						semantic.NewString(section, cfg.Key("npc_class"), semantic.OptsS(semantic.Order(0), semantic.Index(param_index)), semantic.WithLowercaseS(), semantic.WithoutSpacesS()))

					npc_ship_arch.ShipClass = append(npc_ship_arch.ShipClass,
						semantic.NewString(section, cfg.Key("npc_class"), semantic.OptsS(semantic.Order(1), semantic.Index(param_index)), semantic.WithLowercaseS(), semantic.WithoutSpacesS()))
				}

			}
			frelconfig.NpcShips = append(frelconfig.NpcShips, npc_ship_arch)
			frelconfig.NpcShipsByNickname[npc_ship_arch.Nickname.Get()] = npc_ship_arch

			frelconfig.NpcShipsByLoadout[npc_ship_arch.Loadout.Get()] = npc_ship_arch
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
