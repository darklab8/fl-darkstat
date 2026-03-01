package fuse_mapped

import (
	"fmt"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type Fuse struct {
	semantic.Model
	Nickname  *semantic.String
	DeathFuse *semantic.Bool

	DoesDropCargo      bool
	LootableHardpoints map[string]bool
}

type Config struct {
	Files []*iniload.IniLoader

	Fuses   []*Fuse
	FuseMap map[string]*Fuse
}

func Read(configs []*iniload.IniLoader) *Config {
	frelconfig := &Config{
		Files:   configs,
		FuseMap: make(map[string]*Fuse),
	}

	for _, input_file := range configs {

		for i := 0; i < len(input_file.Sections); i++ {

			if input_file.Sections[i].Type == "[fuse]" {

				fuse_section := input_file.Sections[i]
				fuse := &Fuse{
					Nickname:           semantic.NewString(fuse_section, cfg.Key("name"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					DeathFuse:          semantic.NewBool(fuse_section, cfg.Key("death_fuse"), semantic.StrBool),
					LootableHardpoints: make(map[string]bool),
				}

				if fuse.Nickname.Get() == "fuse_suprise_drop_loot" {
					fmt.Println()
				}
				fuse.Map(fuse_section)

				for j := i + 1; j < len(input_file.Sections) && input_file.Sections[j].Type != "[fuse]"; j++ {
					section := input_file.Sections[j]

					switch section.Type {
					case "[destroy_hp_attachment]":
						hardpoint := semantic.NewString(section, cfg.Key("hardpoint"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
						fate := semantic.NewString(section, cfg.Key("fate"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
						if fate.Get() == "loot" {
							fuse.LootableHardpoints[hardpoint.Get()] = true
						}
					case "[dump_cargo]":
						fuse.DoesDropCargo = true
					}

				}

				frelconfig.Fuses = append(frelconfig.Fuses, fuse)
				frelconfig.FuseMap[fuse.Nickname.Get()] = fuse
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
