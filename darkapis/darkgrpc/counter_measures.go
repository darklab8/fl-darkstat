package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetCounterMeasures(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetCounterMeasuresReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []configs_export.CounterMeasure
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulCounterMeasures(s.app_data.Configs.CMs)
	} else {
		input = s.app_data.Configs.CMs
	}
	input = FilterNicknames(in.FilterNicknames, input)

	var items []*pb.CounterMeasure
	for _, item := range input {
		result := &pb.CounterMeasure{
			Name:          item.Name,
			Price:         int64(item.Price),
			HitPts:        int64(item.HitPts),
			AIRange:       int64(item.AIRange),
			Lifetime:      int64(item.Lifetime),
			Range:         int64(item.Range),
			DiversionPctg: int64(item.DiversionPctg),
			Lootable:      item.Lootable,
			Nickname:      item.Nickname,
			NameID:        int64(item.NameID),
			InfoID:        int64(item.InfoID),
			AmmoLimit:     NewAmmoLimit(item.AmmoLimit),
			Mass:          item.Mass,
		}
		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases)
		}
		if in.IncludeTechCompat {
			result.DiscoveryTechCompat = NewTechCompat(item.DiscoveryTechCompat)
		}
		items = append(items, result)
	}
	return &pb.GetCounterMeasuresReply{Items: items}, nil
}
