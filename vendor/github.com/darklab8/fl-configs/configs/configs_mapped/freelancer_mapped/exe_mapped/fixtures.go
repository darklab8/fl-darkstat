package exe_mapped

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

func FixtureFLINIConfig() *Config {
	test_directory := utils.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, FILENAME_FL_INI)))
	config := Read(iniload.NewLoader(fileref).Scan())
	return config
}
