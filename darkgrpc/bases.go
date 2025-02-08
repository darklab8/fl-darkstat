package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetBases(_ context.Context, in *pb.GetBasesInput) (*pb.GetBasesReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

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
			StridName:              int64(base.StridName),
			InfocardID:             int64(base.InfocardID),
			File:                   base.File.ToString(),
			BGCSBaseRunBy:          base.BGCS_base_run_by,
			Pos:                    NewPos(base.Pos),
			SectorCoord:            base.SectorCoord,
			IsTransportUnreachable: base.IsTransportUnreachable,
			Reachable:              base.Reachable,
			IsPob:                  base.IsPob,
		}

		if in.IncludeMarketGoods {
			base.MarketGoodsPerNick = make(map[configs_export.CommodityKey]*configs_export.MarketGood)
			for key, good := range base.MarketGoodsPerNick {
				item.MarketGoodsPerNick[string(key)] = NewMarketGood(good)
			}
		}

		bases = append(bases, item)
	}
	return &pb.GetBasesReply{Items: bases}, nil
}
