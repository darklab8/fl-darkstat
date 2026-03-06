package semantic

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/typelog"
)

type Bool struct {
	*Value
	bool_type     BoolType
	default_value bool
}

type BoolType int64

const (
	IntBool BoolType = iota
	StrBool
	FLSRActiveBool
)

type BoolOption func(s *Bool)

func WithDefaultB(default_value bool) BoolOption {
	return func(s *Bool) { s.default_value = default_value }
}

func OptsB(opts ...ValueOption) BoolOption {
	return func(s *Bool) {
		for _, opt := range opts {
			opt(s.Value)
		}
	}
}

func NewBool(section *inireader.Section, key cfg.ParamKey, bool_type BoolType, opts ...BoolOption) *Bool {
	v := NewValue(section, key)

	s := &Bool{
		Value:     v,
		bool_type: bool_type,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Bool) get() bool {
	if s.optional && len(s.section.ParamMap[s.key]) == 0 {
		return s.default_value
	}
	switch s.bool_type {
	case IntBool:
		return int(s.section.ParamMap[s.key][s.index].Values[s.order].(inireader.ValueNumber).Value) == 1
	case StrBool:
		return strings.Contains(strings.ToLower(s.section.ParamMap[s.key][s.index].Values[s.order].AsString()), "true")
	case FLSRActiveBool:
		return strings.ToLower(s.section.ParamMap[s.key][s.index].Values[s.order].AsString()) == "active"
	}

	panic("not expected bool type")
}

func (s *Bool) Get() bool {
	defer handleGetCrashReporting(s.Value)
	return s.get()
}

func (s *Bool) GetValue() (bool, bool) {
	var value bool = s.default_value
	var ok bool = true
	func() {
		defer func() {
			if r := recover(); r != nil {
				logus.Log.Debug("Recovered from int GetValue Error:\n", typelog.Any("recover", r))
				ok = false
			}
		}()
		value = s.get()
	}()

	return value, ok
}

func (s *Bool) Set(value bool) {
	var processed_value inireader.UniValue
	if s.isComment() {
		s.Delete()
	}

	switch s.bool_type {
	case IntBool:
		var int_bool int
		if value {
			int_bool = 1
		}
		processed_value = inireader.UniParseInt(int_bool)
	case StrBool:
		if value {
			processed_value = inireader.UniParseStr("true")
		} else {
			processed_value = inireader.UniParseStr("false")
		}
	}

	if len(s.section.ParamMap[s.key]) == 0 {
		s.section.AddParamToStart(s.key, (&inireader.Param{IsParamAsComment: s.isComment()}).AddValue(processed_value))
	}
	// implement SetValue in Section
	s.section.ParamMap[s.key][0].First = processed_value
	s.section.ParamMap[s.key][0].Values[0] = processed_value
}

func (s *Bool) Delete() {
	delete(s.section.ParamMap, s.key)
	for index, param := range s.section.Params {
		if param.Key == s.key {
			s.section.Params = append(s.section.Params[:index], s.section.Params[index+1:]...)
		}
	}
}
