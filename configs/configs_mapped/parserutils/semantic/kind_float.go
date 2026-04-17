package semantic

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/typelog"
)

type Precision int

type Float struct {
	*Value
	precision     Precision
	default_value float64
}

type FloatOption func(s *Float)

func WithDefaultF(default_value float64) FloatOption {
	return func(s *Float) { s.default_value = default_value }
}

func OptsF(opts ...ValueOption) FloatOption {
	return func(s *Float) {
		for _, opt := range opts {
			opt(s.Value)
		}
	}
}

func NewFloat(section *inireader.Section, key cfg.ParamKey, precision Precision, opts ...FloatOption) *Float {
	v := NewValue(section, key)

	s := &Float{
		Value:     v,
		precision: precision,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Float) get() float64 {
	var section *inireader.Section = s.section
	if _, ok := s.section.ParamMap[s.key]; !ok {
		if inherit_value, ok := s.section.ParamMap[InheritKey]; ok {
			inherit_nick := inherit_value[0].First.AsString()
			if found_section, ok := s.section.INIFile.SectionMapByNick[inherit_nick]; ok {
				section = found_section
			}
		}
	}

	if s.optional && len(section.ParamMap[s.key]) == 0 {
		return 0
	}

	switch v := section.ParamMap[s.key][s.index].Values[s.order].(type) {
	case inireader.ValueNumber:
		return v.Value
	case inireader.ValueString:
		if strings.ToLower(v.AsString()) == "infinity" {
			return 2_000_000_000 // int32 max should be enough
		} else {
			panic("not recognized string instead of number")
		}
	default:
		logus.Log.Panic("not supported uni value type")
	}
	panic("not possible")
}

func (s *Float) Get() float64 {
	defer handleGetCrashReporting(s.Value)
	return s.get()
}

func (s *Float) GetValue() (float64, bool) {
	var value float64 = s.default_value
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

func (s *Float) Set(value float64) {
	if s.isComment() {
		s.Delete()
	}

	processed_value := inireader.UniParseFloat(value, int(s.precision))
	if len(s.section.ParamMap[s.key]) == 0 {
		s.section.AddParamToStart(s.key, (&inireader.Param{IsParamAsComment: s.isComment()}).AddValue(processed_value))
	}
	// implement SetValue in Section
	s.section.ParamMap[s.key][0].First = processed_value
	s.section.ParamMap[s.key][0].Values[0] = processed_value
}

func (s *Float) Delete() {
	delete(s.section.ParamMap, s.key)
	for index, param := range s.section.Params {
		if param.Key == s.key {
			s.section.Params = append(s.section.Params[:index], s.section.Params[index+1:]...)
		}
	}
}
