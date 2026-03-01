package shipclasses_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

const (
	FILENAME = "shipclasses.ini"
)

type ShipClass struct {
	semantic.Model
	Nickname *semantic.String
	Member   []*semantic.String
}

type Config struct {
	*iniload.IniLoader

	ShipClassByMember map[string][]*ShipClass
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader:         input_file,
		ShipClassByMember: make(map[string][]*ShipClass),
	}
	if sections, ok := frelconfig.SectionMap["[shipclass]"]; ok {
		for _, section := range sections {
			ship_class_info := &ShipClass{
				Nickname: semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
			}
			ship_class_info.Map(section)

			if param, ok := section.ParamMap[cfg.Key("member")]; ok {
				for index_order, _ := range param {
					ship_class_info.Member = append(ship_class_info.Member,
						semantic.NewString(section, cfg.Key("member"), semantic.OptsS(semantic.Index(index_order)), semantic.WithLowercaseS(), semantic.WithoutSpacesS()))
				}

			}

			for _, member := range ship_class_info.Member {
				frelconfig.ShipClassByMember[member.Get()] = append(frelconfig.ShipClassByMember[member.Get()], ship_class_info)
			}
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
