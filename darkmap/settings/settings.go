package settings

import (
	_ "embed"

	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/fl-darkstat/darkcore/envers/darkflag"
	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkmapEnvVars struct {
	utils_settings.UtilsEnvs
	configs_settings.ConfEnvVars
	SiteRoot string
	IndexUrl string
	Enver    *enverant.Enverant
}

var Env DarkmapEnvVars

func init() {
	env := enverant.NewEnverant(enverant.WithPrefix("DARKMAP_"))
	Env = DarkmapEnvVars{
		Enver:       env,
		UtilsEnvs:   utils_settings.GetEnvs(),
		ConfEnvVars: configs_settings.GetEnvs(),
		SiteRoot:    env.GetStr("SITE_ROOT", enverant.OrStr(*darkflag.ArgMapRoot)),
		IndexUrl:    env.GetStr("INDEX_URL", enverant.OrStr("index.html"), enverant.WithDesc("change map index filename, but only for `map web` and `map build` commands")),
	}
}
