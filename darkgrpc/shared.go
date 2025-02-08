package darkgrpc

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/go-utils/utils/ptr"
)

func IntTo32(value *int) *int32 {
	if value == nil {
		return nil
	}
	q := int32(*value)
	return &q
}

func NewMarketGood(good *configs_export.MarketGood) *pb.MarketGood {
	return &pb.MarketGood{
		Nickname: good.Nickname,

		ShipNickname:         good.ShipNickname,
		Name:                 good.Name,
		PriceBase:            int32(good.PriceBase),
		HpType:               good.HpType,
		Category:             good.Category,
		LevelRequired:        int32(good.LevelRequired),
		RepRequired:          good.RepRequired,
		PriceBaseBuysFor:     IntTo32(good.PriceBaseBuysFor),
		PriceBaseSellsFor:    int32(good.PriceBaseSellsFor),
		Volume:               good.Volume,
		ShipClass:            int64(good.ShipClass),
		BaseSells:            good.BaseSells,
		IsServerSideOverride: good.IsServerSideOverride,

		NotBuyable:             good.NotBuyable,
		IsTransportUnreachable: good.IsTransportUnreachable,
		BaseNickname:           good.BaseName,
		BaseName:               good.BaseName,
		SystemName:             good.SystemName,
		Region:                 good.Region,
		FactionName:            good.FactionName,
		BasePos:                NewPos(good.BasePos),
		SectorCoord:            good.SectorCoord,
	}
}

func NewPos(pos cfg.Vector) *pb.Pos {
	return &pb.Pos{
		X: pos.X,
		Y: pos.Y,
		Z: pos.Z,
	}
}

type Nicknamable interface {
	GetNickname() string
}

type Marketable interface {
	Nicknamable
	GetBases() map[cfg.BaseUniNick]*configs_export.MarketGood
}

func GetMarketGoods[T Marketable](items []T, in *pb.GetMarketGoodsInput) (*pb.GetMarketGoodsReply, error) {
	items_by_nick := make(map[string]T)
	for _, item := range items {
		items_by_nick[string(item.GetNickname())] = item
	}
	var answers []*pb.MarketGoodAnswer
	for _, nickname := range in.Nicknames {
		answer := &pb.MarketGoodAnswer{Nickname: string(nickname)}
		if base, ok := items_by_nick[nickname]; ok {
			for _, good := range base.GetBases() {
				answer.MarketGoods = append(answer.MarketGoods, NewMarketGood(good))
			}
		} else {
			answer.Error = ptr.Ptr("not existing base")
		}
		answers = append(answers, answer)
	}
	return &pb.GetMarketGoodsReply{Answers: answers}, nil
}

func NewBases(Bases map[cfg.BaseUniNick]*configs_export.MarketGood) map[string]*pb.MarketGood {
	result := make(map[string]*pb.MarketGood)
	for key, item := range Bases {
		result[string(key)] = NewMarketGood(item)
	}
	return result
}

type TechCompatable interface {
	Nicknamable
	GetDiscoveryTechCompat() *configs_export.DiscoveryTechCompat
}

func GetTechCompat[T TechCompatable](items []T, in *pb.GetTechCompatInput) (*pb.GetTechCompatReply, error) {
	items_by_nick := make(map[string]T)
	for _, item := range items {
		items_by_nick[string(item.GetNickname())] = item
	}
	var answers []*pb.TechCompatAnswer

	for _, nickname := range in.Nicknames {

		answer := &pb.TechCompatAnswer{Nickname: string(nickname)}
		if item, ok := items_by_nick[nickname]; ok {
			answer.TechCompat = NewTechCompat(item.GetDiscoveryTechCompat())
		} else {
			answer.Error = ptr.Ptr("not existing nickname")
		}
		answers = append(answers, answer)

	}
	return &pb.GetTechCompatReply{Answers: answers}, nil
}

func NewTechCompat(tech_compat *configs_export.DiscoveryTechCompat) *pb.DiscoveryTechCompat {
	if tech_compat == nil {
		return nil
	}

	answer := &pb.DiscoveryTechCompat{
		TechcompatByID: make(map[string]float64),
		TechCell:       tech_compat.TechCell,
	}
	for key, value := range tech_compat.TechcompatByID {
		answer.TechcompatByID[string(key)] = value
	}
	return answer
}
