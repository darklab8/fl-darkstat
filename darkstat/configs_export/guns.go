package configs_export

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-configs/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type DamageBonus struct {
	Type     string
	Modifier float64
}

type GunDetailed struct {
	FlashParticleName string
	ConstEffect       string
	MunitionHitEffect string
}

func (g Gun) GetNickname() string                 { return g.Nickname }
func (g Gun) GetTechCompat() *DiscoveryTechCompat { return g.DiscoveryTechCompat }

type Gun struct {
	Nickname     string
	NicknameHash flhash.HashCode
	MunitionHash flhash.HashCode
	HpTypeHash   flhash.HashCode
	Name         string
	Type         string
	Price        int
	Class        string
	HpType       string
	IdsName      int
	IdsInfo      int
	Volume       float64

	HitPts       string
	PowerUsage   float64
	Refire       float64
	Range        float64
	Toughness    float64
	IsAutoTurret bool
	Lootable     bool

	RequiredAmmo bool
	// AmmoPrice     int
	// AmmoBases     []*GoodAtBase
	// AmmoName      string
	HullDamage      int
	EnergyDamage    int
	ShieldDamage    int
	AvgShieldDamage int
	DamageType      string
	LifeTime        float64
	Speed           float64
	GunTurnRate     float64
	DispersionAngle float64

	HullDamagePerSec       float64
	AvgShieldDamagePerSec  float64
	EnergyDamagePerSec     float64
	PowerUsagePerSec       float64
	AvgEfficiency          float64
	HullEfficiency         float64
	ShieldEfficiency       float64
	EnergyDamageEfficiency float64

	Bases         map[cfgtype.BaseUniNick]*GoodAtBase
	DamageBonuses []DamageBonus

	Missile
	*DiscoveryTechCompat
	GunDetailed

	NumBarrels *int
	BurstFire  *BurstFire
	AmmoLimit  AmmoLimit

	Mass float64

	DiscoGun *DiscoGun
}

type DiscoGun struct {
	ArmorPen float64
}

type BurstFire struct {
	SustainedRefire float64
	Ammo            int
	ReloadTime      float64

	SustainedHullDamagePerSec      float64
	SustainedAvgShieldDamagePerSec float64
	SustainedEnergyDamagePerSec    float64
	SustainedPowerUsagePerSec      float64
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

func (e *Exporter) getGunInfo(gun_info *equip_mapped.Gun, ids []Tractor, buyable_ship_tech map[string]bool) (Gun, error) {
	gun_nickname := gun_info.Nickname.Get()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			fmt.Println("recovered gun_nickname", gun_nickname)
			panic(r)
		}
	}()

	gun := Gun{
		Nickname:     gun_nickname,
		NicknameHash: flhash.HashNickname(gun_nickname),
		IdsName:      gun_info.IdsName.Get(),
		IdsInfo:      gun_info.IdsInfo.Get(),
		Class:        getGunClass(gun_info),
		HitPts:       gun_info.HitPts.Get(),
		PowerUsage:   gun_info.PowerUsage.Get(),
		Refire:       float64(1 / gun_info.RefireDelay.Get()),
		Speed:        gun_info.MuzzleVelosity.Get(),
		Toughness:    gun_info.Toughness.Get(),
		Lootable:     gun_info.Lootable.Get(),
		Bases:        make(map[cfgtype.BaseUniNick]*GoodAtBase),
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
	gun.HpTypeHash = flhash.HashNickname(gun.HpType)

	munition, found_munition := e.Configs.Equip.MunitionMap[gun_info.ProjectileArchetype.Get()]

	if e.Configs.FLSR != nil && !found_munition && gun.Nickname == "gd_ww_turret_laser_light02" && gun_info.ProjectileArchetype.Get() == "gd_ww_laser_light02_ammo" {
		logus.Log.Error("gun does not have defined munition",
			typelog.Any("nickname", gun.Nickname),
			typelog.Any("projectile_archetype", gun_info.ProjectileArchetype.Get()))
		return gun, errors.New("not defined munition")
	}

	if gun.Nickname == "gd_ww_turret_laser_light02" {
		logus.Log.Warn("FLSR broken gun potentially",
			typelog.String("gun.Nickname", gun.Nickname),
			typelog.String("projectile", gun_info.ProjectileArchetype.Get()),
			typelog.Bool("is_flsr", e.Configs.FLSR != nil),
			typelog.Bool("found_munition", found_munition),
		)
	}

	gun.MunitionHash = flhash.HashNickname(munition.Nickname.Get())
	gun.FlashParticleName, _ = gun_info.FlashParticleName.GetValue()
	gun.ConstEffect, _ = munition.ConstEffect.GetValue()
	gun.MunitionHitEffect, _ = munition.MunitionHitEffect.GetValue()

	if hull_damange, ok := munition.HullDamage.GetValue(); ok {
		// regular gun or turret
		gun.HullDamage = hull_damange
		gun.EnergyDamage = munition.EnergyDamange.Get()
	} else {

		if explosion_arch, ok := munition.ExplosionArch.GetValue(); ok {
			// rocket launcher
			explosion := e.Configs.Equip.ExplosionMap[explosion_arch]
			gun.HullDamage = explosion.HullDamage.Get()
			gun.EnergyDamage = explosion.EnergyDamange.Get()
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

	gun.DamageType = "undefined"
	if weapon_type, ok := munition.WeaponType.GetValue(); ok {
		gun.DamageType = weapon_type
	}

	if lifetime, ok := munition.LifeTime.GetValue(); ok {
		gun.LifeTime = lifetime
	} else {
		gun.LifeTime = 100000000
	}
	gun.Range = gun.LifeTime * gun.Speed

	if weapon_type, ok := e.Configs.WeaponMods.WeaponTypesMap[gun.DamageType]; ok {
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
	if good_info, ok := e.Configs.Goods.GoodsMap[gun.Nickname]; ok {
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

	e.exportInfocards(InfocardKey(gun.Nickname), gun.IdsInfo)

	gun.ShieldDamage = int(float64(gun.HullDamage)*float64(e.Configs.Consts.ShieldEquipConsts.HULL_DAMAGE_FACTOR.Get()) + float64(gun.EnergyDamage))

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

	gun.AvgEfficiency = (float64(gun.HullDamage) + float64(gun.AvgShieldDamage)) / (gun.PowerUsage * 2)
	gun.HullEfficiency = float64(gun.HullDamage) / gun.PowerUsage
	gun.ShieldEfficiency = float64(gun.AvgShieldDamage) / gun.PowerUsage
	gun.EnergyDamageEfficiency = float64(gun.EnergyDamage) / gun.PowerUsage
	gun.GunTurnRate, _ = gun_info.TurnRate.GetValue()
	gun.DispersionAngle, _ = gun_info.DispersionAngle.GetValue()

	if gun.IsAutoTurret {
		gun.Type = "turret"
	} else {
		gun.Type = "gun"
	}

	// fmt.Println("CalculateTEchCompat", e.configs.Discovery != nil, gun.Nickname)
	gun.DiscoveryTechCompat = CalculateTechCompat(e.Configs.Discovery, ids, gun.Nickname)

	if e.Configs.Discovery != nil {
		gun.DiscoGun = &DiscoGun{}
		if armor_pen, ok := munition.ArmorPen.GetValue(); ok {
			gun.DiscoGun.ArmorPen = armor_pen
		}

		if explosion_arch, ok := munition.ExplosionArch.GetValue(); ok {
			// rocket launcher
			explosion := e.Configs.Equip.ExplosionMap[explosion_arch]
			if armor_pen, ok := explosion.ArmorPen.GetValue(); ok {
				gun.DiscoGun.ArmorPen = armor_pen
			}
		}
	}

	e.Hashes[gun.Nickname] = gun.NicknameHash
	e.Hashes[munition.Nickname.Get()] = gun.MunitionHash
	e.Hashes[gun.HpType] = gun.HpTypeHash

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

func (e *Exporter) GetGuns(ids []Tractor, buyable_ship_tech map[string]bool) []Gun {
	var guns []Gun

	for _, gun_info := range e.Configs.Equip.Guns {
		gun, err := e.getGunInfo(gun_info, ids, buyable_ship_tech)

		if err != nil {
			continue
		}

		munition := e.Configs.Equip.MunitionMap[gun_info.ProjectileArchetype.Get()]
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
