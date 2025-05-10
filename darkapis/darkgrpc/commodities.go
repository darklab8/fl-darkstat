package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetCommodities(_ context.Context, in *pb.GetCommoditiesInput) (*pb.GetCommoditiesReply, error) {
	if s.app_data != nil {
		s.app_data.RLock()
		defer s.app_data.RUnlock()
	}

	var input []*configs_export.Commodity
	if in.FilterToUseful {
		input = s.app_data.Configs.FilterToUsefulCommodities(s.app_data.Configs.Commodities)
	} else {
		input = s.app_data.Configs.Commodities
	}
	input = FilterNicknames(in.FilterNicknames, input)

	var items []*pb.Commodity
	for _, item := range input {
		result := &pb.Commodity{
			Nickname:              item.Nickname,
			PriceBase:             int64(item.PriceBase),
			Name:                  item.Name,
			Combinable:            item.Combinable,
			Volume:                item.Volume,
			ShipClass:             NewShipClass(item.ShipClass),
			NameId:                int64(item.NameID),
			InfocardId:            int64(item.InfocardID),
			PriceBestBaseBuysFor:  int64(item.PriceBestBaseBuysFor),
			PriceBestBaseSellsFor: int64(item.PriceBestBaseSellsFor),
			ProffitMargin:         int64(item.ProffitMargin),
			Mass:                  item.Mass,
		}
		if in.IncludeMarketGoods {
			result.Bases = NewBases(item.Bases)
		}
		items = append(items, result)
	}
	return &pb.GetCommoditiesReply{Items: items}, nil
}
