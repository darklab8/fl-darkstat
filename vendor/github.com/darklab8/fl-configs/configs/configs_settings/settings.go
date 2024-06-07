package configs_settings

import (
	"os"

	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type ConfEnvVars struct {
	FallbackInfonamesToNickname bool
	Strict                      bool
	FreelancerFolder            utils_types.FilePath
}

var Env ConfEnvVars

func init() {
	Env = ConfEnvVars{
		FallbackInfonamesToNickname: os.Getenv("CONFIGS_FALLBACK_TO_NICKNAMES") == "true",
		Strict:                      os.Getenv("CONFIGS_STRICT") != "false",
		FreelancerFolder:            getGameLocation(),
	}
}

func getGameLocation() utils_types.FilePath {
	var folder utils_types.FilePath = utils_types.FilePath(os.Getenv("FREELANCER_FOLDER"))

	if folder == "" {
		workdir, _ := os.Getwd()
		folder = utils_types.FilePath(workdir)
	}

	return folder
}
