package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetMines(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetMinesReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []configs_export.Mine
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulMines(s.app_data.Configs.Mines)
	} else {
		input = s.app_data.Configs.Mines
	}

	var items []*pb.Mine
	for _, item := range input {
		result := &pb.Mine{
			Name:                item.Name,
			Price:               int64(item.Price),
			AmmoPrice:           int64(item.AmmoPrice),
			Nickname:            item.Nickname,
			ProjectileArchetype: item.ProjectileArchetype,
			IdsName:             int64(item.IdsName),
			IdsInfo:             int64(item.IdsInfo),
			HullDamage:          int64(item.HullDamage),
			EnergyDamange:       int64(item.EnergyDamange),
			ShieldDamage:        int64(item.ShieldDamage),
			PowerUsage:          item.PowerUsage,
			Value:               item.Value,
			Refire:              item.Refire,
			DetonationDistance:  item.DetonationDistance,
			Radius:              item.Radius,
			SeekDistance:        int64(item.SeekDistance),
			TopSpeed:            int64(item.TopSpeed),
			Acceleration:        int64(item.Acceleration),
			LinearDrag:          item.LinearDrag,
			LifeTime:            item.LifeTime,
			OwnerSafe:           int64(item.OwnerSafe),
			Toughness:           item.Toughness,
			HitPts:              int64(item.HitPts),
			Lootable:            item.Lootable,
			AmmoLimit:           NewAmmoLimit(item.AmmoLimit),
			Mass:                item.Mass,
		}
		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases)
		}
		if in.IncludeTechCompat {
			result.DiscoveryTechCompat = NewTechCompat(item.DiscoveryTechCompat)
		}
		items = append(items, result)
	}
	return &pb.GetMinesReply{Items: items}, nil
}
