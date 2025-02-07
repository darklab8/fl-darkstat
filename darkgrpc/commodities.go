package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
)

func (s *Server) GetCommodities(_ context.Context, in *pb.Empty) (*pb.GetCommoditiesReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var items []*pb.Commodity
	for _, item := range s.app_data.Configs.Commodities {
		item := &pb.Commodity{
			Nickname:              item.Nickname,
			PriceBase:             int32(item.PriceBase),
			Name:                  item.Name,
			Combinable:            item.Combinable,
			Volume:                item.Volume,
			ShipClass:             int64(item.ShipClass),
			NameID:                int32(item.NameID),
			InfocardID:            int32(item.InfocardID),
			Bases:                 NewBases(item.Bases),
			PriceBestBaseBuysFor:  int32(item.PriceBestBaseBuysFor),
			PriceBestBaseSellsFor: int32(item.PriceBestBaseSellsFor),
			ProffitMargin:         int32(item.ProffitMargin),
			Mass:                  item.Mass,
		}

		items = append(items, item)
	}
	return &pb.GetCommoditiesReply{Items: items}, nil
}

func (s *Server) GetCommoditiesMarketGoods(_ context.Context, in *pb.GetMarketGoodsInput) (*pb.GetMarketGoodsReply, error) {
	return GetMarketGoods(s.app_data.Configs.Commodities, in)
}
