package utils_settings

import "github.com/darklab8/go-utils/utils/enverant"

type UtilsEnvs struct {
	IsDevEnv             bool
	AreTestsRegenerating bool
}

var Envs UtilsEnvs

func init() {
	envs := enverant.NewEnverant()
	GetEnvs(envs)
}

func GetEnvs(envs *enverant.Enverant) UtilsEnvs {
	Envs = UtilsEnvs{
		IsDevEnv:             envs.GetBoolOr("DEV_ENV", false),
		AreTestsRegenerating: envs.GetBoolOr("DARK_TEST_REGENERATE", false),
	}
	return Envs
}
