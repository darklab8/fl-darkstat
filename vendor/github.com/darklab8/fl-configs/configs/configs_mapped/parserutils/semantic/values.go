/*
ORM mapper for Freelancer ini reader. Easy mapping values to change.
*/
package semantic

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/inireader"
)

// ORM values

type ValueType int64

const (
	TypeComment ValueType = iota
	TypeVisible
)

type Value struct {
	section    *inireader.Section
	key        string
	optional   bool
	value_type ValueType
	order      int
	index      int
}

func NewValue(
	section *inireader.Section,
	key string,
) *Value {
	return &Value{
		section:    section,
		key:        key,
		value_type: TypeVisible,
	}
}

func (v Value) isComment() bool {
	return v.value_type == TypeComment
}

type ValueOption func(i *Value)

func Order(order int) ValueOption {
	return func(i *Value) {
		i.order = order
	}
}

func Index(index int) ValueOption {
	return func(i *Value) {
		i.index = index
	}
}

func Optional() ValueOption {
	return func(i *Value) {
		i.optional = true
	}
}
