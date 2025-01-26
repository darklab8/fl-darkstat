package semantic

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

type String struct {
	*Value
	remove_spaces bool
	lowercase     bool
}

type StringOption func(s *String)

func WithoutSpacesS() StringOption {
	return func(s *String) { s.remove_spaces = true }
}

func WithLowercaseS() StringOption {
	return func(s *String) { s.lowercase = true }
}

func OptsS(opts ...ValueOption) StringOption {
	return func(s *String) {
		for _, opt := range opts {
			opt(s.Value)
		}
	}
}

func NewString(section *inireader.Section, key cfg.ParamKey, opts ...StringOption) *String {
	v := NewValue(section, key)
	s := &String{Value: v}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *String) get() string {
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
		return ""
	}
	value := section.ParamMap[s.key][s.index].Values[s.order].AsString()
	if s.remove_spaces {
		value = strings.ReplaceAll(value, " ", "")
	}
	if s.lowercase {
		value = strings.ToLower(value)
	}
	return value

}

func (s *String) Get() string {

	defer handleGetCrashReporting(s.Value)
	return s.get()
}

func (s *String) GetValue() (string, bool) {
	var value string
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

func (s *String) Set(value string) {
	if s.isComment() {
		s.Delete()
	}

	processed_value := inireader.UniParseStr(value)
	if len(s.section.ParamMap[s.key]) == 0 {
		s.section.AddParamToStart(s.key, (&inireader.Param{IsParamAsComment: s.isComment()}).AddValue(processed_value))
	}
	// implement SetValue in Section
	s.section.ParamMap[s.key][0].First = processed_value
	s.section.ParamMap[s.key][0].Values[0] = processed_value
}

func (s *String) Delete() {
	delete(s.section.ParamMap, s.key)
	for index, param := range s.section.Params {
		if param.Key == s.key {
			s.section.Params = append(s.section.Params[:index], s.section.Params[index+1:]...)
		}
	}
}
