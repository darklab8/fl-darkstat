package equip_mapped

import (
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type Item struct {
	semantic.Model

	Category string
	Nickname *semantic.String
	IdsName  *semantic.Int
	IdsInfo  *semantic.Int
}

type Commodity struct {
	semantic.Model

	Nickname          *semantic.String
	IdsName           *semantic.Int
	IdsInfo           *semantic.Int
	UnitsPerContainer *semantic.Int
	PodApperance      *semantic.String
	LootAppearance    *semantic.String
	DecayPerSecond    *semantic.Int
	Volume            *semantic.Float
	HitPts            *semantic.Int
}

type Munition struct {
	semantic.Model
	Nickname           *semantic.String
	ExplosionArch      *semantic.String
	RequiredAmmo       *semantic.Bool
	HullDamage         *semantic.Int
	EnergyDamange      *semantic.Int
	HealintAmount      *semantic.Int
	WeaponType         *semantic.String
	LifeTime           *semantic.Float
	Mass               *semantic.Int
	Motor              *semantic.String
	MaxAngularVelocity *semantic.Float

	IdsName *semantic.Int
}

type Explosion struct {
	semantic.Model
	Nickname      *semantic.String
	HullDamage    *semantic.Int
	EnergyDamange *semantic.Int
	Radius        *semantic.Int
}

type Gun struct {
	semantic.Model
	Nickname            *semantic.String
	IdsName             *semantic.Int
	IdsInfo             *semantic.Int
	HitPts              *semantic.String // not able to read hit_pts = 5E+13 as any number yet
	PowerUsage          *semantic.Float
	RefireDelay         *semantic.Float
	MuzzleVelosity      *semantic.Float
	Toughness           *semantic.Float
	IsAutoTurret        *semantic.Bool
	TurnRate            *semantic.Float
	ProjectileArchetype *semantic.String
	HPGunType           *semantic.String
	Lootable            *semantic.Bool
}

type Mine struct {
	semantic.Model
	Nickname           *semantic.String
	ExplosionArch      *semantic.String
	AmmoLimit          *semantic.Int
	HitPts             *semantic.Int
	Lifetime           *semantic.Float
	IdsName            *semantic.Int
	IdsInfo            *semantic.Int
	SeekDist           *semantic.Int
	TopSpeed           *semantic.Int
	Acceleration       *semantic.Int
	OwnerSafeTime      *semantic.Int
	DetonationDistance *semantic.Int
	LinearDrag         *semantic.Float
}

type MineDropper struct {
	semantic.Model

	Nickname            *semantic.String
	IdsName             *semantic.Int
	IdsInfo             *semantic.Int
	HitPts              *semantic.Int
	ChildImpulse        *semantic.Float
	PowerUsage          *semantic.Float
	RefireDelay         *semantic.Float
	MuzzleVelocity      *semantic.Float
	Toughness           *semantic.Float
	ProjectileArchetype *semantic.String
	Lootable            *semantic.Bool
}

type ShieldGenerator struct {
	semantic.Model

	Nickname           *semantic.String
	IdsName            *semantic.Int
	IdsInfo            *semantic.Int
	HitPts             *semantic.Int
	Volume             *semantic.Int
	RegenerationRate   *semantic.Int
	MaxCapacity        *semantic.Int
	Toughness          *semantic.Float
	HpType             *semantic.String
	ConstPowerDraw     *semantic.Int
	RebuildPowerDraw   *semantic.Int
	OfflineRebuildTime *semantic.Int
	Lootable           *semantic.Bool
	ShieldType         *semantic.String
}

type Thruster struct {
	semantic.Model

	Nickname *semantic.String
	IdsName  *semantic.Int
	IdsInfo  *semantic.Int
	HitPts   *semantic.Int
	Lootable *semantic.Bool

	MaxForce   *semantic.Int
	PowerUsage *semantic.Int
}

type Engine struct {
	semantic.Model
	Nickname        *semantic.String
	IdsName         *semantic.Int
	IdsInfo         *semantic.Int
	CruiseSpeed     *semantic.Int
	LinearDrag      *semantic.Int
	MaxForce        *semantic.Int
	ReverseFraction *semantic.Float

	HpType           *semantic.String
	FlameEffect      *semantic.String
	TrailEffect      *semantic.String
	CruiseChargeTime *semantic.Int
}

type Power struct {
	semantic.Model
	Nickname       *semantic.String
	IdsName        *semantic.Int
	IdsInfo        *semantic.Int
	Capacity       *semantic.Int
	ChargeRate     *semantic.Int
	ThrustCapacity *semantic.Int
	ThrustRecharge *semantic.Int
}

type Tractor struct {
	semantic.Model
	Nickname   *semantic.String
	IdsName    *semantic.Int
	IdsInfo    *semantic.Int
	MaxLength  *semantic.Int
	ReachSpeed *semantic.Int
	Lootable   *semantic.Bool
}

type CounterMeasureDropper struct {
	semantic.Model
	Nickname *semantic.String
	IdsName  *semantic.Int
	IdsInfo  *semantic.Int
	Lootable *semantic.Bool

	ProjectileArchetype *semantic.String
	HitPts              *semantic.Int
	AIRange             *semantic.Int
}

type CounterMeasure struct {
	semantic.Model
	Nickname      *semantic.String
	IdsName       *semantic.Int
	IdsInfo       *semantic.Int
	AmmoLimit     *semantic.Int
	Lifetime      *semantic.Int
	Range         *semantic.Int
	DiversionPctg *semantic.Int
}

type Config struct {
	Files []*iniload.IniLoader

	Commodities    []*Commodity
	CommoditiesMap map[string]*Commodity

	Guns        []*Gun
	GunMap      map[string]*Gun
	Munitions   []*Munition
	MunitionMap map[string]*Munition

	Explosions   []*Explosion
	ExplosionMap map[string]*Explosion

	MineDroppers []*MineDropper
	Mines        []*Mine
	MinesMap     map[string]*Mine

	Items    []*Item
	ItemsMap map[string]*Item

	ShieldGens  []*ShieldGenerator
	ShidGenMap  map[string]*ShieldGenerator
	Thrusters   []*Thruster
	ThrusterMap map[string]*Thruster

	Engines    []*Engine
	EnginesMap map[string]*Engine
	Powers     []*Power
	PowersMap  map[string]*Power

	CounterMeasureDroppers []*CounterMeasureDropper
	CounterMeasure         []*CounterMeasure
	CounterMeasureMap      map[string]*CounterMeasure

	Tractors []*Tractor
}

const (
	FILENAME_SELECT_EQUIP utils_types.FilePath = "select_equip.ini"
)

func Read(files []*iniload.IniLoader) *Config {
	frelconfig := &Config{
		Files:             files,
		Guns:              make([]*Gun, 0, 100),
		Munitions:         make([]*Munition, 0, 100),
		MineDroppers:      make([]*MineDropper, 0, 100),
		MunitionMap:       make(map[string]*Munition),
		ExplosionMap:      make(map[string]*Explosion),
		MinesMap:          make(map[string]*Mine),
		EnginesMap:        make(map[string]*Engine),
		PowersMap:         make(map[string]*Power),
		CounterMeasureMap: make(map[string]*CounterMeasure),
		GunMap:            make(map[string]*Gun),
		ShidGenMap:        make(map[string]*ShieldGenerator),
		ThrusterMap:       make(map[string]*Thruster),
	}
	frelconfig.Commodities = make([]*Commodity, 0, 100)
	frelconfig.CommoditiesMap = make(map[string]*Commodity)
	frelconfig.Items = make([]*Item, 0, 100)
	frelconfig.ItemsMap = make(map[string]*Item)

	for _, file := range files {
		for _, section := range file.Sections {
			item := &Item{}
			item.Map(section)
			item.Category = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(string(section.Type), "[", ""), "]", ""))
			item.Nickname = semantic.NewString(section, "nickname", semantic.OptsS(semantic.Optional()), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			item.IdsName = semantic.NewInt(section, "ids_name", semantic.Optional())
			item.IdsInfo = semantic.NewInt(section, "ids_info", semantic.Optional())
			frelconfig.Items = append(frelconfig.Items, item)
			frelconfig.ItemsMap[item.Nickname.Get()] = item

			switch section.Type {
			case "[Commodity]":
				commodity := &Commodity{}
				commodity.Map(section)
				commodity.Nickname = semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				commodity.IdsName = semantic.NewInt(section, "ids_name")
				commodity.IdsInfo = semantic.NewInt(section, "ids_info")
				commodity.UnitsPerContainer = semantic.NewInt(section, "units_per_container")
				commodity.PodApperance = semantic.NewString(section, "pod_appearance")
				commodity.LootAppearance = semantic.NewString(section, "loot_appearance")
				commodity.DecayPerSecond = semantic.NewInt(section, "decay_per_second")
				commodity.Volume = semantic.NewFloat(section, "volume", semantic.Precision(6))
				commodity.HitPts = semantic.NewInt(section, "hit_pts")

				frelconfig.Commodities = append(frelconfig.Commodities, commodity)
				frelconfig.CommoditiesMap[commodity.Nickname.Get()] = commodity
			case "[Gun]":
				gun := &Gun{}
				gun.Map(section)

				gun.Nickname = semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				gun.IdsName = semantic.NewInt(section, "ids_name")
				gun.IdsInfo = semantic.NewInt(section, "ids_info")
				gun.HitPts = semantic.NewString(section, "hit_pts")
				gun.PowerUsage = semantic.NewFloat(section, "power_usage", semantic.Precision(2))
				gun.RefireDelay = semantic.NewFloat(section, "refire_delay", semantic.Precision(2))
				gun.MuzzleVelosity = semantic.NewFloat(section, "muzzle_velocity", semantic.Precision(2))
				gun.Toughness = semantic.NewFloat(section, "toughness", semantic.Precision(2))
				gun.IsAutoTurret = semantic.NewBool(section, "auto_turret", semantic.StrBool)
				gun.TurnRate = semantic.NewFloat(section, "turn_rate", semantic.Precision(2))
				gun.ProjectileArchetype = semantic.NewString(section, "projectile_archetype", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				gun.HPGunType = semantic.NewString(section, "hp_gun_type")
				gun.Lootable = semantic.NewBool(section, "lootable", semantic.StrBool)
				frelconfig.Guns = append(frelconfig.Guns, gun)
				frelconfig.GunMap[gun.Nickname.Get()] = gun
			case "[Munition]":
				munition := &Munition{}
				munition.Nickname = semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				munition.IdsName = semantic.NewInt(section, "ids_name")
				munition.ExplosionArch = semantic.NewString(section, "explosion_arch")
				munition.RequiredAmmo = semantic.NewBool(section, "requires_ammo", semantic.StrBool)
				munition.HullDamage = semantic.NewInt(section, "hull_damage")
				munition.EnergyDamange = semantic.NewInt(section, "energy_damage")
				munition.HealintAmount = semantic.NewInt(section, "damage")
				munition.WeaponType = semantic.NewString(section, "weapon_type", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				munition.LifeTime = semantic.NewFloat(section, "lifetime", semantic.Precision(2))
				munition.Mass = semantic.NewInt(section, "mass")
				munition.Motor = semantic.NewString(section, "motor", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				munition.MaxAngularVelocity = semantic.NewFloat(section, "max_angular_velocity", semantic.Precision(4))
				frelconfig.Munitions = append(frelconfig.Munitions, munition)
				frelconfig.MunitionMap[munition.Nickname.Get()] = munition
			case "[Explosion]":
				explosion := &Explosion{}
				explosion.Nickname = semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				explosion.HullDamage = semantic.NewInt(section, "hull_damage")
				explosion.EnergyDamange = semantic.NewInt(section, "energy_damage")
				explosion.Radius = semantic.NewInt(section, "radius")
				frelconfig.Explosions = append(frelconfig.Explosions, explosion)
				frelconfig.ExplosionMap[explosion.Nickname.Get()] = explosion
			case "[MineDropper]":
				mine_dropper := &MineDropper{
					Nickname:            semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					IdsName:             semantic.NewInt(section, "ids_name"),
					IdsInfo:             semantic.NewInt(section, "ids_info"),
					HitPts:              semantic.NewInt(section, "hit_pts"),
					ChildImpulse:        semantic.NewFloat(section, "child_impulse", semantic.Precision(2)),
					PowerUsage:          semantic.NewFloat(section, "power_usage", semantic.Precision(2)),
					RefireDelay:         semantic.NewFloat(section, "refire_delay", semantic.Precision(2)),
					MuzzleVelocity:      semantic.NewFloat(section, "muzzle_velocity", semantic.Precision(2)),
					Toughness:           semantic.NewFloat(section, "toughness", semantic.Precision(2)),
					ProjectileArchetype: semantic.NewString(section, "projectile_archetype", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					Lootable:            semantic.NewBool(section, "lootable", semantic.StrBool),
				}

				frelconfig.MineDroppers = append(frelconfig.MineDroppers, mine_dropper)
			case "[Mine]":
				mine := &Mine{
					Nickname:           semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					ExplosionArch:      semantic.NewString(section, "explosion_arch", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					AmmoLimit:          semantic.NewInt(section, "ammo_limit"),
					HitPts:             semantic.NewInt(section, "hit_pts"),
					Lifetime:           semantic.NewFloat(section, "lifetime", semantic.Precision(2)),
					IdsName:            semantic.NewInt(section, "ids_name"),
					IdsInfo:            semantic.NewInt(section, "ids_info"),
					SeekDist:           semantic.NewInt(section, "seek_dist"),
					TopSpeed:           semantic.NewInt(section, "top_speed"),
					Acceleration:       semantic.NewInt(section, "acceleration"),
					OwnerSafeTime:      semantic.NewInt(section, "owner_safe_time"),
					DetonationDistance: semantic.NewInt(section, "detonation_dist"),
					LinearDrag:         semantic.NewFloat(section, "linear_drag", semantic.Precision(6)),
				}
				frelconfig.Mines = append(frelconfig.Mines, mine)
				frelconfig.MinesMap[mine.Nickname.Get()] = mine
			case "[ShieldGenerator]":
				shield := &ShieldGenerator{
					Nickname:           semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					IdsName:            semantic.NewInt(section, "ids_name"),
					IdsInfo:            semantic.NewInt(section, "ids_info"),
					HitPts:             semantic.NewInt(section, "hit_pts"),
					Volume:             semantic.NewInt(section, "volume"),
					RegenerationRate:   semantic.NewInt(section, "regeneration_rate"),
					MaxCapacity:        semantic.NewInt(section, "max_capacity"),
					Toughness:          semantic.NewFloat(section, "toughness", semantic.Precision(2)),
					HpType:             semantic.NewString(section, "hp_type", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					ConstPowerDraw:     semantic.NewInt(section, "constant_power_draw"),
					RebuildPowerDraw:   semantic.NewInt(section, "rebuild_power_draw"),
					OfflineRebuildTime: semantic.NewInt(section, "offline_rebuild_time"),
					Lootable:           semantic.NewBool(section, "lootable", semantic.StrBool),
					ShieldType:         semantic.NewString(section, "shield_type", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
				}
				frelconfig.ShieldGens = append(frelconfig.ShieldGens, shield)
				frelconfig.ShidGenMap[shield.Nickname.Get()] = shield
			case "[Thruster]":
				thruster := &Thruster{
					Nickname:   semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					IdsName:    semantic.NewInt(section, "ids_name"),
					IdsInfo:    semantic.NewInt(section, "ids_info"),
					HitPts:     semantic.NewInt(section, "hit_pts"),
					Lootable:   semantic.NewBool(section, "lootable", semantic.StrBool),
					MaxForce:   semantic.NewInt(section, "max_force"),
					PowerUsage: semantic.NewInt(section, "power_usage"),
				}
				frelconfig.Thrusters = append(frelconfig.Thrusters, thruster)
				frelconfig.ThrusterMap[thruster.Nickname.Get()] = thruster
			case "[Power]":
				power := &Power{
					Nickname:       semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					IdsName:        semantic.NewInt(section, "ids_name"),
					IdsInfo:        semantic.NewInt(section, "ids_info"),
					Capacity:       semantic.NewInt(section, "capacity"),
					ChargeRate:     semantic.NewInt(section, "charge_rate"),
					ThrustCapacity: semantic.NewInt(section, "thrust_capacity"),
					ThrustRecharge: semantic.NewInt(section, "thrust_charge_rate"),
				}
				frelconfig.Powers = append(frelconfig.Powers, power)
				frelconfig.PowersMap[power.Nickname.Get()] = power
			case "[Engine]":
				engine := &Engine{
					Nickname:        semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					IdsName:         semantic.NewInt(section, "ids_name"),
					IdsInfo:         semantic.NewInt(section, "ids_info"),
					CruiseSpeed:     semantic.NewInt(section, "cruise_speed"),
					LinearDrag:      semantic.NewInt(section, "linear_drag"),
					MaxForce:        semantic.NewInt(section, "max_force"),
					ReverseFraction: semantic.NewFloat(section, "reverse_fraction", semantic.Precision(2)),

					HpType:           semantic.NewString(section, "hp_type", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					FlameEffect:      semantic.NewString(section, "flame_effect", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					TrailEffect:      semantic.NewString(section, "trail_effect", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					CruiseChargeTime: semantic.NewInt(section, "cruise_charge_time"),
				}
				frelconfig.Engines = append(frelconfig.Engines, engine)
				frelconfig.EnginesMap[engine.Nickname.Get()] = engine
			case "[Tractor]":
				tractor := &Tractor{
					Nickname:   semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					IdsName:    semantic.NewInt(section, "ids_name"),
					IdsInfo:    semantic.NewInt(section, "ids_info"),
					MaxLength:  semantic.NewInt(section, "max_length"),
					ReachSpeed: semantic.NewInt(section, "reach_speed"),
					Lootable:   semantic.NewBool(section, "lootable", semantic.StrBool),
				}
				frelconfig.Tractors = append(frelconfig.Tractors, tractor)
			case "[CounterMeasureDropper]":
				item := &CounterMeasureDropper{
					Nickname: semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					IdsName:  semantic.NewInt(section, "ids_name"),
					IdsInfo:  semantic.NewInt(section, "ids_info"),
					Lootable: semantic.NewBool(section, "lootable", semantic.StrBool),

					ProjectileArchetype: semantic.NewString(section, "projectile_archetype", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					HitPts:              semantic.NewInt(section, "hit_pts"),
					AIRange:             semantic.NewInt(section, "ai_range"),
				}
				frelconfig.CounterMeasureDroppers = append(frelconfig.CounterMeasureDroppers, item)
			case "[CounterMeasure]":
				item := &CounterMeasure{
					Nickname: semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					IdsName:  semantic.NewInt(section, "ids_name"),
					IdsInfo:  semantic.NewInt(section, "ids_info"),

					AmmoLimit:     semantic.NewInt(section, "ammo_limit"),
					Lifetime:      semantic.NewInt(section, "lifetime"),
					Range:         semantic.NewInt(section, "range"),
					DiversionPctg: semantic.NewInt(section, "diversion_pctg"),
				}
				frelconfig.CounterMeasure = append(frelconfig.CounterMeasure, item)
				frelconfig.CounterMeasureMap[item.Nickname.Get()] = item
			}
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
