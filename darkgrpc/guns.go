package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetGuns(_ context.Context, in *pb.GetGunsInput) (*pb.GetGunsReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}
	var input []configs_export.Gun
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulGun(s.app_data.Configs.Guns)
	} else {
		input = s.app_data.Configs.Guns
	}
	var items []*pb.Gun
	for _, item := range input {

		items = append(items, NewGun(item, in))
	}
	return &pb.GetGunsReply{Items: items}, nil
}

func (s *Server) GetMissiles(_ context.Context, in *pb.GetGunsInput) (*pb.GetGunsReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}
	var input []configs_export.Gun
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulGun(s.app_data.Configs.Missiles)
	} else {
		input = s.app_data.Configs.Missiles
	}
	var items []*pb.Gun
	for _, item := range input {

		items = append(items, NewGun(item, in))
	}
	return &pb.GetGunsReply{Items: items}, nil
}

func NewGun(item configs_export.Gun, in *pb.GetGunsInput) *pb.Gun {
	result := &pb.Gun{
		Nickname:               item.Nickname,
		Name:                   item.Name,
		Type:                   item.Type,
		Price:                  int64(item.Price),
		Class:                  item.Class,
		HpType:                 item.HpType,
		IdsName:                int64(item.IdsName),
		IdsInfo:                int64(item.IdsInfo),
		Volume:                 item.Volume,
		HitPts:                 item.HitPts,
		PowerUsage:             item.PowerUsage,
		Refire:                 item.Refire,
		Range:                  item.Range,
		Toughness:              item.Toughness,
		IsAutoTurret:           item.IsAutoTurret,
		Lootable:               item.Lootable,
		RequiredAmmo:           item.RequiredAmmo,
		HullDamage:             int64(item.HullDamage),
		EnergyDamage:           int64(item.EnergyDamage),
		ShieldDamage:           int64(item.ShieldDamage),
		AvgShieldDamage:        int64(item.AvgShieldDamage),
		DamageType:             item.DamageType,
		LifeTime:               item.LifeTime,
		Speed:                  item.Speed,
		GunTurnRate:            item.GunTurnRate,
		DispersionAngle:        item.DispersionAngle,
		HullDamagePerSec:       item.HullDamagePerSec,
		AvgShieldDamagePerSec:  item.AvgShieldDamagePerSec,
		EnergyDamagePerSec:     item.EnergyDamagePerSec,
		PowerUsagePerSec:       item.PowerUsagePerSec,
		AvgEfficiency:          item.AvgEfficiency,
		HullEfficiency:         item.HullEfficiency,
		ShieldEfficiency:       item.ShieldEfficiency,
		EnergyDamageEfficiency: item.EnergyDamageEfficiency,
		Mass:                   item.Mass,
		NumBarrels:             NewInt64(item.NumBarrels),
		Missile: &pb.Missile{
			MaxAngularVelocity: item.MaxAngularVelocity,
		},
		GunDetailed: &pb.GunDetailed{
			FlashParticleName: item.FlashParticleName,
			ConstEffect:       item.ConstEffect,
			MunitionHitEffect: item.MunitionHitEffect,
		},
		AmmoLimit: NewAmmoLimit(item.AmmoLimit),
	}
	if item.BurstFire != nil {
		result.BurstFire = &pb.BurstFire{
			SustainedRefire:                item.BurstFire.SustainedRefire,
			Ammo:                           int64(item.BurstFire.Ammo),
			ReloadTime:                     item.BurstFire.ReloadTime,
			SustainedHullDamagePerSec:      item.BurstFire.SustainedHullDamagePerSec,
			SustainedAvgShieldDamagePerSec: item.BurstFire.SustainedAvgShieldDamagePerSec,
			SustainedEnergyDamagePerSec:    item.BurstFire.SustainedEnergyDamagePerSec,
			SustainedPowerUsagePerSec:      item.BurstFire.SustainedPowerUsagePerSec,
		}
	}
	if item.DiscoGun != nil {
		result.DiscoGun = &pb.DiscoGun{
			ArmorPen: item.DiscoGun.ArmorPen,
		}
	}
	if in.IncludeDamageBonuses {
		for _, dmg_bonus := range item.DamageBonuses {
			result.DamageBonuses = append(result.DamageBonuses, &pb.DamageBonus{
				Type:     dmg_bonus.Type,
				Modifier: dmg_bonus.Modifier,
			})
		}
	}
	if in.IncludeMarketGoods {
		result.Bases = NewBases(item.Bases)
	}
	if in.IncludeTechCompat {
		result.DiscoveryTechCompat = NewTechCompat(item.DiscoveryTechCompat)
	}
	return result
}
