package solararch_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils/utils_types"
)

// material_library = solar\planets\planet_rckmnt.txm
// material_library = solar\planets\detailmaps\detailmap_rock01.txm
// material_library = solar\planets\atmosphere.txm
// material_library = solar\planets\planet_rckmnt\planet_rckmnt.mat

type Solar struct {
	semantic.Model
	Nickname        *semantic.String
	DockingSpheres  []*semantic.String
	Fuses           []*semantic.String
	Destructible    *semantic.Bool
	CargoLimit      *semantic.Int
	ShapeName       *semantic.String
	MaterialLibrary []*semantic.Path
	SolarRadius     *semantic.Float
}

const (
	DockingSphereJump       = "jump"
	DockingSphereRing       = "ring"
	DockingSphereBerth      = "berth"
	DockingSphereMoorMedium = "moor_medium"
	DockingSphereMoorLarge  = "moor_large"
)

const (
	MissionPropertyAllowsBerth      = "can_use_berths"
	MissionPropertyAllowsMoorMedium = "can_use_med_moors"
	MissionPropertyAllowsMoorLarge  = "can_use_large_moors"
)

type DockableOptions struct {
	IsDisco                 bool // disco allows for transports jump
	WithDiscoFreighterPaths cfg.WithDiscoFreighterPaths

	PlayersCanDockBerth      bool
	PlayersCanDockMoorMedium bool
	PlayersCanDockMoorLarge  bool
	ShowInitialWorldBlocked  bool
}

type DockableResult struct {
	IsDockable bool
}

/*
important knowldge to calculate IsDockable right

mission_property    This parameter sets where the ship may dock.
Possible options: can_use_berths, can_use_med_moors, can_use_large_moors.
Berths are the small docks, moors are in space and are disabled in vanilla.
You can enable them with adoxas moors plugin. See dacom.ini and dacomsrv.ini for info on how to add the dll.
*/
func (solar *Solar) IsDockable(options DockableOptions) DockableResult {
	var result DockableResult
	for _, docking_sphere := range solar.DockingSpheres {
		if docking_sphere_name, dockable := docking_sphere.GetValue(); dockable {
			if docking_sphere_name == DockingSphereJump {
				result.IsDockable = true
			}
			if options.PlayersCanDockBerth && docking_sphere_name == DockingSphereRing {
				result.IsDockable = true
			}
			if options.PlayersCanDockBerth && docking_sphere_name == DockingSphereBerth {
				result.IsDockable = true
			}
			if options.PlayersCanDockMoorMedium && docking_sphere_name == DockingSphereMoorMedium {
				result.IsDockable = true
			}
			if options.PlayersCanDockMoorLarge && docking_sphere_name == DockingSphereMoorLarge {
				result.IsDockable = true
			}
		}
	}
	return result
}

type Config struct {
	*iniload.IniLoader
	Solars       []*Solar
	SolarsByNick map[string]*Solar
}

const (
	FILENAME utils_types.FilePath = "solararch.ini"
)

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader:    input_file,
		SolarsByNick: make(map[string]*Solar),
	}

	for _, section := range input_file.SectionMap["[solar]"] {

		solar := &Solar{
			Nickname:     semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
			Destructible: semantic.NewBool(section, cfg.ParamKey("destructible"), semantic.StrBool),
			CargoLimit:   semantic.NewInt(section, cfg.ParamKey("cargo_limit")),
			SolarRadius:  semantic.NewFloat(section, cfg.ParamKey("solar_radius"), semantic.Precision(2)),
			ShapeName:    semantic.NewString(section, cfg.ParamKey("shape_name"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
		}
		solar.Map(section)

		docking_sphere_key := cfg.Key("docking_sphere")
		for good_index, _ := range section.ParamMap[docking_sphere_key] {
			solar.DockingSpheres = append(solar.DockingSpheres,
				semantic.NewString(section, cfg.Key("docking_sphere"), semantic.WithLowercaseS(), semantic.OptsS(semantic.Index(good_index)), semantic.WithoutSpacesS()))
		}

		for good_index, _ := range section.ParamMap["fuse"] {
			solar.Fuses = append(solar.Fuses,
				semantic.NewString(section, cfg.Key("fuse"), semantic.WithLowercaseS(), semantic.OptsS(semantic.Index(good_index)), semantic.WithoutSpacesS()))
		}

		for index, _ := range section.ParamMap["material_library"] {
			solar.MaterialLibrary = append(solar.MaterialLibrary,
				semantic.NewPath(section, cfg.Key("material_library"), semantic.WithoutSpacesP(), semantic.WithLowercaseP(), semantic.OptsP(semantic.Index(index))))
		}

		frelconfig.Solars = append(frelconfig.Solars, solar)
		frelconfig.SolarsByNick[solar.Nickname.Get()] = solar

	}
	return frelconfig

}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
