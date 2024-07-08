package configs_export

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
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

type Gun struct {
	Nickname string
	Name     string
	Type     string
	Price    int
	Class    string
	HpType   string
	IdsName  int
	IdsInfo  int

	HitPts       string
	PowerUsage   float64
	PowerPerSec  float64
	Refire       float64
	Range        float64
	Toughness    float64
	IsAutoTurret bool
	TurnRate     float64
	Lootable     bool

	RequiredAmmo bool
	// AmmoPrice     int
	// AmmoBases     []*GoodAtBase
	// AmmoName      string
	HullDamage      int
	EnergyDamange   int
	ShieldDamage    int
	AvgShieldDamage int
	DamageType      string
	LifeTime        float64
	Speed           float64

	HullDamagePerSec      float64
	AvgShieldDamagePerSec float64
	PowerUsagePerSec      float64
	AvgEfficiency         float64
	HullEfficiency        float64
	ShieldEfficiency      float64
	Value                 float64
	Rating                float64

	Bases         []*GoodAtBase
	DamageBonuses []DamageBonus

	Missile
	*DiscoveryTechCompat
	GunDetailed
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

func (e *Exporter) getGunInfo(gun_info *equip_mapped.Gun, ids []Tractor, buyable_ship_tech map[string]bool) Gun {
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
		TurnRate:   gun_info.TurnRate.Get(),
		Lootable:   gun_info.Lootable.Get(),
	}

	gun.IsAutoTurret, _ = gun_info.IsAutoTurret.GetValue()

	gun.HpType, _ = gun_info.HPGunType.GetValue()

	gun.PowerPerSec = gun.PowerUsage * gun.Refire
	munition := e.configs.Equip.MunitionMap[gun_info.ProjectileArchetype.Get()]

	gun.FlashParticleName, _ = gun_info.FlashParticleName.GetValue()
	gun.ConstEffect, _ = munition.ConstEffect.GetValue()
	gun.MunitionHitEffect, _ = munition.MunitionHitEffect.GetValue()

	if hull_damange, ok := munition.HullDamage.GetValue(); ok {
		// regular gun or turret
		gun.HullDamage = hull_damange
		gun.EnergyDamange = munition.EnergyDamange.Get()
	} else {

		if explosion_arch, ok := munition.ExplosionArch.GetValue(); ok {
			// rocket launcher
			explosion := e.configs.Equip.ExplosionMap[explosion_arch]
			gun.HullDamage = explosion.HullDamage.Get()
			gun.EnergyDamange = explosion.EnergyDamange.Get()
		} else {
			// healing gun
			gun.HullDamage = -1
		}

	}

	if required_ammo, ok := munition.RequiredAmmo.GetValue(); ok {
		gun.RequiredAmmo = required_ammo
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

	if weapon_type, ok := e.configs.WeaponMods.WeaponTypesMap[gun.DamageType]; ok {
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
	if good_info, ok := e.configs.Goods.GoodsMap[gun.Nickname]; ok {
		if price, ok := good_info.Price.GetValue(); ok {
			gun.Price = price
			gun.Bases = e.GetAtBasesSold(GetAtBasesInput{
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

	gun.ShieldDamage = int(float64(gun.HullDamage)*float64(e.configs.Consts.ShieldEquipConsts.HULL_DAMAGE_FACTOR.Get()) + float64(gun.EnergyDamange))

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
		shield_modifier_count += 1
	}

	avgShieldModifier := avg_shield_modifier / float64(shield_modifier_count)
	gun.AvgShieldDamage = int(float64(gun.ShieldDamage) * avgShieldModifier)

	gun.HullDamagePerSec = float64(gun.HullDamage) * gun.Refire
	gun.AvgShieldDamagePerSec = float64(gun.AvgShieldDamage) * gun.Refire
	gun.PowerUsagePerSec = float64(gun.PowerUsage) * gun.Refire
	gun.AvgEfficiency = (float64(gun.HullDamage) + float64(gun.AvgShieldDamage)) / (gun.PowerUsage * 2)
	gun.HullEfficiency = float64(gun.HullDamage) / gun.PowerUsage
	gun.ShieldEfficiency = float64(gun.AvgShieldDamage) / gun.PowerUsage
	gun.Value = math.Max(float64(gun.HullDamagePerSec), float64(gun.AvgShieldDamagePerSec)) / float64(gun.Price) * 1000
	gun.Rating = gun.AvgEfficiency * gun.Value

	if gun.IsAutoTurret {
		gun.Type = "turret"
	} else {
		gun.Type = "gun"
	}

	// fmt.Println("CalculateTEchCompat", e.configs.Discovery != nil, gun.Nickname)
	gun.DiscoveryTechCompat = CalculateTechCompat(e.configs.Discovery, ids, gun.Nickname)
	return gun
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

	for _, gun_info := range e.configs.Equip.Guns {
		gun := e.getGunInfo(gun_info, ids, buyable_ship_tech)

		munition := e.configs.Equip.MunitionMap[gun_info.ProjectileArchetype.Get()]
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
