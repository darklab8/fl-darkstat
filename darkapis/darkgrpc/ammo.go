package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetAmmos(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetAmmoReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []configs_export.Ammo
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulAmmo(s.app_data.Configs.Ammos)
	} else {
		input = s.app_data.Configs.Ammos
	}
	input = FilterNicknames(in.FilterNicknames, input)

	var items []*pb.Ammo
	for _, item := range input {
		result := &pb.Ammo{
			Name:             item.Name, //
			Price:            int64(item.Price),
			HitPts:           int64(item.HitPts),
			Volume:           item.Volume,
			MunitionLifetime: item.MunitionLifetime,
			Nickname:         item.Nickname,
			NameID:           int64(item.NameID),
			InfoID:           int64(item.InfoID),
			SeekerType:       item.SeekerType,
			SeekerRange:      int64(item.SeekerRange),
			SeekerFovDeg:     int64(item.SeekerFovDeg),
			AmmoLimit:        NewAmmoLimit(item.AmmoLimit),
			Mass:             item.Mass,
		}
		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases)
		}
		if in.IncludeTechCompat {
			result.DiscoveryTechCompat = NewTechCompat(item.DiscoveryTechCompat)
		}
		items = append(items, result)
	}
	return &pb.GetAmmoReply{Items: items}, nil
}

func NewAmmoLimit(AmmoLimit configs_export.AmmoLimit) *pb.AmmoLimit {
	return &pb.AmmoLimit{
		AmountInCatridge: NewInt64(AmmoLimit.AmountInCatridge),
		MaxCatridges:     NewInt64(AmmoLimit.MaxCatridges),
	}
}
