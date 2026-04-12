package export_map

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/loadouts_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/solararch_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-darkstat/darkmap/search_bar"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
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
	Zones     []*Zone
	PoBs      []*Obj
}

func (s System) GetSquareScale() float64 {
	return 35.0 / s.NavMapScale
}

type Obj struct {
	Nickname         string
	Name             string
	Pos              cfg.Vector
	ShapeName        string
	VisibleByDefault bool
	Kind             ObjKind
	UseFallback      bool

	Star              Star
	PlanetSolarRadius float64

	Rotation cfg.Vector

	IsPoBWithKnownDockingPerms bool
}

type Zone struct {
	Obj
	PropertyFlags           int
	ZoneShape               string
	Size                    cfg.Vector
	PropertyFogColor        cfg.Vector
	PropertyFogColorEnabled bool
}

type Star struct {
	AtmosphereRange int
	StarRadius      float64
	StarGlow        Glow
	StarCenter      Glow
}

type Glow struct {
	Scale      float64
	InnerColor cfg.Vector
	OuterColor cfg.Vector
}

type ObjKind int8

const (
	ObjUnknown ObjKind = iota
	ObjJumphole
	ObjTradelane
	ObjStar
	ObjBase
	ObjPlayerBase
	ObjMining
	ObjPlanet
	ObjWreck
	ObjZone
	ObjOthers
	ObjDiscoEncounter
)

func (o ObjKind) ToNick() string {
	switch o {
	case ObjJumphole:
		return "jumphole"
	case ObjTradelane:
		return "tradelane"
	case ObjStar:
		return "star"
	case ObjPlayerBase:
		return "playerbase"
	case ObjMining:
		return "mining"
	case ObjBase:
		return "base"
	case ObjPlanet:
		return "planet"
	case ObjWreck:
		return "wreck"
	case ObjZone:
		return "zone"
	case ObjOthers:
		return "obj_other"
	case ObjDiscoEncounter:
		return "obj_disco_encounter"
	case ObjUnknown:
		return "unknown"
	}

	panic("you forgot to declare ObjKind ToNick() for some entity")
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

func IsPlanetByShape(shape_name string) bool {
	if shape_name == "nnm_sm_medium_rocky_moon" {
		return true
	} else if shape_name == "nnm_sm_medium_forest_moon" {
		return true
	} else if shape_name == "nnm_sm_small_ice_moon" {
		return true
	} else if shape_name == "nnm_sm_small_desert_moon" {
		return true
	}

	//  else if shape_name == "nnm_sm_rock_asteroid" {
	// 	all_bases[base.Nickname.Get()] = base
	// }

	return false
}

func (e *Export) EnrichSystemWithObjects(
	configs *configs_mapped.MappedConfigs,
	system_to_add *System,
	stats *MissingShapes,
) {
	system_info := configs.Systems.SystemsMap[system_to_add.Nickname]
	handled_objects := make(map[string]bool)

	if system_info == nil {
		return
	}

	all_bases := make(map[string]*systems_mapped.Base)
	if true {
		// technically this one is not interesting
		for _, bases := range system_info.AllBasesByBases {
			for _, base := range bases {
				all_bases[base.Nickname.Get()] = base
			}
		}
	}
	for _, bases := range system_info.AllBasesByDockWith {
		for _, base := range bases {
			all_bases[base.Nickname.Get()] = base
		}
	}

	for _, base_obj := range system_info.Objects {
		base := systems_mapped.NewBase(base_obj.RenderModel(), nil)
		archetype := base.Archetype.Get()
		solararch := e.Mapped.Solararch.SolarsByNick[archetype]
		shape_name, _ := solararch.ShapeName.GetValue()
		if _, ok := e.Shapes.ShapesByNick[shape_name]; !ok {
			if IsPlanetByShape(shape_name) {
				all_bases[base.Nickname.Get()] = base
			}
		} else { // if ok
			if shape_name == "nav_blackholehazard" { // permit fisher blackhole
				all_bases[base.Nickname.Get()] = base
			}
		}
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

		if strings.Contains(archetype, "invisible") {
			jumphole.VisibleByDefault = false
		}

		shape_name, found_shape := solararch.ShapeName.GetValue()
		if _, ok := e.Shapes.ShapesByNick[strings.ToLower(shape_name)]; ok {
			e.Shapes.PermittedShapes[strings.ToLower(shape_name)] = true
		} else {
			if e.Exp.Mapped.FLSR != nil {
				logus.Log.Error("FLSR can't find shape for jumphole. Using fallback",
					typelog.Any("shape", strings.ToLower(shape_name)),
					typelog.Any("obj_nick", strings.ToLower(jumphole.Nickname)),
				)
				shape_name = "nav_jumphole"
				if _, ok := e.Shapes.ShapesByNick[strings.ToLower(shape_name)]; ok {
					e.Shapes.PermittedShapes[strings.ToLower(shape_name)] = true
				} else {
					logus.Log.Panic("fallback for jumphole model is not found",
						typelog.Any("shape", strings.ToLower(shape_name)),
						typelog.Any("obj_nickname", strings.ToLower(jumphole.Nickname)),
					)
				}
			} else {
				logus.Log.Panic("can't find shape for jumphole",
					typelog.Any("shape", strings.ToLower(shape_name)),
					typelog.Any("obj_nick", strings.ToLower(jumphole.Nickname)),
				)
			}
		}
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
		handled_objects[jumphole.Nickname] = true
		system_to_add.Jumpholes = append(system_to_add.Jumpholes, jumphole)
	}

	for _, obj_info := range system_info.Tradelanes {
		obj := &Obj{
			Nickname: obj_info.Nickname.Get(),
			Pos:      obj_info.Pos.Get(),
			Kind:     ObjTradelane,
		}
		handled_objects[obj.Nickname] = true
		if value, ok := obj_info.Rotate.GetValue(); ok {
			obj.Rotation = value
		}

		obj.Name = configs.GetInfocardName(obj_info.IdsName.Get(), obj.Nickname)

		archetype := obj_info.Archetype.Get()
		solararch := e.Mapped.Solararch.SolarsByNick[archetype]

		fallback_shape_name, found_shape := solararch.ShapeName.GetValue()
		if _, ok := e.Shapes.ShapesByNick[strings.ToLower(fallback_shape_name)]; ok {
			e.Shapes.PermittedShapes[strings.ToLower(fallback_shape_name)] = true
		} else {
			logus.Log.Error("can't find shape for tradelane, going for fallback",
				typelog.Any("shape", strings.ToLower(fallback_shape_name)),
				typelog.Any("obj_nickname", strings.ToLower(obj.Nickname)),
			)
			fallback_shape_name = "nav_tradelanering"
			if _, ok := e.Shapes.ShapesByNick[strings.ToLower(fallback_shape_name)]; ok {
				e.Shapes.PermittedShapes[strings.ToLower(fallback_shape_name)] = true
			} else {
				logus.Log.Panic("fallback for tradelane model is not found",
					typelog.Any("shape", strings.ToLower(fallback_shape_name)),
					typelog.Any("obj_nickname", strings.ToLower(obj.Nickname)),
				)
			}
		}
		if !found_shape {
			stats.solars_without_shapes[archetype] = true
		}
		obj.ShapeName = strings.ToLower(fallback_shape_name)

		if _, ok := e.Shapes.ShapesByNick[obj.ShapeName]; !ok && obj.ShapeName != "" {
			stats.shape_without_images[obj.ShapeName] = true
		}
		obj.VisibleByDefault = true

		reputation_nickname, _ := obj_info.RepNickname.GetValue()
		var factionName string
		if group, exists := e.Mapped.InitialWorld.GroupsMap[reputation_nickname]; exists {
			factionName = e.GetInfocardName(group.IdsName.Get(), reputation_nickname)
		}

		e.ExportInfocard(obj_info.IDsInfo, obj.Nickname, obj.Name, obj.Pos, obj_info.IdsName, factionName)
		system_to_add.Objs = append(system_to_add.Objs, obj)
	}

	for _, star_info := range system_info.Stars {
		star := &Obj{
			Nickname: star_info.Nickname.Get(),
			Pos:      star_info.Pos.Get(),
			Kind:     ObjStar,
			Star: Star{
				AtmosphereRange: star_info.AtmosphereRange.Get(),
			},
		}

		stararch_nick := star_info.Star.Get()

		stararch := e.Mapped.Stararch.StarsByNick[stararch_nick]

		star.Star.StarRadius = stararch.Radius.Get()
		star_glow := e.Mapped.Stararch.GlowsByNick[stararch.StarGlow.Get()]

		star.Star.StarGlow = Glow{
			Scale:      star_glow.Scale.Get(),
			OuterColor: star_glow.OuterColor.Get(),
		}
		if inner_color, ok := star_glow.InnerColor.GetValue(); ok {
			star.Star.StarGlow.InnerColor = inner_color
		} else {
			star.Star.StarGlow.InnerColor = star.Star.StarGlow.OuterColor
		}

		if star_center_nick, ok := stararch.StarCenter.GetValue(); ok {
			star_center := e.Mapped.Stararch.GlowsByNick[star_center_nick]
			star.Star.StarCenter = Glow{
				Scale:      star_center.Scale.Get(),
				OuterColor: star_center.OuterColor.Get(),
			}
			if inner_color, ok := star_center.InnerColor.GetValue(); ok {
				star.Star.StarCenter.InnerColor = inner_color
			} else {
				star.Star.StarCenter.InnerColor = star.Star.StarGlow.OuterColor
			}
		}

		star.Name = configs.GetInfocardName(star_info.IdsName.Get(), star.Nickname)
		e.ExportInfocard(star_info.IDsInfo, star.Nickname, star.Name, star.Pos, star_info.IdsName, "")

		star.VisibleByDefault = true
		handled_objects[star.Nickname] = true
		system_to_add.Objs = append(system_to_add.Objs, star)

		if false {
			// if u will wish adding to search
			e.SearchEntries[star.Nickname] = search_bar.NewEntry(
				star.Name,
				"STAR",
				fmt.Sprintf("%s", system_to_add.Name),
				"star",
				"S",
				system_to_add.Nickname,
				star.Nickname,
			)
		}
	}

	for _, base_info := range all_bases {

		base := &Obj{
			Nickname: base_info.Nickname.Get(),
			Pos:      base_info.Pos.Get(),
			Kind:     ObjBase,
		}

		archetype := base_info.Archetype.Get()
		solararch := e.Mapped.Solararch.SolarsByNick[archetype]
		shape_name, found_shape := solararch.ShapeName.GetValue()
		if _, ok := e.Shapes.ShapesByNick[shape_name]; !ok {
			// if shape_name == "nnm_sm_depot" { // TODO deprecated? linker fallbacker takes care of it
			// 	shape_name = "nav_depot"
			// } else if shape_name == "nnm_sm_communications" {
			// 	shape_name = "nav_outpost"
			// } else if shape_name == "nnm_sm_mining" {
			// 	shape_name = "nav_outpost"
			// } else
			if shape_name == "nnm_sm_medium_rocky_moon" {
				base.Kind = ObjPlanet
			} else if shape_name == "nnm_sm_medium_forest_moon" {
				base.Kind = ObjPlanet
			} else if shape_name == "nnm_sm_small_ice_moon" {
				base.Kind = ObjPlanet
			} else if shape_name == "nnm_sm_rock_asteroid" {
				base.Kind = ObjPlanet
			} else if shape_name == "nnm_sm_small_desert_moon" {
				base.Kind = ObjPlanet
			}

			if base.Kind == ObjPlanet {
				material := solararch.MaterialLibrary[0].Get()
				shape_name = material.Base().ToString()
				if strings.Contains(shape_name, ".") {
					shape_name = strings.Split(shape_name, ".")[0]
				}
				base.PlanetSolarRadius = solararch.SolarRadius.Get()
			}

			if shape_name == "indust" { // disco hardcoded fix
				shape_name = "earthlike"
			}
		} else { // if shape found
			if shape_name == "nav_blackholehazard" { // permit fisher blackhole
				base.Kind = ObjPlanet
				base.PlanetSolarRadius = solararch.SolarRadius.Get()
			}
		}
		if strings.Contains(archetype, "docking_fixture") {
			if _, ok := e.Shapes.ShapesByNick[shape_name]; ok {
				e.Shapes.PermittedShapes[shape_name] = true
			} else {
				continue
			}
		} else if shape_name == "" {
			logus.Log.Error("can't find shape for base, letting linker fallbacker sorting it out",
				typelog.Any("shape", shape_name),
				typelog.Any("obj_nick", base.Nickname),
			)
			// shape_name = "nav_depot" // TODO deprecated? linker fallbacker takes care of it
		}

		if _, ok := e.Shapes.ShapesByNick[shape_name]; ok {
			e.Shapes.PermittedShapes[shape_name] = true
		} else {
			stats.solars_without_shapes[shape_name] = true
			logus.Log.Error("can't find shape for base, letting linker fallbacker sorting it out",
				typelog.Any("shape", shape_name),
				typelog.Any("obj_nick", base.Nickname),
			)
			// shape_name = "nav_depot" // TODO deprecated? linker fallbacker takes care of it
		}

		if !found_shape {
			stats.shape_without_images[archetype] = true
		}

		base.ShapeName = shape_name

		if _, ok := e.Shapes.ShapesByNick[base.ShapeName]; !ok && base.ShapeName != "" {
			stats.shape_without_images[base.ShapeName] = true
		}

		base.Name = configs.GetInfocardName(base_info.IdsName.Get(), base.Nickname)
		reputation_nickname, _ := base_info.RepNickname.GetValue()
		var factionName string
		if group, exists := e.Mapped.InitialWorld.GroupsMap[reputation_nickname]; exists {
			factionName = e.GetInfocardName(group.IdsName.Get(), reputation_nickname)
		}

		e.ExportInfocard(base_info.IDsInfo, base.Nickname, base.Name, base.Pos, base_info.IdsName, factionName)

		// dockable := solararch.IsDockable(solararch_mapped.DockableOptions{
		// 	IsDisco:                  e.Mapped.Discovery != nil,
		// 	PlayersCanDockBerth:      true,
		// 	PlayersCanDockMoorMedium: true,
		// 	PlayersCanDockMoorLarge:  true,
		// })
		// if dockable.IsDockable {
		// 	base.VisibleByDefault = true
		// }

		base.VisibleByDefault = true

		if archetype == "invisible_base" {
			base.VisibleByDefault = false
		}
		handled_objects[base.Nickname] = true
		system_to_add.Objs = append(system_to_add.Objs, base)

		e.SearchEntries[base.Nickname] = search_bar.NewEntry(
			base.Name,
			"BASE",
			fmt.Sprintf("%s", system_to_add.Name),
			"base",
			"B",
			system_to_add.Nickname,
			base.Nickname,
		)
	}

	for _, base_info := range e.PobsBySystemNick[system_info.Nickname] {

		base := &Obj{
			Nickname: string(base_info.Nickname),
			Kind:     ObjPlayerBase,
		}

		if base_info.BasePos != nil {
			base.Pos = *base_info.BasePos
		}

		if base_info.PoBCore.DefenseMode != nil {
			base.IsPoBWithKnownDockingPerms = true
		}

		base.ShapeName = "nav_depot"
		e.Shapes.PermittedShapes[base.ShapeName] = true
		base.Name = configs.GetInfocardName(int(flhash.HashNickname(base_info.Nickname)), base.Nickname)
		base.Name += " 🛰️"
		// base.Name += " ⬢╬⬢"

		if base_info.BasePos != nil {
			base.VisibleByDefault = true
		}

		var info infocarder.Infocard
		var infocard_addition infocarder.InfocardBuilder
		faction_name := ""
		if base_info.FactionName != nil {
			faction_name = *base_info.FactionName
		}
		infocard_addition = TechnicalInfoWrite(
			infocard_addition,
			base.Nickname,
			base.Pos,
			HiddenID,
			HiddenID,
			faction_name,
		)
		if value, ok := e.Exp.GetInfocard2(infocarder.InfocardKey(base.Nickname)); ok {
			info = value
		}
		e.Exp.PutInfocard(infocarder.InfocardKey(base.Nickname), append(info, infocard_addition.Lines...))
		handled_objects[base.Nickname] = true
		system_to_add.Objs = append(system_to_add.Objs, base)

		e.SearchEntries[base.Nickname] = search_bar.NewEntry(
			base.Name,
			"POB",
			fmt.Sprintf("%s", system_to_add.Name),
			"pob",
			"P",
			system_to_add.Nickname,
			base.Nickname,
		)
	}
	for _, base_info := range e.MiningBySystemNick[system_info.Nickname] {

		base := &Obj{
			Nickname: string(base_info.Nickname),
			Kind:     ObjMining,
			Pos:      base_info.Pos,
		}

		base.Name = base_info.Name
		base.Name += " ⛏️"

		if _, ok := e.MiningUsefulByNick[string(base.Nickname)]; ok {
			base.VisibleByDefault = true
		}
		handled_objects[base.Nickname] = true
		system_to_add.Objs = append(system_to_add.Objs, base)

		e.SearchEntries[base.Nickname] = search_bar.NewEntry(
			base.Name,
			"MINE",
			fmt.Sprintf("%s", system_to_add.Name),
			"mine",
			"M",
			system_to_add.Nickname,
			base.Nickname,
		)
	}

	for _, wreck := range system_info.Wrecks {

		obj := &Obj{
			Nickname: wreck.Nickname.Get(),
			Pos:      wreck.Pos.Get(),
			Kind:     ObjWreck,
		}

		loots, _ := e.Exp.ProcessWreck(configs_export.Wreck{ // check code logic there
			LoadoutNickname: wreck.Loadout.Get(),
			Archetype:       wreck.Archetype.Get(),
			Nickname:        wreck.Nickname.Get(),
			Pos:             wreck.Pos.Get(),
			Kind:            configs_export.LootWreck,
		}, system_info)

		loadout_nickname, _ := wreck.Loadout.GetValue()
		var Cargos []*loadouts_mapped.Cargo
		if loadout, ok := e.Mapped.Loadouts.LoadoutsByNick[loadout_nickname]; ok {
			Cargos = loadout.Cargos
		}

		if len(loots) == 0 && len(Cargos) == 0 {
			continue
		}

		archetype := wreck.Archetype.Get()
		solararch := e.Mapped.Solararch.SolarsByNick[archetype]
		shape_name, found_shape := solararch.ShapeName.GetValue()

		if !found_shape {
			shape_name = "nav_surprisex"
		}
		if _, ok := e.Shapes.ShapesByNick[shape_name]; ok {
			e.Shapes.PermittedShapes[shape_name] = true
		}
		obj.ShapeName = shape_name

		if _, ok := e.Shapes.ShapesByNick[obj.ShapeName]; !ok && obj.ShapeName != "" {
			stats.shape_without_images[obj.ShapeName] = true
		}

		obj.Name = configs.GetInfocardName(wreck.IdsName.Get(), obj.Nickname)
		e.ExportInfocard(wreck.IDsInfo, obj.Nickname, obj.Name, obj.Pos, wreck.IdsName, "")

		obj.VisibleByDefault = true
		handled_objects[obj.Nickname] = true
		system_to_add.Objs = append(system_to_add.Objs, obj)

		e.SearchEntries[obj.Nickname] = search_bar.NewEntry(
			obj.Name,
			"WRECK",
			fmt.Sprintf("%s", system_to_add.Name),
			"wre",
			"W",
			system_to_add.Nickname,
			obj.Nickname,
		)
	}
	for _, zone_info := range system_info.Zones {
		zone := &Zone{
			Obj: Obj{Nickname: zone_info.Nickname.Get(),
				Kind: ObjZone,
			},
		}

		zone.ZoneShape, _ = zone_info.ZoneShape.GetValue()

		zone.Rotation, _ = zone_info.Rotate.GetValue()

		zone.Size.X, _ = zone_info.SizeX.GetValue()
		zone.Size.Y, _ = zone_info.SizeY.GetValue()
		zone.Size.Z, _ = zone_info.SizeZ.GetValue()

		zone.PropertyFogColor, zone.PropertyFogColorEnabled = zone_info.PropertyFogColor.GetValue()

		zone.Pos, _ = zone_info.Pos.GetValue()

		property_flag, found_property_flag := zone_info.PropertyFlags.GetValue()

		if found_property_flag {
			zone.PropertyFlags = property_flag
		} else {
			continue
		}

		zone.Name = configs.GetInfocardName(zone_info.IdsName.Get(), zone.Nickname)
		if strings.Contains(strings.ToLower(zone.Name), "object unknown") {
			zone.Name = ""
		}
		e.ExportInfocard(zone_info.IDsInfo, zone.Nickname, zone.Name, zone.Pos, zone_info.IdsName, "")
		handled_objects[zone.Nickname] = true
		system_to_add.Zones = append(system_to_add.Zones, zone)
	}

	for _, obj_info := range system_info.Objects {
		// all other objects that have navmap defined

		if obj_info.Nickname.Get() == strings.ToLower("BAF_Encounter02_marker") {
			fmt.Print()
		}

		if _, ok := obj_info.Pos.GetValue(); !ok {
			continue
		}

		obj := &Obj{
			Nickname: obj_info.Nickname.Get(),
			Pos:      obj_info.Pos.Get(),
			Kind:     ObjOthers,
		}
		if _, ok := handled_objects[obj.Nickname]; ok {
			continue
		} else {
			handled_objects[obj.Nickname] = true
		}

		archetype := obj_info.Archetype.Get()
		solararch := e.Mapped.Solararch.SolarsByNick[archetype]
		shape_name, found_shape := solararch.ShapeName.GetValue()
		if !found_shape {
			continue
		}
		if _, ok := e.Shapes.ShapesByNick[shape_name]; ok {
			e.Shapes.PermittedShapes[shape_name] = true
		}
		obj.ShapeName = shape_name

		if _, ok := e.Shapes.ShapesByNick[obj.ShapeName]; !ok && obj.ShapeName != "" {
			stats.shape_without_images[obj.ShapeName] = true
		}

		obj.Name = configs.GetInfocardName(obj_info.IdsName.Get(), obj.Nickname)
		e.ExportInfocard(obj_info.IDsInfo, obj.Nickname, obj.Name, obj.Pos, obj_info.IdsName, "")

		if strings.Contains(strings.ToLower(obj.Name), "object unknown") {
			continue
		}

		if strings.Contains(strings.ToLower(obj.Name), "encounter") {
			obj.Kind = ObjDiscoEncounter
			obj.VisibleByDefault = true

			e.SearchEntries[obj.Nickname] = search_bar.NewEntry(
				obj.Name,
				"ENC",
				fmt.Sprintf("%s", system_to_add.Name),
				"enc",
				"E",
				system_to_add.Nickname,
				obj.Nickname,
			)
		} else {
			obj.VisibleByDefault = false
		}
		system_to_add.Objs = append(system_to_add.Objs, obj)

	}

	e.SearchEntries[system_to_add.Nickname] = search_bar.NewEntry(
		system_to_add.Name,
		"SYSTEM",
		fmt.Sprintf("%s", system_to_add.Name),
		"sys",
		"S",
		system_to_add.Nickname,
		system_to_add.Nickname,
	)
}

func (e *Export) ExportInfocard(
	ids_info *semantic.Int,
	nickname string,
	name string,
	Pos cfg.Vector,
	ids_name *semantic.Int,
	faction_name string,
) {
	var ids_info_num int
	if ids_info, ok := ids_info.GetValue(); ok && ids_info != 0 {
		e.Exp.ExportInfocards(infocarder.InfocardKey(nickname), e.GetFullInfocardIds(ids_info)...)
		ids_info_num = ids_info
	}

	ids_name_num, _ := ids_name.GetValue()

	var info infocarder.InfocardBuilder
	if value, ok := e.Exp.GetInfocard2(infocarder.InfocardKey(nickname)); ok {
		info.Lines = value
	}
	var base_name_as_infocard infocarder.Infocard = []infocarder.InfocardLine{{Phrases: []infocarder.InfocardPhrase{{Bold: true, Phrase: strings.ToUpper(name)}}}}

	info = TechnicalInfoWrite(info, nickname, Pos, ids_name_num, ids_info_num, faction_name)
	e.Exp.PutInfocard(infocarder.InfocardKey(nickname), append(base_name_as_infocard, info.Lines...))
}

const HiddenID = -1

func TechnicalInfoWrite(
	info infocarder.InfocardBuilder,
	nickname string,
	Pos cfg.Vector,
	ids_name_num int,
	ids_info_num int,
	faction_name string,
) infocarder.InfocardBuilder {
	info.WriteLineStrBold("Technical info")
	// It belongs to Alaska Security Forces.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("This object with internal nickname %s", nickname))
	sb.WriteString(fmt.Sprintf(" is located on the coordinates (%.0f,%.0f,%.0f).", Pos.X, Pos.Y, Pos.Z))
	if ids_info_num != HiddenID || ids_name_num != HiddenID {
		sb.WriteString(fmt.Sprintf(" It has name infocard number %d and infocard number %d.", ids_name_num, ids_info_num))
	}
	if faction_name != "" {
		sb.WriteString(fmt.Sprintf(" It belongs to %s.", faction_name))
	}
	info.WriteLineStr(sb.String())
	return info
}
