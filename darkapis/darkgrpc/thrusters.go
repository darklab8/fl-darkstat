package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetThrusters(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetThrustersReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []configs_export.Thruster
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulThrusters(s.app_data.Configs.Thrusters)
	} else {
		input = s.app_data.Configs.Thrusters
	}
	input = FilterNicknames(in.FilterNicknames, input)

	var items []*pb.Thruster
	for _, item := range input {
		result := &pb.Thruster{
			Name:       item.Name,
			Price:      int64(item.Price),
			MaxForce:   int64(item.MaxForce),
			PowerUsage: int64(item.PowerUsage),
			Efficiency: item.Efficiency,
			Value:      item.Value,
			HitPts:     int64(item.HitPts),
			Lootable:   item.Lootable,
			Nickname:   item.Nickname,
			NameID:     int64(item.NameID),
			InfoID:     int64(item.InfoID),
			Mass:       item.Mass,
		}
		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases)
		}
		if in.IncludeTechCompat {
			result.DiscoveryTechCompat = NewTechCompat(item.DiscoveryTechCompat)
		}
		items = append(items, result)
	}
	return &pb.GetThrustersReply{Items: items}, nil
}
