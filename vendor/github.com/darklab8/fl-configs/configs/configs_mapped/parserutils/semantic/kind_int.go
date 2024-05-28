package semantic

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-configs/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

type Int struct {
	*Value
}

func NewInt(section *inireader.Section, key string, opts ...ValueOption) *Int {
	v := NewValue(section, key)
	for _, opt := range opts {
		opt(v)
	}
	s := &Int{Value: v}

	return s
}

func (s *Int) Get() int {
	if s.optional && len(s.section.ParamMap[s.key]) == 0 {
		return 0
	}
	return int(s.section.ParamMap[s.key][s.index].Values[s.order].(inireader.ValueNumber).Value)
}

func (s *Int) GetValue() (int, bool) {
	var value int
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

func (s *Int) Set(value int) {
	if s.isComment() {
		s.Delete()
	}

	processed_value := inireader.UniParseInt(value)
	if len(s.section.ParamMap[s.key]) == 0 {
		s.section.AddParamToStart(s.key, (&inireader.Param{IsComment: s.isComment()}).AddValue(processed_value))
	}
	// implement SetValue in Section
	s.section.ParamMap[s.key][0].First = processed_value
	s.section.ParamMap[s.key][0].Values[0] = processed_value
}

func (s *Int) Delete() {
	delete(s.section.ParamMap, s.key)
	for index, param := range s.section.Params {
		if param.Key == s.key {
			s.section.Params = append(s.section.Params[:index], s.section.Params[index+1:]...)
		}
	}
}
