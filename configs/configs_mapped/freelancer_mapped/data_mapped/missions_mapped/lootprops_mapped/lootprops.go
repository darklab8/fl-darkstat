package lootprops_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

const (
	FILENAME = "lootprops.ini"
)

type MLootProps struct {
	semantic.Model
	Nickname *semantic.String
}

type Config struct {
	semantic.ConfigModel
	File *iniload.IniLoader

	LootProps []*MLootProps
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		File: input_file,
	}
	for i := 0; i < len(input_file.Sections); i++ {
		if input_file.Sections[i].Type == "[mlootprops]" {

			section := input_file.Sections[i]
			loot_prop := &MLootProps{}
			loot_prop.Map(section)
			loot_prop.Nickname = semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())

			frelconfig.LootProps = append(frelconfig.LootProps, loot_prop)
		}

	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.File.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
