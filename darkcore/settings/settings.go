package settings

import (
	"fmt"

	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkcoreEnvVars struct {
	utils_settings.UtilsEnvs
	Password            string
	Secret              string
	ExtraCookieHost     string
	IsDiscoOauthEnabled bool
	CacheControl        string
	Enver               *enverant.Enverant
}

var Env DarkcoreEnvVars

func GetEnvs() DarkcoreEnvVars {
	envs := enverant.NewEnverant(enverant.WithPrefix("DARKCORE_"), enverant.WithDescription("DARKCORE set of envs for a web framework based on templ to implement static site generator with backend fallback"))

	Env = DarkcoreEnvVars{
		UtilsEnvs:           utils_settings.GetEnvs(),
		Password:            envs.GetStrOr("PASSWORD", "", enverant.WithDesc("protect access to web interface of darkstat with ?password=query_param")),
		Secret:              envs.GetStrOr("SECRET", "passphrasewhichneedstobe32bytes!", enverant.WithDesc("secret to persist authentifications with query param password or oauth, required if using auths")),
		CacheControl:        envs.GetStrOr("CACHE_CONTROL", ""),
		IsDiscoOauthEnabled: envs.GetBool("DISCO_OAUTH", enverant.WithDesc("an option to turn auth of darkstat for Discovery freelancer a protected dev instance of darkstat")),
		Enver:               envs,
	}
	return Env
}

func init() {
	Env = GetEnvs()
	fmt.Sprintln("conf=", Env)
}
