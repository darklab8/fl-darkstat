package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetScanners(_ context.Context, in *pb.GetEquipmentInput) (*pb.GetScannersReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []configs_export.Scanner
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUserfulScanners(s.app_data.Configs.Scanners)
	} else {
		input = s.app_data.Configs.Scanners
	}
	input = FilterNicknames(in.FilterNicknames, input)

	var items []*pb.Scanner
	for _, item := range input {
		result := &pb.Scanner{
			Name:           item.Name,
			Price:          int64(item.Price),
			Range:          int64(item.Range),
			CargoScanRange: int64(item.CargoScanRange),
			Lootable:       item.Lootable,
			Nickname:       item.Nickname,
			NameID:         int64(item.NameID),
			InfoID:         int64(item.InfoID),
			Mass:           item.Mass,
		}
		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases, in.FilterMarketGoodCategory)
		}
		if in.IncludeTechCompat {
			result.DiscoveryTechCompat = NewTechCompat(item.DiscoveryTechCompat)
		}
		items = append(items, result)
	}
	return &pb.GetScannersReply{Items: items}, nil
}
