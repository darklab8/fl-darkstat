package exe_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils"
)

var KEY_BASE_TERRAINS = [...]string{"terrain_tiny", "terrain_sml", "terrain_mdm", "terrain_lrg", "terrain_dyna_01", "terrain_dyna_02"}

const (
	FILENAME_FL_INI = "freelancer.ini"
)

func (r *Config) GetDlls() []string {
	dlls := utils.CompL(r.Dlls, func(x *semantic.String) string { return x.Get() })
	return append([]string{"resources.dll"}, dlls...)
}

type Config struct {
	*iniload.IniLoader

	Dlls     []*semantic.String
	Markets  []*semantic.Path
	Goods    []*semantic.Path
	Equips   []*semantic.Path
	Universe []*semantic.Path
	Ships    []*semantic.Path
	Loadouts []*semantic.Path
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{IniLoader: input_file}

	if resources, ok := input_file.SectionMap["[resources]"]; ok {

		for dll_index, _ := range resources[0].ParamMap[cfg.Key("dll")] {
			frelconfig.Dlls = append(frelconfig.Dlls,
				semantic.NewString(resources[0], cfg.Key("dll"), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(dll_index))),
			)
		}
	}

	if resources, ok := input_file.SectionMap["[data]"]; ok {
		for equipment_index, _ := range resources[0].ParamMap[cfg.Key("equipment")] {
			frelconfig.Equips = append(frelconfig.Equips,
				semantic.NewPath(resources[0], cfg.Key("equipment"), semantic.WithoutSpacesP(), semantic.WithLowercaseP(), semantic.OptsP(semantic.Index(equipment_index))),
			)
		}
		for index, _ := range resources[0].ParamMap[cfg.Key("loadouts")] {
			frelconfig.Loadouts = append(frelconfig.Loadouts,
				semantic.NewPath(resources[0], cfg.Key("loadouts"), semantic.WithoutSpacesP(), semantic.WithLowercaseP(), semantic.OptsP(semantic.Index(index))),
			)
		}
		for equipment_index, _ := range resources[0].ParamMap[cfg.Key("markets")] {
			frelconfig.Markets = append(frelconfig.Markets,
				semantic.NewPath(resources[0], cfg.Key("markets"), semantic.WithoutSpacesP(), semantic.WithLowercaseP(), semantic.OptsP(semantic.Index(equipment_index))),
			)
		}
		for equipment_index, _ := range resources[0].ParamMap[cfg.Key("universe")] {
			frelconfig.Universe = append(frelconfig.Universe,
				semantic.NewPath(resources[0], cfg.Key("universe"), semantic.WithoutSpacesP(), semantic.WithLowercaseP(), semantic.OptsP(semantic.Index(equipment_index))),
			)
		}
		for equipment_index, _ := range resources[0].ParamMap[cfg.Key("goods")] {
			frelconfig.Goods = append(frelconfig.Goods,
				semantic.NewPath(resources[0], cfg.Key("goods"), semantic.WithoutSpacesP(), semantic.WithLowercaseP(), semantic.OptsP(semantic.Index(equipment_index))),
			)
		}
		for equipment_index, _ := range resources[0].ParamMap[cfg.Key("ships")] {
			frelconfig.Ships = append(frelconfig.Ships,
				semantic.NewPath(resources[0], cfg.Key("ships"), semantic.WithoutSpacesP(), semantic.WithLowercaseP(), semantic.OptsP(semantic.Index(equipment_index))),
			)
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
