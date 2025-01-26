package semantic

import (
	"path/filepath"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/utils_types"
)

// Linux friendly filepath, that can be returned to Windows way from linux
type Path struct {
	*Value
	remove_spaces bool
	lowercase     bool
}

type PathOption func(s *Path)

func WithoutSpacesP() PathOption {
	return func(s *Path) { s.remove_spaces = true }
}

func WithLowercaseP() PathOption {
	return func(s *Path) { s.lowercase = true }
}

func OptsP(opts ...ValueOption) PathOption {
	return func(s *Path) {
		for _, opt := range opts {
			opt(s.Value)
		}
	}
}

func NewPath(section *inireader.Section, key cfg.ParamKey, opts ...PathOption) *Path {
	v := NewValue(section, key)
	s := &Path{Value: v}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Path) FileName() utils_types.FilePath {
	if s.optional && len(s.section.ParamMap[s.key]) == 0 {
		return ""
	}
	value := s.section.ParamMap[s.key][s.index].Values[s.order].AsString()
	value = strings.ReplaceAll(value, "\\", PATH_SEPARATOR)
	value = filepath.Base(value)
	if s.remove_spaces {
		value = strings.ReplaceAll(value, " ", "")
	}
	if s.lowercase {
		value = strings.ToLower(value)
	}
	return utils_types.FilePath(value)
}

func (s *Path) get() utils_types.FilePath {
	if s.optional && len(s.section.ParamMap[s.key]) == 0 {
		return ""
	}
	value := s.section.ParamMap[s.key][s.index].Values[s.order].AsString()
	value = strings.ReplaceAll(value, "\\", PATH_SEPARATOR)
	if s.remove_spaces {
		value = strings.ReplaceAll(value, " ", "")
	}
	if s.lowercase {
		value = strings.ToLower(value)
	}
	return utils_types.FilePath(value)
}

func (s *Path) Get() utils_types.FilePath {
	defer handleGetCrashReporting(s.Value)
	return s.get()
}

func (s *Path) GetValue() (utils_types.FilePath, bool) {
	var value utils_types.FilePath
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

func (s *Path) Set(value utils_types.FilePath) {
	if s.isComment() {
		s.Delete()
	}

	processed_value := inireader.UniParseStr(string(value))
	if len(s.section.ParamMap[s.key]) == 0 {
		s.section.AddParamToStart(s.key, (&inireader.Param{IsParamAsComment: s.isComment()}).AddValue(processed_value))
	}
	// implement SetValue in Section
	s.section.ParamMap[s.key][0].First = processed_value
	s.section.ParamMap[s.key][0].Values[0] = processed_value
}

func (s *Path) Delete() {
	delete(s.section.ParamMap, s.key)
	for index, param := range s.section.Params {
		if param.Key == s.key {
			s.section.Params = append(s.section.Params[:index], s.section.Params[index+1:]...)
		}
	}
}
