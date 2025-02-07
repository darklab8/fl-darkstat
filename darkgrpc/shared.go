package darkgrpc

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
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
