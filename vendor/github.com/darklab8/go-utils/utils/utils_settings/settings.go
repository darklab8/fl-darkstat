package utils_settings

import "github.com/darklab8/go-utils/utils/utils_env"

type UtilsEnvs struct {
	IsDevEnv             utils_env.EnvBool
	AreTestsRegenerating utils_env.EnvBool
}

var Envs UtilsEnvs

func init() {
	envs := utils_env.NewEnvConfig()
	Envs = UtilsEnvs{
		IsDevEnv:             envs.GetEnvBool("DEV_ENV"),
		AreTestsRegenerating: envs.GetEnvBool("DARK_TEST_REGENERATE"),
	}
}
