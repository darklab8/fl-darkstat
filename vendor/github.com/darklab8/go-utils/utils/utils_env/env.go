package utils_env

import "os"

type EnvBool string

func (e EnvBool) ToBool() bool {
	return e == "true"
}

type EnvConfig struct{}

func NewEnvConfig() EnvConfig {
	return EnvConfig{}
}

func (e EnvConfig) GetEnvWithDefault(key string, default_ string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		return default_
	}
}

func (e EnvConfig) GetEnvBool(key string) EnvBool {
	return EnvBool(os.Getenv(key))
}

func (e EnvConfig) GetEnv(key string) string {
	return os.Getenv(key)
}
