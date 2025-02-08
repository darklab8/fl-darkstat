package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
)

func (s *Server) GetAmmos(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetAmmoReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var items []*pb.Ammo
	for _, item := range s.app_data.Configs.Ammos {
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
			AmountInCatridge: NewInt64(item.AmmoLimit.AmountInCatridge),
			MaxCatridges:     NewInt64(item.AmmoLimit.MaxCatridges),
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
