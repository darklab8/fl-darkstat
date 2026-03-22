package export_map

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/solararch_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type System struct {
	Nickname    string   `json:"nickname"`
	Name        string   `json:"name"`
	Pos         Coords2D `json:"galaxy_pos"`
	NavMapScale float64

	Region Region `json:"region"`

	SystemGraphInfo

	Objs      []*Obj
	Jumpholes []*Jumphole
}

func (s System) GetSquareScale() float64 {
	return 30.0 / s.NavMapScale
}

type Obj struct {
	Nickname         string
	Name             string
	Pos              cfg.Vector
	ShapeName        string
	VisibleByDefault bool
	Kind             ObjKind
}
type ObjKind int8

const (
	ObjUnknown ObjKind = iota
	ObjJumphole
	ObjTradelane
)

func (o ObjKind) ToNick() string {
	switch o {
	case ObjJumphole:
		return "jumphole"
	case ObjTradelane:
		return "tradelane"
	}
	return "unknown"
}

type Jumphole struct {
	Obj
	GotoSystem     string
	GotoSystemName string
	Kind           JumpConnectionKind
}

type Region struct {
	Name string `json:"name"`
}

type Coords2D struct {
	X *float64
	Y *float64
}

func (e *Export) ExportSystems(configs *configs_mapped.MappedConfigs) []*System {
	stats := &MissingShapes{
		solars_without_shapes: make(map[string]bool),
		shape_without_images:  make(map[string]bool),
	}
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

		if navmapscale, ok := system.NavMapScale.GetValue(); ok {
			system_to_add.NavMapScale = navmapscale
		} else {
			system_to_add.NavMapScale = 1.0
		}

		system_to_add.Name = configs.GetInfocardName(system.StridName.Get(), system.Nickname.Get())
		if strings.Contains(strings.ToLower(system_to_add.Name), "pennsyl") {
			// Fixed connection lines laying onto each other
			system_to_add.Pos.Y = ptr.Ptr(*system_to_add.Pos.Y + 0.25)
		}

		e.EnrichSystemWithObjects(configs, system_to_add, stats)
		systems = append(systems, system_to_add)
	}

	for archetype, _ := range stats.solars_without_shapes {
		logus.Log.Warn("solar without archetype", typelog.Any("archetype", archetype))
	}
	for shape, _ := range stats.shape_without_images {
		logus.Log.Warn("shape without image", typelog.Any("shape", shape))
	}

	return systems
}

type MissingShapes struct {
	solars_without_shapes map[string]bool
	shape_without_images  map[string]bool
}

func (e *Export) EnrichSystemWithObjects(
	configs *configs_mapped.MappedConfigs,
	system_to_add *System,
	stats *MissingShapes,
) {
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
		base_obj := &Obj{
			Nickname: base.Nickname.Get(),
			Pos:      base.Pos.Get(),
		}
		base_obj.Name = configs.GetInfocardName(base.IdsName.Get(), base_obj.Nickname)

		archetype := base.Archetype.Get()
		solararch := e.Mapped.Solararch.SolarsByNick[archetype]

		shape_name, found_shape := solararch.ShapeName.GetValue()
		if !found_shape {
			stats.solars_without_shapes[archetype] = true
		}

		base_obj.ShapeName = strings.ToLower(shape_name)

		if _, ok := e.Shapes.ShapesByNick[base_obj.ShapeName]; !ok && base_obj.ShapeName != "" {
			stats.shape_without_images[base_obj.ShapeName] = true
		}
		// TODO export infocards
		// if ids_info, ok := base.IDsInfo.GetValue(); ok && ids_info != 0 {
		// 	e.Exp.ExportInfocards(infocarder.InfocardKey(base_obj.Nickname), ids_info)
		// }
		base_obj.VisibleByDefault = true

		// TODO add bases
		// system_to_add.Objs = append(system_to_add.Objs, base_obj)
	}

	for _, jh_info := range system_info.Jumpholes {
		jumphole := &Jumphole{
			Obj: Obj{
				Nickname: jh_info.Nickname.Get(),
				Pos:      jh_info.Pos.Get(),
				Kind:     ObjJumphole,
			},
		}
		jumphole.Name = configs.GetInfocardName(jh_info.IdsName.Get(), jumphole.Nickname)

		archetype := jh_info.Archetype.Get()
		solararch := e.Mapped.Solararch.SolarsByNick[archetype]

		dockable := solararch.IsDockable(solararch_mapped.DockableOptions{
			IsDisco:                  e.Mapped.Discovery != nil,
			PlayersCanDockBerth:      true,
			PlayersCanDockMoorMedium: true,
			PlayersCanDockMoorLarge:  true,
		})
		if dockable.IsDockable {
			jumphole.VisibleByDefault = true
		}

		shape_name, found_shape := solararch.ShapeName.GetValue()
		e.Shapes.PermittedShapes[strings.ToLower(shape_name)] = true
		if !found_shape {
			stats.solars_without_shapes[archetype] = true
		}

		jumphole.ShapeName = strings.ToLower(shape_name)

		if target, ok := jh_info.GotoSystem.GetValue(); ok {
			jumphole.GotoSystem = target

			if value, ok := e.Mapped.Universe.SystemMap[universe_mapped.SystemNickname(target)]; ok {
				jumphole.GotoSystemName = configs.GetInfocardName(value.StridName.Get(), target)
			}
		}

		jumphole.Kind = e.GetJumpConnectionKind(jh_info)

		system_to_add.Jumpholes = append(system_to_add.Jumpholes, jumphole)
	}

	for _, obj_info := range system_info.Tradelanes {
		obj := &Obj{
			Nickname: obj_info.Nickname.Get(),
			Pos:      obj_info.Pos.Get(),
			Kind:     ObjTradelane,
		}
		obj.Name = configs.GetInfocardName(obj_info.IdsName.Get(), obj.Nickname)

		archetype := obj_info.Archetype.Get()
		solararch := e.Mapped.Solararch.SolarsByNick[archetype]

		shape_name, found_shape := solararch.ShapeName.GetValue()
		e.Shapes.PermittedShapes[strings.ToLower(shape_name)] = true
		if !found_shape {
			stats.solars_without_shapes[archetype] = true
		}
		obj.ShapeName = strings.ToLower(shape_name)

		if _, ok := e.Shapes.ShapesByNick[obj.ShapeName]; !ok && obj.ShapeName != "" {
			stats.shape_without_images[obj.ShapeName] = true
		}
		system_to_add.Objs = append(system_to_add.Objs, obj)
	}
}
