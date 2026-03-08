package cfg

import (
	"strconv"

	"github.com/darklab8/go-utils/utils/ptr"
)

type CtxKey string

type Vector struct {
	X float64
	Y float64
	Z float64
}

type TractorID string

type FactionNick string

type Milliseconds = float64
type Seconds = float64

type MillisecondsI = int32
type SecondsI = int

type BaseUniNick string

func (b BaseUniNick) ToStr() string { return string(b) }

// ShipClass is used to show item volume for specific ship class, used for commoditiies
type ShipClass int64

func ShipClassToKey(s *ShipClass) string {
	if s == nil {
		return "nil"
	}

	return strconv.Itoa(int(*s))
}

func (s ShipClass) ToStr() string {
	switch s {
	case 10:
		return "liner"
	case 14:
		return "miner"
	default:
		return ""
	}
}

// Gob friendly
type ErrP string
type Err = *ErrP

func NewErr(msg string) Err {
	return Err(ptr.Ptr(msg))
}

func (r *ErrP) Error() string {
	return string(*r)
}

type WithDiscoFreighterPaths bool

type Reachability struct {
	IsTransportReachable bool `json:"is_transport_reachable"` // Check if base is NOT reachable from manhattan by Transport through Graph method (at Discovery base has to have Transport dockable spheres)
	IsFreighterReachable bool `json:"is_freighter_reachable"` // is base IS Rechable by freighter from Manhattan
	IsFrigateReachable   bool `json:"is_frigate_reachhable"`  // is base IS Rechable by frigate from Manhattan
}
