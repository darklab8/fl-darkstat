package enverant

/*
Manager for getting values from Environment variables
*/

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/darklab8/go-utils/utils/ptr"
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

func (m *Enverant) GetValidating() *Enverant {
	var clone *Enverant = &Enverant{}
	*clone = *m
	clone.validate_missing = true
	return clone
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

func (e *Enverant) GetStrOr(key string, default_ string, opts ...ValueOption) string {
	value, _ := e.GetString(key, append([]ValueOption{OrStr(default_)}, opts...)...)
	return value
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
		return EnrichStr(value), true
	}

	if params.default_ != nil {
		return params.default_.(string), true
	}

	if e.validate_missing {
		panic(fmt.Sprintln("enverant value is not defined, key=", key))
	}

	return "", false
}

func (e *Enverant) GetBoolOr(key string, default_ bool, opts ...ValueOption) bool {
	value, _ := e.GetBoolean(key, append([]ValueOption{OrBool(default_)}, opts...)...)
	return value
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

	if params.default_ != nil {
		return params.default_.(bool), true
	}

	if e.validate_missing {
		panic(fmt.Sprintln("enverant value is not defined, key=", key))
	}

	return false, false
}

func (e *Enverant) GetIntOr(key string, default_ int, opts ...ValueOption) int {
	value, _ := e.GetInteger(key, append([]ValueOption{OrInt(default_)}, opts...)...)
	return value
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
		if err != nil {
			panic(fmt.Sprintln(err, "expected to be int, key=", key))
		}
		return int_value, true
	}

	if params.default_ != nil {
		return params.default_.(int), true
	}

	if e.validate_missing {
		panic(fmt.Sprintln("enverant value is not defined, key=", key))
	}

	return 0, false
}
