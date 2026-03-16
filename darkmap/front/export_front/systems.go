package export_front

import (
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/go-utils/utils/ptr"
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

type Coords2D struct {
	X *float64
	Y *float64
}

func ExportSystems(configs *configs_mapped.MappedConfigs) []System {
	var systems []System
	for _, system := range configs.Universe.Systems {

		var pos_x *float64
		var pos_y *float64
		if posx, ok := system.PosX.GetValue(); ok {
			pos_x = ptr.Ptr(posx)
		}
		if posy, ok := system.PosY.GetValue(); ok {
			pos_y = ptr.Ptr(posy)
		}
		system_to_add := System{
			Nickname: system.Nickname.Get(),
			Pos: Coords2D{
				X: pos_x,
				Y: pos_y,
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
