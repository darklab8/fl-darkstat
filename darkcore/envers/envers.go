package envers

import (
	"encoding/json"
	"fmt"

	"github.com/darklab8/fl-darkstat/darkcore/envers/darkflag"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	map_settings "github.com/darklab8/fl-darkstat/darkmap/settings"
	stat_settings "github.com/darklab8/fl-darkstat/darkstat/settings"

	"github.com/darklab8/go-utils/utils/enverant"
)

var Enverants = []*enverant.Enverant{
	stat_settings.Env.Enver,
	stat_settings.Env.DarkcoreEnvVars.Enver,
	stat_settings.Env.ConfEnvVars.Enver,
	stat_settings.Env.UtilsEnvs.Enver,
	map_settings.Env.Enver,
}

var Envs = []any{
	stat_settings.Env,
	stat_settings.Env.DarkcoreEnvVars,
	stat_settings.Env.ConfEnvVars,
	stat_settings.Env.UtilsEnvs,
	map_settings.Env,
}

func init() {
	darkflag.Parse()
}

func PrintSettings() {
	env_data, err := json.Marshal(Envs)
	logus.Log.CheckPanic(err, "failed encoding settings")
	fmt.Println("freelancer_folder=", stat_settings.Env.FreelancerFolder, " configuration=", string(env_data))
}
