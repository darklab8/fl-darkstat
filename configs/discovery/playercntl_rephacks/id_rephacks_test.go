package playercntl_rephacks

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	test_directory := utils_os.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, "playercntl_rephacks.cfg")))

	config := Read(iniload.NewLoader(fileref).Scan())

	assert.Greater(t, len(config.DefaultReps), 0)

	assert.Greater(t, len(config.RephacksByID), 0)

	for rephack_index, rephack := range config.RephacksByID {
		_ = rephack_index
		rephack.ID.Get()
		for _, faction := range rephack.Reps {
			faction.Rep.Get()
			faction.GetRepType()
		}
	}

}
