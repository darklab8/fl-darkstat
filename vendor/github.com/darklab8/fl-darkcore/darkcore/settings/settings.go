package settings

import (
	"fmt"

	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkcoreEnvVars struct {
	utils_settings.UtilsEnvs
}

var Env DarkcoreEnvVars

func init() {
	env := enverant.NewEnverant()
	Env = DarkcoreEnvVars{
		UtilsEnvs: utils_settings.GetEnvs(env),
	}
	fmt.Sprintln("conf=", Env)
}
