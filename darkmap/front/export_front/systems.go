package export_front

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
)

type System struct {
	Nickname string   `json:"nickname"`
	Name     string   `json:"name"`
	Pos      Coords2D `json:"galaxy_pos"`

	Region Region `json:"region"`
}

type Region struct {
	Name string `json:"name"`
}

func (region Region) ToHexColor() string {
	region_name := strings.ToLower(region.Name)
	if strings.Contains(region_name, "liberty") {
		return "#2299F5"
	}
	if strings.Contains(region_name, "gallia") {
		return "#5961ff"
	}
	if strings.Contains(region_name, "nomad") {
		return "#368674"
	}
	if strings.Contains(region_name, "kusari") {
		return "#FED433"
	}
	if strings.Contains(region_name, "bretonia") {
		return "#d7363c"
	}
	if strings.Contains(region_name, "gallia") {
		return "#5961ff"
	}
	if strings.Contains(region_name, "gallia") {
		return "#5961ff"
	}
	if strings.Contains(region_name, "rheinland") {
		return "#00AD1D"
	}
	if strings.Contains(region_name, "independent") {
		return "#DBDBDB"
	}
	if strings.Contains(region_name, "outcasts") {
		return "#A42AFA"
	}
	if strings.Contains(region_name, "corsairs") {
		return "#A42AFA"
	}
	if strings.Contains(region_name, "edge") {
		return "#414141" // #C35B2B
	}

	return "#959597"
}

type Coords2D struct {
	X float64
	Y float64
}

func ExportSystems(configs *configs_mapped.MappedConfigs) []System {
	var systems []System
	for _, system := range configs.Universe.Systems {
		system_to_add := System{
			Nickname: system.Nickname.Get(),
			Pos: Coords2D{
				X: system.PosX.Get(),
				Y: system.PosY.Get(),
			},
			Region: Region{
				Name: configs.GetRegionName(system),
			},
		}
		system_to_add.Name = configs.GetInfocardName(system.StridName.Get(), system.Nickname.Get())

		systems = append(systems, system_to_add)
	}

	return systems
}
