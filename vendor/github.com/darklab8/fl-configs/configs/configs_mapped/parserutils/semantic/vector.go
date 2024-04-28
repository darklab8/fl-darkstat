package semantic

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/conftypes"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-configs/configs/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

type Vect struct {
	Model
	X *Float
	Y *Float
	Z *Float
}

func NewVector(section *inireader.Section, key string, precision Precision, opts ...ValueOption) *Vect {
	v := &Vect{
		X: NewFloat(section, key, precision, Order(0)),
		Y: NewFloat(section, key, precision, Order(1)),
		Z: NewFloat(section, key, precision, Order(2)),
	}
	v.Map(section)
	return v
}

func (s *Vect) Get() conftypes.Vector {
	return conftypes.Vector{
		X: s.X.Get(),
		Y: s.Y.Get(),
		Z: s.Z.Get(),
	}
}

func (s *Vect) GetValue() (conftypes.Vector, bool) {
	var value conftypes.Vector
	var ok bool = true
	func() {
		defer func() {
			if r := recover(); r != nil {
				logus.Log.Debug("Recovered from int GetValue Error:\n", typelog.Any("recover", r))
				ok = false
			}
		}()
		value = s.Get()
	}()

	return value, ok
}
