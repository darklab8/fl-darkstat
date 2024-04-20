package configs_fixtures

import (
	"os"

	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

func FixtureGameLocation() utils_types.FilePath {

	var folder utils_types.FilePath = utils_types.FilePath(os.Getenv("FREELANCER_FOLDER"))

	if folder == "" {
		workdir, _ := os.Getwd()
		folder = utils_types.FilePath(workdir)
	}

	return folder
}
