package faction_props_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

const (
	FILENAME = "faction_prop.ini"
)

type FactionProp struct {
	semantic.Model
	Affiliation *semantic.String
	NpcShips    []*semantic.String
}

type Config struct {
	*iniload.IniLoader

	FactionProps             []*FactionProp
	FactionPropMapByNickname map[string]*FactionProp
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader:                input_file,
		FactionPropMapByNickname: make(map[string]*FactionProp),
	}
	if sections, ok := frelconfig.SectionMap["[factionprops]"]; ok {
		for _, section := range sections {
			faction_prop := &FactionProp{}
			faction_prop.Map(section)
			faction_prop.Affiliation = semantic.NewString(section, cfg.Key("affiliation"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())

			KEY_NPC_SHIP := cfg.Key("npc_ship")
			for index, _ := range section.ParamMap[KEY_NPC_SHIP] {
				faction_prop.NpcShips = append(faction_prop.NpcShips,
					semantic.NewString(section, KEY_NPC_SHIP, semantic.OptsS(semantic.Index(index)), semantic.WithLowercaseS(), semantic.WithoutSpacesS()))
			}
			frelconfig.FactionProps = append(frelconfig.FactionProps, faction_prop)
			frelconfig.FactionPropMapByNickname[faction_prop.Affiliation.Get()] = faction_prop
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
