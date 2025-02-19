package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetShields(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetShieldsReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []configs_export.Shield
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulShields(s.app_data.Configs.Shields)
	} else {
		input = s.app_data.Configs.Shields
	}

	var items []*pb.Shield
	for _, item := range input {
		result := &pb.Shield{
			Name:              item.Name,
			Class:             item.Class,
			Type:              item.Type,
			Technology:        item.Technology,
			Price:             int64(item.Price),
			Capacity:          int64(item.Capacity),
			RegenerationRate:  int64(item.RegenerationRate),
			ConstantPowerDraw: int64(item.ConstantPowerDraw),
			Value:             item.Value,
			RebuildPowerDraw:  int64(item.RebuildPowerDraw),
			OffRebuildTime:    int64(item.OffRebuildTime),
			Toughness:         item.Toughness,
			HitPts:            int64(item.HitPts),
			Lootable:          item.Lootable,
			Nickname:          item.Nickname,
			HpType:            item.HpType,
			IdsName:           int64(item.IdsName),
			IdsInfo:           int64(item.IdsInfo),
			Mass:              item.Mass,
		}
		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases)
		}
		if in.IncludeTechCompat {
			result.DiscoveryTechCompat = NewTechCompat(item.DiscoveryTechCompat)
		}
		items = append(items, result)
	}
	return &pb.GetShieldsReply{Items: items}, nil
}
