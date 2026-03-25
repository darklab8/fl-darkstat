package settings

import (
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/envers/darkflag"
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
	EnableUnixSockets   bool
	Enver               *enverant.Enverant
	WebPort             int
	AppStart            time.Time
}

var Env DarkcoreEnvVars

func GetEnvs() DarkcoreEnvVars {
	envs := enverant.NewEnverant(enverant.WithPrefix("DARKCORE_"), enverant.WithDescription("DARKCORE set of envs for a web framework based on templ to implement static site generator with backend fallback"))

	Env = DarkcoreEnvVars{
		UtilsEnvs:           utils_settings.GetEnvs(),
		Password:            envs.GetStrOr("PASSWORD", *darkflag.ArgPassword, enverant.WithDesc("protect access to web interface of darkstat with ?password=query_param")),
		CacheControl:        envs.GetStrOr("CACHE_CONTROL", ""), // refactor to boolean and set as true
		IsDiscoOauthEnabled: envs.GetBool("DISCO_OAUTH", enverant.WithDesc("an option to turn auth of darkstat for Discovery freelancer a protected dev instance of darkstat")),
		Secret:              envs.GetStrOr("SECRET", "passphrasewhichneedstobe32bytes!", enverant.WithDesc("secret to persist authentifications with query param password or oauth, required if using auths")),
		EnableUnixSockets:   envs.GetBoolOr("ENABLE_UNIX_SOCKETS", *darkflag.ArgEnableUnixSockets, enverant.WithDesc("creating unix sockets, requires /tmp/darkstat or /tmp/darkstat-{environment} folder defined")),
		WebPort:             envs.GetIntOr("WEB_PORT", *darkflag.ArgWebPort, enverant.WithDesc("specify web port")),
		AppStart:            time.Now(),
		Enver:               envs,
	}
	return Env
}

func init() {
	Env = GetEnvs()
}
