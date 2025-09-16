package darkgrpc

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/go-utils/utils/ptr"
)

func NewInt64(value *int) *int64 {
	if value == nil {
		return nil
	}
	q := int64(*value)
	return &q
}

func NewShipClass(ship_class *cfg.ShipClass) *int64 {
	if ship_class == nil {
		return nil
	}

	return ptr.Ptr(int64(*ship_class))
}

func NewPos(pos *cfg.Vector) *pb.Pos {
	if pos == nil {
		return nil
	}
	return &pb.Pos{
		X: pos.X,
		Y: pos.Y,
		Z: pos.Z,
	}
}
