package configs_settings

import (
	"os"

	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

var FallbackInfonamesToNickname bool

var Strict bool

func init() {
	FallbackInfonamesToNickname = os.Getenv("CONFIGS_FALLBACK_TO_NICKNAMES") == "true"

	Strict = os.Getenv("CONFIGS_STRICT") != "false"
}

func GetGameLocation() utils_types.FilePath {
	var folder utils_types.FilePath = utils_types.FilePath(os.Getenv("FREELANCER_FOLDER"))

	if folder == "" {
		workdir, _ := os.Getwd()
		folder = utils_types.FilePath(workdir)
	}

	return folder
}
