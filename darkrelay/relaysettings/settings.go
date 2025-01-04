package relaysettings

import (
	"fmt"

	"github.com/darklab8/fl-configs/configs/configs_settings"
	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkrelayEnvVars struct {
	utils_settings.UtilsEnvs
	configs_settings.ConfEnvVars
	AppVersion string
}

var Env DarkrelayEnvVars

func init() {
	env := enverant.NewEnverant()
	Env = DarkrelayEnvVars{
		UtilsEnvs:   utils_settings.GetEnvs(env),
		ConfEnvVars: configs_settings.GetEnvs(env),
		AppVersion:  env.GetStrOr("BUILD_VERSION", "v0.0.0-dev"),
	}
	fmt.Sprintln("conf=", Env)
}