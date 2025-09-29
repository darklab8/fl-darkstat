package ship_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type Ship struct {
	semantic.Model
	Nickname  *semantic.String // matches value `ship` in goods
	ShipClass *semantic.Int
	Type      *semantic.String
	IdsName   *semantic.Int
	IdsInfo   *semantic.Int
	IdsInfo1  *semantic.Int
	IdsInfo2  *semantic.Int
	IdsInfo3  *semantic.Int

	Nanobots    *semantic.Int
	Batteries   *semantic.Int
	Mass        *semantic.Float
	LinearDrag  *semantic.Float
	HoldSize    *semantic.Int
	HitPts      *semantic.Int
	NudgeForce  *semantic.Float
	StrafeForce *semantic.Float
	ShieldLink  *ShieldLink
	HpTypes     []*HpType

	SteeringTorque   *semantic.Vect
	AngularDrag      *semantic.Vect
	RotationIntertia *semantic.Vect

	/*
		Some info in Goods with category shiphull, it has link from [Ship] to hulll
		Some is good ship, it has stuff leading to [Power], [Engine], [Scanner] and [ShieldGenerator]

	*/
	ArmorMult *semantic.Float // disco only
}

type HpType struct {
	semantic.Model
	Nickname          *semantic.String
	AllowedEquipments []*semantic.String
}

type ShieldLink struct {
	semantic.Model
	ShieldClass *semantic.String
}

type Config struct {
	Files []*iniload.IniLoader

	Ships    []*Ship
	ShipsMap map[string]*Ship
}

func Read(files []*iniload.IniLoader) *Config {
	frelconfig := &Config{Files: files}
	frelconfig.Ships = make([]*Ship, 0, 100)
	frelconfig.ShipsMap = make(map[string]*Ship)

	for _, Iniconfig := range files {

		for _, section := range Iniconfig.SectionMap["[ship]"] {
			ship := &Ship{
				Nickname:    semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
				Type:        semantic.NewString(section, cfg.Key("type")),
				ShipClass:   semantic.NewInt(section, cfg.Key("ship_class")),
				IdsName:     semantic.NewInt(section, cfg.Key("ids_name"), semantic.Optional()),
				IdsInfo:     semantic.NewInt(section, cfg.Key("ids_info"), semantic.Optional()),
				IdsInfo1:    semantic.NewInt(section, cfg.Key("ids_info1")),
				IdsInfo2:    semantic.NewInt(section, cfg.Key("ids_info2")),
				IdsInfo3:    semantic.NewInt(section, cfg.Key("ids_info3")),
				Nanobots:    semantic.NewInt(section, cfg.Key("nanobot_limit")),
				Batteries:   semantic.NewInt(section, cfg.Key("shield_battery_limit")),
				Mass:        semantic.NewFloat(section, cfg.Key("mass"), semantic.Precision(2)),
				LinearDrag:  semantic.NewFloat(section, cfg.Key("linear_drag"), semantic.Precision(2)),
				HoldSize:    semantic.NewInt(section, cfg.Key("hold_size")),
				NudgeForce:  semantic.NewFloat(section, cfg.Key("nudge_force"), semantic.Precision(2)),
				StrafeForce: semantic.NewFloat(section, cfg.Key("strafe_force"), semantic.Precision(2)),
				ShieldLink: &ShieldLink{
					ShieldClass: semantic.NewString(section, cfg.Key("shield_link"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
				},
				HitPts: semantic.NewInt(section, cfg.Key("hit_pts")),

				SteeringTorque:   semantic.NewVector(section, cfg.Key("steering_torque"), semantic.Precision(2)),
				AngularDrag:      semantic.NewVector(section, cfg.Key("angular_drag"), semantic.Precision(2)),
				RotationIntertia: semantic.NewVector(section, cfg.Key("rotation_inertia"), semantic.Precision(2)),

				ArmorMult: semantic.NewFloat(section, cfg.Key("armor"), semantic.Precision(2), semantic.WithDefaultF(1.0)),
			}
			ship.Map(section)
			ship.ShieldLink.Map(section)

			for index, param := range section.ParamMap[cfg.Key("hp_type")] {
				hp_type := &HpType{
					Nickname: semantic.NewString(section, cfg.Key("hp_type"),
						semantic.OptsS(semantic.Index(index)),
						semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
				}

				for order_i, _ := range param.Values[1:] {
					hp_type.AllowedEquipments = append(hp_type.AllowedEquipments,
						semantic.NewString(section, cfg.Key("hp_type"),
							semantic.OptsS(semantic.Index(index), semantic.Order(order_i+1)),
							semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					)
				}

				hp_type.Map(section)
				ship.HpTypes = append(ship.HpTypes, hp_type)
			}

			// ids_name, ids_name_exists := ship.IdsName.GetValue()

			frelconfig.Ships = append(frelconfig.Ships, ship)
			frelconfig.ShipsMap[ship.Nickname.Get()] = ship

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
