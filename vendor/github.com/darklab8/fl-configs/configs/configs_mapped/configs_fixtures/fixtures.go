package configs_fixtures

import (
	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

func FixtureGameLocation() utils_types.FilePath {
	current_folder := utils.GetCurrentFolder()
	game_location := utils_filepath.Dir(utils_filepath.Dir(utils_filepath.Dir(current_folder)))
	return game_location
}
