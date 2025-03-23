package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetTractors(_ context.Context, in *pb.GetTractorsInput) (*pb.GetTractorsReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []*configs_export.Tractor
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulTractors(s.app_data.Configs.Tractors)
	} else {
		input = s.app_data.Configs.Tractors
	}
	input = FilterNicknames(in.FilterNicknames, input)

	var items []*pb.Tractor
	for _, item := range input {
		result := &pb.Tractor{
			Name:       item.Name, //
			Price:      int64(item.Price),
			MaxLength:  int64(item.MaxLength),
			ReachSpeed: int64(item.ReachSpeed),
			Lootable:   item.Lootable,
			Nickname:   string(item.Nickname),
			NameId:     int64(item.NameID),
			InfoId:     int64(item.InfoID),
			Mass:       item.Mass,
		}
		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases)
		}
		if in.IncludeRephacks {
			rephacks := item.GetRephacksList()
			for _, rephack := range rephacks {
				result.Rephacks = append(result.Rephacks, &pb.Rephack{
					FactionName:     rephack.FactionName,
					FactionNickname: string(rephack.FactionNick),
					Reputation:      rephack.Reputation,
					RepType:         int32(rephack.RepType),
				})
			}

		}
		items = append(items, result)
	}
	return &pb.GetTractorsReply{Items: items}, nil
}
