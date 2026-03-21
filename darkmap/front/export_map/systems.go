package export_map

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/go-utils/utils/ptr"
)

type System struct {
	Nickname string   `json:"nickname"`
	Name     string   `json:"name"`
	Pos      Coords2D `json:"galaxy_pos"`

	Region Region `json:"region"`

	SystemGraphInfo

	BaseObjs []*BaseObj
}

type BaseObj struct {
	Nickname string
	Name     string
	Pos      cfg.Vector
}

type Region struct {
	Name string `json:"name"`
}

type Coords2D struct {
	X *float64
	Y *float64
}

func (e *Export) ExportSystems(configs *configs_mapped.MappedConfigs) []*System {
	var systems []*System
	for _, system := range configs.Universe.Systems {

		var pos_x *float64
		var pos_y *float64
		if posx, ok := system.PosX.GetValue(); ok {
			pos_x = ptr.Ptr(posx)
		}
		if posy, ok := system.PosY.GetValue(); ok {
			pos_y = ptr.Ptr(posy)
		}
		system_to_add := &System{
			Nickname: system.Nickname.Get(),
			Pos: Coords2D{
				X: pos_x,
				Y: pos_y,
			},
			Region: Region{
				Name: configs.GetRegionName(system),
			},
			SystemGraphInfo: SystemGraphInfo{
				LeadsTo: make(map[string]*JumpConnection),
			},
		}

		system_to_add.Name = configs.GetInfocardName(system.StridName.Get(), system.Nickname.Get())
		if strings.Contains(strings.ToLower(system_to_add.Name), "pennsyl") {
			// Fixed connection lines laying onto each other
			system_to_add.Pos.Y = ptr.Ptr(*system_to_add.Pos.Y + 0.25)
		}

		e.EnrichSystemWithObjects(configs, system_to_add)
		systems = append(systems, system_to_add)
	}

	return systems
}

func (e *Export) EnrichSystemWithObjects(configs *configs_mapped.MappedConfigs, system_to_add *System) {
	system_info := configs.Systems.SystemsMap[system_to_add.Nickname]

	if system_info == nil {
		return
	}

	all_bases := make(map[string]*systems_mapped.Base)
	for _, bases := range system_info.AllBasesByBases {
		for _, base := range bases {
			all_bases[base.Nickname.Get()] = base
		}
	}
	for _, bases := range system_info.AllBasesByDockWith {
		for _, base := range bases {
			all_bases[base.Nickname.Get()] = base
		}
	}

	for _, base := range all_bases {
		base_obj := &BaseObj{
			Nickname: base.Nickname.Get(),
			Pos:      base.Pos.Get(),
		}
		base_obj.Name = configs.GetInfocardName(base.IdsName.Get(), base_obj.Nickname)

		// TODO export infocards
		// if ids_info, ok := base.IDsInfo.GetValue(); ok && ids_info != 0 {
		// 	e.Exp.ExportInfocards(infocarder.InfocardKey(base_obj.Nickname), ids_info)
		// }
		system_to_add.BaseObjs = append(system_to_add.BaseObjs, base_obj)
	}
}
