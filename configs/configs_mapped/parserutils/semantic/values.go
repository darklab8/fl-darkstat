/*
ORM mapper for Freelancer ini reader. Easy mapping values to change.
*/
package semantic

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

// ORM values

type ValueType int64

const (
	TypeComment ValueType = iota
	TypeVisible
)

type Value struct {
	section    *inireader.Section
	key        cfg.ParamKey
	optional   bool
	value_type ValueType
	order      int
	index      int
}

func NewValue(
	section *inireader.Section,
	key cfg.ParamKey,
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

func Comment() ValueOption {
	return func(i *Value) {
		i.value_type = TypeComment
	}
}

func quickJson(value any) string {
	result, err := json.Marshal(value)
	if err != nil {
		return err.Error()
	}
	return string(result)
}

func handleGetCrashReporting(value *Value) {
	if r := recover(); r != nil {
		if value == nil {
			logus.Log.Panic("value is not defined. not possible. ;)")
			return
		} else {
			var section strings.Builder
			section.WriteString(string(value.section.Type))
			for _, param := range value.section.Params {
				section.WriteString(fmt.Sprintf("\"%s\"", param.ToString(inireader.WithComments(true))))
			}
			logus.Log.Error("unable to Get() from semantic.",
				typelog.Any("value", quickJson(value)),
				typelog.Any("key", value.key),
				typelog.String("section", section.String()),
				typelog.Any("r", r),
			)
		}
		panic(r)
	}
}
