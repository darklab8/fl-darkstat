package enverant

/*
Manager for getting values from Environment variables
*/

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

type Enverant struct {
	file_envs        map[string]interface{}
	validate_missing bool
}

type EnverantOption func(m *Enverant)

func NewEnverant(opts ...EnverantOption) *Enverant {
	e := &Enverant{
		file_envs: map[string]interface{}{},
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func WithEnvFile(path string) EnverantOption {
	return func(m *Enverant) {
		m.file_envs = ReadJson(path)
	}
}

func WithValidate(validate bool) EnverantOption {
	return func(m *Enverant) {
		m.validate_missing = validate
	}
}

func EnrichStr(value string) string {
	// unhardcode later
	if strings.Contains(value, "${env:HOME}") {
		value = strings.ReplaceAll(value, "${env:HOME}", os.Getenv("HOME"))
	}
	return value
}

type ValueParams struct {
	default_ any
}
type ValueOption func(m *ValueParams)

func OrStr(default_ string) ValueOption {
	return func(m *ValueParams) {
		m.default_ = default_
	}
}

func OrInt(default_ int) ValueOption {
	return func(m *ValueParams) {
		m.default_ = default_
	}
}

func OrBool(default_ bool) ValueOption {
	return func(m *ValueParams) {
		m.default_ = default_
	}
}

func (e *Enverant) GetStr(key string, opts ...ValueOption) string {
	if value, ok := e.GetString(key, opts...); ok {
		return value
	}
	return ""
}

func (e *Enverant) GetPtrStr(key string, opts ...ValueOption) *string {
	if value, ok := e.GetString(key, opts...); ok {
		return ptr.Ptr(value)
	}
	return nil
}

func (e *Enverant) GetString(key string, opts ...ValueOption) (string, bool) {
	params := &ValueParams{}
	for _, opt := range opts {
		opt(params)
	}

	if value, ok := os.LookupEnv(key); ok {
		return value, true
	}

	if value, ok := e.file_envs[key]; ok {
		return EnrichStr(value.(string)), true
	}

	if params.default_ != nil {
		return params.default_.(string), true
	}

	if e.validate_missing {
		utils_logus.Log.Panic("enverant value is not defined", typelog.String("key", key))
	}

	return "", false
}

func (e *Enverant) GetBool(key string, opts ...ValueOption) bool {
	if value, ok := e.GetBoolean(key, opts...); ok {
		return value
	}
	return false
}

func (e *Enverant) GetPtrBool(key string, opts ...ValueOption) *bool {
	if value, ok := e.GetBoolean(key, opts...); ok {
		return ptr.Ptr(value)
	}
	return nil
}

func (e *Enverant) GetBoolean(key string, opts ...ValueOption) (bool, bool) {
	params := &ValueParams{}
	for _, opt := range opts {
		opt(params)
	}

	if value, ok := os.LookupEnv(key); ok {
		return value == "true", true
	}

	if value, ok := e.file_envs[key]; ok {
		switch v := value.(type) {
		case bool:
			return v, true
		case string:
			return v == "true", true
		default:
			panic(fmt.Sprintln("unrecognized type for value", v, " in GetBoolOr"))
		}
	}

	if params.default_ != nil {
		return params.default_.(bool), true
	}

	if e.validate_missing {
		utils_logus.Log.Panic("enverant value is not defined", typelog.String("key", key))
	}

	return false, false
}

func (e *Enverant) GetInt(key string, opts ...ValueOption) int {
	if value, ok := e.GetInteger(key, opts...); ok {
		return value
	}
	return 0
}

func (e *Enverant) GetPtrInt(key string, opts ...ValueOption) *int {
	if value, ok := e.GetInteger(key, opts...); ok {
		return ptr.Ptr(value)
	}
	return nil
}

func (e *Enverant) GetInteger(key string, opts ...ValueOption) (int, bool) {
	params := &ValueParams{}
	for _, opt := range opts {
		opt(params)
	}

	if value, ok := os.LookupEnv(key); ok {
		int_value, err := strconv.Atoi(value)
		utils_logus.Log.CheckPanic(err, "expected to be int", typelog.String("key", key))
		return int_value, true
	}

	if value, ok := e.file_envs[key]; ok {
		switch v := value.(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		case string:
			int_value, err := strconv.Atoi(v)
			utils_logus.Log.CheckPanic(err, "expected to be int", typelog.String("key", key))
			return int_value, true
		default:
			panic(fmt.Sprintln("unrecognized type for value", v, " in GetIntOr"))
		}
	}

	if params.default_ != nil {
		return params.default_.(int), true
	}

	if e.validate_missing {
		utils_logus.Log.Panic("enverant value is not defined", typelog.String("key", key))
	}

	return 0, false
}
