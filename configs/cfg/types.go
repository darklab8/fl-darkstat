package cfg

import (
	"strconv"

	"github.com/darklab8/go-utils/utils/ptr"
)

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
