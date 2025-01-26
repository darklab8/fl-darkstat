package settings

import (
	"fmt"

	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkcoreEnvVars struct {
	utils_settings.UtilsEnvs
	Password        string
	Secret          string
	ExtraCookieHost string
}

var Env DarkcoreEnvVars

func GetEnvs(envs *enverant.Enverant) DarkcoreEnvVars {
	Env = DarkcoreEnvVars{
		UtilsEnvs: utils_settings.GetEnvs(envs),
		Password:  envs.GetStrOr("DARKCORE_PASSWORD", ""),
		Secret:    envs.GetStrOr("DARKCORE_SECRET", "passphrasewhichneedstobe32bytes!"),
	}
	return Env
}

func init() {
	env := enverant.NewEnverant()
	Env = GetEnvs(env)
	fmt.Sprintln("conf=", Env)
}
