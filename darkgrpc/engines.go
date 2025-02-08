package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
)

func (s *Server) GetEngines(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetEnginesReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var items []*pb.Engine
	for _, item := range s.app_data.Configs.Engines {
		result := &pb.Engine{
			Name:             item.Name,
			Price:            int64(item.Price),
			CruiseSpeed:      int64(item.CruiseSpeed),
			CruiseChargeTime: int64(item.CruiseChargeTime),
			LinearDrag:       int64(item.LinearDrag),
			MaxForce:         int64(item.MaxForce),
			ReverseFraction:  item.ReverseFraction,
			ImpulseSpeed:     item.ImpulseSpeed,
			HpType:           item.HpType,
			FlameEffect:      item.FlameEffect,
			TrailEffect:      item.TrailEffect,
			Nickname:         item.Nickname,
			NameID:           int64(item.NameID),
			InfoID:           int64(item.InfoID),
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
	return &pb.GetEnginesReply{Items: items}, nil
}
