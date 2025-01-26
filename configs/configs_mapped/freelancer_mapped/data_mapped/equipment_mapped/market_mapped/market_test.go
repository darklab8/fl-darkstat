package market_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"

	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile(FILENAME_SHIPS)

	loaded_market_ships := Read([]*iniload.IniLoader{iniload.NewLoader(fileref).Scan()})

	assert.Greater(t, len(loaded_market_ships.BaseGoods), 0, "market ships sections were not scanned")
}

func TestWriter(t *testing.T) {
	input_file := tests.FixtureFileFind().GetFile(FILENAME_SHIPS)

	temp_directory := utils_os.GetCurrrentTempFolder()

	config := Read([]*iniload.IniLoader{iniload.NewLoader(input_file).Scan()})
	config.Files[0].SetOutputPath(utils_filepath.Join(temp_directory, FILENAME_SHIPS))
	config.Write()
}
