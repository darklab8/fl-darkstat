package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func (s *Server) GetFactions(_ context.Context, in *pb.GetFactionsInput) (*pb.GetFactionsReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var input []configs_export.Faction
	if in.FilterToUseful {
		input = configs_export.FilterToUsefulFactions(s.app_data.Configs.Factions)
	} else {
		input = s.app_data.Configs.Factions
	}

	var items []*pb.Faction
	for _, item := range input {
		result := &pb.Faction{
			Name:              item.Name,
			ShortName:         item.ShortName,
			Nickname:          item.Nickname,
			ObjectDestruction: item.ObjectDestruction,
			MissionSuccess:    item.MissionSuccess,
			MissionFailure:    item.MissionFailure,
			MissionAbort:      item.MissionAbort,
			InfonameId:        int64(item.InfonameID),
			InfocardId:        int64(item.InfocardID),
		}
		if in.IncludeBribes {
			for _, bribe := range item.Bribes {
				result.Bribes = append(result.Bribes, &pb.Bribe{
					BaseNickname: bribe.BaseNickname,
					Chance:       bribe.Chance,
					BaseInfo:     NewBaseInfo(bribe.BaseInfo),
				})
			}
		}
		if in.IncludeReputations {
			for _, rep := range item.Reputations {
				result.Reputations = append(result.Reputations, &pb.Reputation{
					Name:     rep.Name,
					Rep:      rep.Rep,
					Empathy:  rep.Empathy,
					Nickname: rep.Nickname,
				})
			}

		}
		items = append(items, result)
	}
	return &pb.GetFactionsReply{Items: items}, nil
}
