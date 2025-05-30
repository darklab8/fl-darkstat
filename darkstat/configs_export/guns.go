package configs_export

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type DamageBonus struct {
	Type     string  `json:"type" validate:"required"`
	Modifier float64 `json:"modifier" validate:"required"`
}

func (g Gun) GetTechCompat() *DiscoveryTechCompat { return g.DiscoveryTechCompat }

type Gun struct {
	Nickname string  `json:"nickname"  validate:"required"`
	Name     string  `json:"name"  validate:"required"`
	Type     string  `json:"type"  validate:"required"`
	Price    int     `json:"price"  validate:"required"`
	Class    string  `json:"class"  validate:"required"`
	HpType   string  `json:"hp_type"  validate:"required"`
	IdsName  int     `json:"ids_name"  validate:"required"`
	IdsInfo  int     `json:"ids_info" validate:"required"`
	Volume   float64 `json:"volume" validate:"required"`

	HitPts       string  `json:"hit_pts"  validate:"required"`
	PowerUsage   float64 `json:"power_usage"  validate:"required"`
	Refire       float64 `json:"refire" validate:"required"`
	Range        float64 `json:"range"  validate:"required"`
	Toughness    float64 `json:"toughness"  validate:"required"`
	IsAutoTurret bool    `json:"is_auto_turret"  validate:"required"`
	Lootable     bool    `json:"lootable"  validate:"required"`

	RequiredAmmo bool `json:"required_ammo"  validate:"required"`
	// AmmoPrice     int
	// AmmoBases     []*GoodAtBase
	// AmmoName      string
	HullDamage      int     `json:"hull_damage"  validate:"required"`
	EnergyDamage    int     `json:"energy_damage"  validate:"required"`
	ShieldDamage    int     `json:"shield_damage" validate:"required"`
	AvgShieldDamage int     `json:"avg_shield_damage"  validate:"required"`
	DamageType      string  `json:"damage_type"  validate:"required"`
	LifeTime        float64 `json:"life_time" validate:"required"`
	Speed           float64 `json:"speed" validate:"required"`
	GunTurnRate     float64 `json:"gun_turn_rate" validate:"required"`
	DispersionAngle float64 `json:"dispersion_angle" validate:"required"`

	HullDamagePerSec       float64 `json:"hull_damage_per_sec"  validate:"required"`
	AvgShieldDamagePerSec  float64 `json:"avg_shield_damage_per_sec" validate:"required"`
	EnergyDamagePerSec     float64 `json:"energy_damage_per_sec" validate:"required"`
	PowerUsagePerSec       float64 `json:"power_usage_per_sec" validate:"required"`
	AvgEfficiency          float64 `json:"avg_efficiency" validate:"required"`
	HullEfficiency         float64 `json:"hull_efficiency" validate:"required"`
	ShieldEfficiency       float64 `json:"shield_efficiency" validate:"required"`
	EnergyDamageEfficiency float64 `json:"energy_damage_efficiency" validate:"required"`

	Bases         map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`
	DamageBonuses []DamageBonus                   `json:"damage_bonuses" validate:"required"`

	Missile
	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`

	NumBarrels *int       `json:"num_barrels,omitempty"`
	BurstFire  *BurstFire `json:"burst_fire,omitempty"`
	AmmoLimit  AmmoLimit  `json:"ammo_limit,omitempty"`

	Mass float64 `json:"mass" validate:"required"`

	DiscoGun *DiscoGun `json:"disco_gun"`
}

func (b Gun) GetNickname() string { return string(b.Nickname) }

func (b Gun) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Gun) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

type DiscoGun struct {
	ArmorPen float64 `json:"armor_pen" validate:"required"`
}

type BurstFire struct {
	SustainedRefire float64 `json:"sustained_fire" validate:"required"`
	Ammo            int     `json:"ammo" validate:"required"`
	ReloadTime      float64 `json:"reload_time" validate:"required"`

	SustainedHullDamagePerSec      float64 `json:"sustained_hull_dmg_per_sec" validate:"required"`
	SustainedAvgShieldDamagePerSec float64 `json:"sustained_avg_shield_dmg_per_sec" validate:"required"`
	SustainedEnergyDamagePerSec    float64 `json:"sustained_energy_dmg_per_sec" validate:"required"`
	SustainedPowerUsagePerSec      float64 `json:"sustained_pwer_usage_per_sec" validate:"required"`
}

func getGunClass(gun_info *equip_mapped.Gun) string {
	var gun_class string
	if gun_type, ok := gun_info.HPGunType.GetValue(); ok {
		splitted := strings.Split(gun_type, "_")
		if len(splitted) > 0 {
			class := splitted[len(splitted)-1]
			if _, err := strconv.Atoi(class); err == nil {
				gun_class = class
			}
		}
	}
	return gun_class
}

func (e *Exporter) getGunInfo(gun_info *equip_mapped.Gun, ids []*Tractor, buyable_ship_tech map[string]bool) (Gun, error) {
	gun_nickname := gun_info.Nickname.Get()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			fmt.Println("recovered gun_nickname", gun_nickname)
			panic(r)
		}
	}()

	gun := Gun{
		Nickname:   gun_nickname,
		IdsName:    gun_info.IdsName.Get(),
		IdsInfo:    gun_info.IdsInfo.Get(),
		Class:      getGunClass(gun_info),
		HitPts:     gun_info.HitPts.Get(),
		PowerUsage: gun_info.PowerUsage.Get(),
		Refire:     float64(1 / gun_info.RefireDelay.Get()),
		Speed:      gun_info.MuzzleVelosity.Get(),
		Toughness:  gun_info.Toughness.Get(),
		Lootable:   gun_info.Lootable.Get(),
		Bases:      make(map[cfg.BaseUniNick]*MarketGood),
	}

	gun.Mass, _ = gun_info.Mass.GetValue()

	num_barrels := 1
	if num_barrels_value, ok := gun_info.NumBarrels.GetValue(); ok {
		num_barrels = num_barrels_value
		gun.NumBarrels = ptr.Ptr(num_barrels_value)
	}

	if ammo, ok := gun_info.BurstAmmo.GetValue(); ok {
		gun.BurstFire = &BurstFire{
			Ammo: ammo,
		}
		gun.BurstFire.ReloadTime = gun_info.BurstReload.Get()

		// (magCapacity * RefireDelay + Reload time) / Mag Capacity = This should be average refire delay

		gun.BurstFire.SustainedRefire = 1 / ((gun_info.RefireDelay.Get()*float64(gun.BurstFire.Ammo) + gun.BurstFire.ReloadTime) / float64(gun.BurstFire.Ammo))
	}

	gun.IsAutoTurret, _ = gun_info.IsAutoTurret.GetValue()
	gun.Volume, _ = gun_info.Volume.GetValue()

	gun.HpType, _ = gun_info.HPGunType.GetValue()

	munition, found_munition := e.Mapped.Equip().MunitionMap[gun_info.ProjectileArchetype.Get()]

	if e.Mapped.FLSR != nil && !found_munition && gun.Nickname == "gd_ww_turret_laser_light02" && gun_info.ProjectileArchetype.Get() == "gd_ww_laser_light02_ammo" {
		logus.Log.Error("gun does not have defined munition",
			typelog.Any("nickname", gun.Nickname),
			typelog.Any("projectile_archetype", gun_info.ProjectileArchetype.Get()))
		return gun, errors.New("not defined munition")
	}

	if gun.Nickname == "gd_ww_turret_laser_light02" {
		logus.Log.Warn("FLSR broken gun potentially",
			typelog.String("gun.Nickname", gun.Nickname),
			typelog.String("projectile", gun_info.ProjectileArchetype.Get()),
			typelog.Bool("is_flsr", e.Mapped.FLSR != nil),
			typelog.Bool("found_munition", found_munition),
		)
	}

	gun.DamageType = "undefined"

	if hull_damange, ok := munition.HullDamage.GetValue(); ok {
		// regular gun or turret
		gun.HullDamage = hull_damange
		gun.EnergyDamage = munition.EnergyDamange.Get()
		if weapon_type, ok := munition.WeaponType.GetValue(); ok {
			gun.DamageType = weapon_type
		}
	} else {

		if explosion_arch, ok := munition.ExplosionArch.GetValue(); ok {
			// rocket launcher
			explosion := e.Mapped.Equip().ExplosionMap[explosion_arch]
			gun.HullDamage = explosion.HullDamage.Get()
			gun.EnergyDamage = explosion.EnergyDamange.Get()
			if weapon_type, ok := explosion.WeaponType.GetValue(); ok {
				gun.DamageType = weapon_type
			}
		} else {
			// healing gun
			gun.HullDamage = -1
		}

	}

	if required_ammo, ok := munition.RequiredAmmo.GetValue(); ok {
		gun.RequiredAmmo = required_ammo
	}

	if value, ok := munition.AmmoLimitAmountInCatridge.GetValue(); ok {
		gun.AmmoLimit.AmountInCatridge = ptr.Ptr(value)
	}
	if value, ok := munition.AmmoLimitMaxCatridges.GetValue(); ok {
		gun.AmmoLimit.MaxCatridges = ptr.Ptr(value)
	}

	if lifetime, ok := munition.LifeTime.GetValue(); ok {
		gun.LifeTime = lifetime
	} else {
		gun.LifeTime = 100000000
	}
	gun.Range = gun.LifeTime * gun.Speed

	if weapon_type, ok := e.Mapped.WeaponMods.WeaponTypesMap[gun.DamageType]; ok {
		for _, weapon_modifier := range weapon_type.ShieldMods {
			gun.DamageBonuses = append(gun.DamageBonuses,
				DamageBonus{
					Type:     weapon_modifier.ShieldType.Get(),
					Modifier: weapon_modifier.DamageModifier.Get(),
				},
			)
		}
	}

	gun.Price = -1
	if good_info, ok := e.Mapped.Goods.GoodsMap[gun.Nickname]; ok {
		if price, ok := good_info.Price.GetValue(); ok {
			gun.Price = price
			gun.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
				Nickname: good_info.Nickname.Get(),
				Price:    price,
			})
		}
	}

	gun.Name = e.GetInfocardName(gun.IdsName, gun.Nickname)

	if NameWithSpacesOnly(gun.Name) {
		gun.Name = "Undefined"
	}

	e.exportInfocards(infocarder.InfocardKey(gun.Nickname), gun.IdsInfo)

	gun.ShieldDamage = int(float64(gun.HullDamage)*float64(e.Mapped.Consts.ShieldEquipConsts.HULL_DAMAGE_FACTOR.Get()) + float64(gun.EnergyDamage))

	avg_shield_modifier := 0.0
	shield_modifier_count := 0
	for _, damage_bonus := range gun.DamageBonuses {
		if _, ok := buyable_ship_tech[damage_bonus.Type]; !ok {
			continue
		}
		avg_shield_modifier += damage_bonus.Modifier
		shield_modifier_count += 1
	}
	if shield_modifier_count == 0 {
		shield_modifier_count = 1
		avg_shield_modifier = 1
	}

	avgShieldModifier := avg_shield_modifier / float64(shield_modifier_count)
	gun.AvgShieldDamage = int(float64(gun.ShieldDamage) * avgShieldModifier)

	gun.HullDamagePerSec = float64(gun.HullDamage) * gun.Refire * float64(num_barrels)
	gun.EnergyDamagePerSec = float64(gun.EnergyDamage) * gun.Refire * float64(num_barrels)
	gun.AvgShieldDamagePerSec = float64(gun.AvgShieldDamage) * gun.Refire * float64(num_barrels)

	gun.PowerUsagePerSec = float64(gun.PowerUsage) * gun.Refire * float64(num_barrels)

	if gun.BurstFire != nil {
		gun.BurstFire.SustainedHullDamagePerSec = float64(gun.HullDamage) * gun.BurstFire.SustainedRefire * float64(num_barrels)
		gun.BurstFire.SustainedEnergyDamagePerSec = float64(gun.EnergyDamage) * gun.BurstFire.SustainedRefire * float64(num_barrels)
		gun.BurstFire.SustainedAvgShieldDamagePerSec = float64(gun.AvgShieldDamage) * gun.BurstFire.SustainedRefire * float64(num_barrels)
		gun.BurstFire.SustainedPowerUsagePerSec = float64(gun.PowerUsage) * gun.BurstFire.SustainedRefire * float64(num_barrels)
	}

	power_usage_for_calcs := gun.PowerUsage
	if power_usage_for_calcs == 0 {
		power_usage_for_calcs = 1
	}
	gun.AvgEfficiency = (float64(gun.HullDamage) + float64(gun.AvgShieldDamage)) / (power_usage_for_calcs * 2)
	gun.HullEfficiency = float64(gun.HullDamage) / power_usage_for_calcs
	gun.ShieldEfficiency = float64(gun.AvgShieldDamage) / power_usage_for_calcs
	gun.EnergyDamageEfficiency = float64(gun.EnergyDamage) / power_usage_for_calcs

	gun.GunTurnRate, _ = gun_info.TurnRate.GetValue()
	gun.DispersionAngle, _ = gun_info.DispersionAngle.GetValue()

	if gun.IsAutoTurret {
		gun.Type = "turret"
	} else {
		gun.Type = "gun"
	}

	// fmt.Println("CalculateTEchCompat", e.mapped.Discovery != nil, gun.Nickname)
	gun.DiscoveryTechCompat = CalculateTechCompat(e.Mapped.Discovery, ids, gun.Nickname)

	if e.Mapped.Discovery != nil {
		gun.DiscoGun = &DiscoGun{}
		if armor_pen, ok := munition.ArmorPen.GetValue(); ok {
			gun.DiscoGun.ArmorPen = armor_pen
		}

		if explosion_arch, ok := munition.ExplosionArch.GetValue(); ok {
			// rocket launcher
			explosion := e.Mapped.Equip().ExplosionMap[explosion_arch]
			if armor_pen, ok := explosion.ArmorPen.GetValue(); ok {
				gun.DiscoGun.ArmorPen = armor_pen
			}
		}
	}

	return gun, nil
}

func (e *Exporter) GetBuyableShields(shields []Shield) map[string]bool {
	var buyable_ship_tech map[string]bool = make(map[string]bool)
	for _, shield := range shields {
		if !e.Buyable(shield.Bases) {
			continue
		}
		buyable_ship_tech[shield.Technology] = true
	}
	return buyable_ship_tech
}

func (e *Exporter) GetGuns(ids []*Tractor, buyable_ship_tech map[string]bool) []Gun {
	var guns []Gun

	for _, gun_info := range e.Mapped.Equip().Guns {
		gun, err := e.getGunInfo(gun_info, ids, buyable_ship_tech)

		if err != nil {
			continue
		}

		munition := e.Mapped.Equip().MunitionMap[gun_info.ProjectileArchetype.Get()]
		if _, ok := munition.Motor.GetValue(); ok {
			// Excluded rocket launching stuff
			continue
		}

		guns = append(guns, gun)
	}

	return guns
}

func (e *Exporter) FilterToUsefulGun(guns []Gun) []Gun {
	var items []Gun = make([]Gun, 0, len(guns))
	for _, gun := range guns {

		if gun.HpType == "" {
			continue
		}

		if strings.Contains(gun.DamageType, "w_npc") || strings.Contains(gun.DamageType, "station") {
			continue
		}
		if strings.Contains(gun.Name, "Obsolete Equipment") {
			continue
		}
		if strings.Contains(gun.Nickname, "_wp_") ||
			strings.Contains(gun.Nickname, "_wps_") ||
			strings.Contains(gun.Nickname, "_station_") ||
			strings.Contains(gun.Nickname, "admin_cannon") {
			continue
		}

		if !e.Buyable(gun.Bases) {
			continue
		}
		items = append(items, gun)
	}
	return items
}
