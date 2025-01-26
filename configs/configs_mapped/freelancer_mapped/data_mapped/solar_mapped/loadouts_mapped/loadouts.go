package loadouts_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Cargo struct {
	semantic.Model
	Nickname *semantic.String
}

type Loadout struct {
	semantic.Model
	Nickname *semantic.String
	Cargos   []*Cargo
}

type Config struct {
	Files          []*iniload.IniLoader
	Loadouts       []*Loadout
	LoadoutsByNick map[string]*Loadout
}

const (
	FILENAME utils_types.FilePath = "loadouts.ini"
)

func Read(files []*iniload.IniLoader) *Config {
	frelconfig := &Config{
		Files:          files,
		LoadoutsByNick: make(map[string]*Loadout),
	}
	for _, input_file := range files {

		for _, section := range input_file.SectionMap["[loadout]"] {

			loadout := &Loadout{
				Nickname: semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
			}
			loadout.Map(section)

			cargo_key := cfg.Key("cargo")
			for good_index, _ := range section.ParamMap[cargo_key] {
				cargo := &Cargo{
					Nickname: semantic.NewString(section, cargo_key, semantic.WithLowercaseS(), semantic.OptsS(semantic.Index(good_index)), semantic.WithoutSpacesS()),
				}
				cargo.Map(section)
				loadout.Cargos = append(loadout.Cargos, cargo)
			}

			frelconfig.Loadouts = append(frelconfig.Loadouts, loadout)
			frelconfig.LoadoutsByNick[loadout.Nickname.Get()] = loadout

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
