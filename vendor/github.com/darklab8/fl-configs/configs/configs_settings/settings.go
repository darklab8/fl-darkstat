package configs_settings

import (
	"os"

	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_settings"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type ConfEnvVars struct {
	utils_settings.UtilsEnvs
	FallbackInfonamesToNickname bool
	Strict                      bool
	FreelancerFolder            utils_types.FilePath
}

var Env ConfEnvVars

func init() {
	env := enverant.NewEnverant(
		enverant.WithEnvFile(utils_os.GetCurrentFolder().Dir().Dir().Join(".enverant", "enverant.json").ToString()),
	)
	Env = GetEnvs(env)
}

func GetEnvs(envs *enverant.Enverant) ConfEnvVars {
	Env = ConfEnvVars{
		UtilsEnvs:                   utils_settings.GetEnvs(envs),
		FallbackInfonamesToNickname: envs.GetBool("CONFIGS_FALLBACK_TO_NICKNAMES", enverant.OrBool(false)),
		Strict:                      envs.GetBool("CONFIGS_STRICT", enverant.OrBool(true)),
		FreelancerFolder:            getGameLocation(envs),
	}
	return Env
}

func getGameLocation(envs *enverant.Enverant) utils_types.FilePath {
	var folder utils_types.FilePath = utils_types.FilePath(
		envs.GetStr("FREELANCER_FOLDER", enverant.OrStr("")),
	)

	if folder == "" {
		workdir, _ := os.Getwd()
		folder = utils_types.FilePath(workdir)
	}

	return folder
}
