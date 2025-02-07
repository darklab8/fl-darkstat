package darkgrpcsrv

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkgrpc"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/go-utils/utils/ptr"
)

// SayHello implements helloworld.GreeterServer
func (s *Server) GetBases(_ context.Context, in *pb.Empty) (*pb.GetBasesReply, error) {
	var bases []*pb.Base
	for _, base := range s.app_data.Configs.Bases {
		item := &pb.Base{
			Name:                   base.Name,
			Archetypes:             base.Archetypes,
			Nickname:               string(base.Nickname),
			FactionName:            base.FactionName,
			System:                 base.System,
			SystemNickname:         base.SystemNickname,
			Region:                 base.Region,
			StridName:              int32(base.StridName),
			InfocardID:             int32(base.InfocardID),
			File:                   base.File.ToString(),
			BGCSBaseRunBy:          base.BGCS_base_run_by,
			Pos:                    NewPos(base.Pos),
			SectorCoord:            base.SectorCoord,
			IsTransportUnreachable: base.IsTransportUnreachable,
			Reachable:              base.Reachable,
			IsPob:                  base.IsPob,
		}

		bases = append(bases, item)
	}
	return &pb.GetBasesReply{Items: bases}, nil
}

func (s *Server) GetBasesMarketGoods(_ context.Context, in *pb.GetMarketGoodsInput) (*pb.GetMarketGoodsReply, error) {
	bases_by_nick := make(map[string]*configs_export.Base)
	for _, base := range s.app_data.Configs.Bases {
		bases_by_nick[string(base.Nickname)] = base
	}
	var answers []*pb.MarketGoodAnswer
	for _, nickname := range in.Nicknames {
		answer := &pb.MarketGoodAnswer{Nickname: string(nickname)}
		if base, ok := bases_by_nick[nickname]; ok {
			for _, good := range base.MarketGoodsPerNick {
				answer.MarketGoods = append(answer.MarketGoods, NewMarketGood(good))
			}
		} else {
			answer.Error = ptr.Ptr("not existing base")
		}
		answers = append(answers, answer)
	}
	return &pb.GetMarketGoodsReply{Answers: answers}, nil
}
