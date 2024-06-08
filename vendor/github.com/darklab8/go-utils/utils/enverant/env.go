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
	"github.com/darklab8/go-utils/utils/utils_logus"
)

type Enverant struct {
	file_envs map[string]interface{}
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

func EnrichStr(value string) string {
	// unhardcode later
	if strings.Contains(value, "${env:HOME}") {
		value = strings.ReplaceAll(value, "${env:HOME}", os.Getenv("HOME"))
	}
	return value
}

func (e *Enverant) GetStr(key string) string {
	return e.GetStrOr(key, "")
}

func (e *Enverant) GetStrOr(key string, default_ string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if value, ok := e.file_envs[key]; ok {
		return EnrichStr(value.(string))
	}

	return default_
}

func (e *Enverant) GetBool(key string) bool {
	return e.GetBoolOr(key, false)
}

func (e *Enverant) GetBoolOr(key string, default_ bool) bool {

	if value, ok := os.LookupEnv(key); ok {
		return value == "true"
	}

	if value, ok := e.file_envs[key]; ok {
		switch v := value.(type) {
		case bool:
			return v
		case string:
			return v == "true"
		default:
			panic(fmt.Sprintln("unrecognized type for value", v, " in GetBoolOr"))
		}
	}

	return default_
}

func (e *Enverant) GetInt(key string) int {
	return e.GetIntOr(key, 0)
}

func (e *Enverant) GetIntOr(key string, default_ int) int {

	if value, ok := os.LookupEnv(key); ok {
		int_value, err := strconv.Atoi(value)
		utils_logus.Log.CheckPanic(err, "expected to be int", typelog.String("key", key))
		return int_value
	}

	if value, ok := e.file_envs[key]; ok {
		switch v := value.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			int_value, err := strconv.Atoi(v)
			utils_logus.Log.CheckPanic(err, "expected to be int", typelog.String("key", key))
			return int_value
		default:
			panic(fmt.Sprintln("unrecognized type for value", v, " in GetIntOr"))
		}
	}

	return default_
}
