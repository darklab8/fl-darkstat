package relaysettings

import (
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkrelayEnvVars struct {
	utils_settings.UtilsEnvs
	configs_settings.ConfEnvVars
}

var Env DarkrelayEnvVars

func init() {
	Env = DarkrelayEnvVars{
		UtilsEnvs:   utils_settings.GetEnvs(),
		ConfEnvVars: configs_settings.GetEnvs(),
	}
}
