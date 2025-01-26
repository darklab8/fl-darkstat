package cfg

type Vector struct {
	X float64
	Y float64
	Z float64
}

type TractorID string

type FactionNick string

type Milliseconds = float64
type Seconds = float64

type MillisecondsI = int
type SecondsI = int

type BaseUniNick string

func (b BaseUniNick) ToStr() string { return string(b) }

type ShipClass int64

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
