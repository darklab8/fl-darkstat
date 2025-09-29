package solararch_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Solar struct {
	semantic.Model
	Nickname       *semantic.String
	DockingSpheres []*semantic.String
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
	IsDisco bool // disco allows for transports jump

	PlayersCanDockBerth      bool
	PlayersCanDockMoorMedium bool
	PlayersCanDockMoorLarge  bool
}

type DockableResult struct {
	IsDockable             bool
	IsDockableByTransports bool // important distinguishing for disco. only jump and moor_large valid khm?
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
				result.IsDockableByTransports = true
				result.IsDockable = true
			}
			if docking_sphere_name == DockingSphereRing {
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

			if options.IsDisco {
				if docking_sphere_name == DockingSphereMoorMedium {
					result.IsDockableByTransports = true
				}
			} else {
				result.IsDockableByTransports = result.IsDockable
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
			Nickname: semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
		}
		solar.Map(section)

		empathy_rate_key := cfg.Key("docking_sphere")
		for good_index, _ := range section.ParamMap[empathy_rate_key] {
			solar.DockingSpheres = append(solar.DockingSpheres,
				semantic.NewString(section, cfg.Key("docking_sphere"), semantic.WithLowercaseS(), semantic.OptsS(semantic.Index(good_index)), semantic.WithoutSpacesS()))
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
