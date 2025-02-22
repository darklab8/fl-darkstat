package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetBasesNpc(_ context.Context, in *pb.GetBasesInput) (*pb.GetBasesReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var bases []*pb.Base
	var input []*configs_export.Base
	if in.FilterToUseful {
		input = configs_export.FilterToUserfulBases(s.app_data.Configs.Bases)
	} else {
		input = s.app_data.Configs.Bases
	}
	input = FilterNicknames(in.FilterNicknames, input)

	for _, base := range input {
		bases = append(bases, NewBase(base, in.IncludeMarketGoods, in.FilterMarketGoodCategory))
	}
	return &pb.GetBasesReply{Items: bases}, nil
}

func (s *Server) GetBasesMiningOperations(_ context.Context, in *pb.GetBasesInput) (*pb.GetBasesReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var bases []*pb.Base
	var input []*configs_export.Base
	if in.FilterToUseful {
		input = configs_export.FilterToUserfulBases(s.app_data.Configs.MiningOperations)
	} else {
		input = s.app_data.Configs.MiningOperations
	}
	input = FilterNicknames(in.FilterNicknames, input)

	for _, base := range input {
		bases = append(bases, NewBase(base, in.IncludeMarketGoods, in.FilterMarketGoodCategory))
	}
	return &pb.GetBasesReply{Items: bases}, nil
}

func (s *Server) GetBasesPoBs(_ context.Context, in *pb.GetBasesInput) (*pb.GetBasesReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []*configs_export.Base = s.app_data.Configs.PoBsToBases(s.app_data.Configs.PoBs)
	input = FilterNicknames(in.FilterNicknames, input)
	var bases []*pb.Base
	for _, base := range input {
		bases = append(bases, NewBase(base, in.IncludeMarketGoods, in.FilterMarketGoodCategory))
	}
	return &pb.GetBasesReply{Items: bases}, nil
}

func NewBase(base *configs_export.Base, include_market_goods bool, filter_market_good_category []string) *pb.Base {
	item := &pb.Base{
		Name:                   base.Name,
		Archetypes:             base.Archetypes,
		Nickname:               string(base.Nickname),
		FactionName:            base.FactionName,
		System:                 base.System,
		SystemNickname:         base.SystemNickname,
		Region:                 base.Region,
		StridName:              int64(base.StridName),
		InfocardId:             int64(base.InfocardID),
		File:                   base.File.ToString(),
		BgcsBaseRunBy:          base.BGCS_base_run_by,
		Pos:                    NewPos(&base.Pos),
		SectorCoord:            base.SectorCoord,
		IsTransportUnreachable: base.IsTransportUnreachable,
		Reachable:              base.Reachable,
		IsPob:                  base.IsPob,
	}

	if include_market_goods {
		base.MarketGoodsPerNick = make(map[configs_export.CommodityKey]*configs_export.MarketGood)
		for key, good := range FilterMarketGoodCategory(filter_market_good_category, base.MarketGoodsPerNick) {
			item.MarketGoodsPerNick[key] = NewMarketGood(good)
		}
	}

	return item
}
