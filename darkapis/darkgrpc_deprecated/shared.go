package darkgrpc_deprecated

import (
	"strconv"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated/statproto_deprecated"
	"github.com/darklab8/go-utils/utils/ptr"
)

func NewInt64S(value int) pb.NumString {
	q := strconv.Itoa(value)
	return &q
}

func NewInt64(value *int) pb.NumString {
	if value == nil {
		return nil
	}
	q := strconv.Itoa(*value)
	return &q
}

func NewShipClass(ship_class *cfg.ShipClass) *string {
	if ship_class == nil {
		return nil
	}

	return ptr.Ptr(strconv.Itoa(int(*ship_class)))
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
