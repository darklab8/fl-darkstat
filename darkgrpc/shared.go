package darkgrpc

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func NewInt64(value *int) *int64 {
	if value == nil {
		return nil
	}
	q := int64(*value)
	return &q
}

func NewMarketGood(good *configs_export.MarketGood) *pb.MarketGood {
	return &pb.MarketGood{
		Nickname: good.Nickname,

		ShipNickname:         good.ShipNickname,
		Name:                 good.Name,
		PriceBase:            int64(good.PriceBase),
		HpType:               good.HpType,
		Category:             good.Category,
		LevelRequired:        int64(good.LevelRequired),
		RepRequired:          good.RepRequired,
		PriceBaseBuysFor:     NewInt64(good.PriceBaseBuysFor),
		PriceBaseSellsFor:    int64(good.PriceBaseSellsFor),
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
