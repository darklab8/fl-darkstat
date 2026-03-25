package settings

import (
	_ "embed"
	"flag"

	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/fl-darkstat/darkcore/envers/darkflag"
	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkmapEnvVars struct {
	utils_settings.UtilsEnvs
	configs_settings.ConfEnvVars
	SiteRoot string
	Enver    *enverant.Enverant
}

var Env DarkmapEnvVars

var (
	ArgMapRoot = flag.String("map-site-root", "/", "map site root")
)

func init() {
	darkflag.Parse()
	env := enverant.NewEnverant(enverant.WithPrefix("DARKMAP_"))
	Env = DarkmapEnvVars{
		Enver:       env,
		UtilsEnvs:   utils_settings.GetEnvs(),
		ConfEnvVars: configs_settings.GetEnvs(),
		SiteRoot:    env.GetStr("SITE_ROOT", enverant.OrStr(*ArgMapRoot)),
	}
}
