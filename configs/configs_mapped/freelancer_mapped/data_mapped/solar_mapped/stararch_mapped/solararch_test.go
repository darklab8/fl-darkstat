package stararch_mapped

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/darklab8/go-utils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile(FILENAME)
	loaded_market_ships := Read([]*iniload.IniLoader{iniload.NewLoader(fileref).Scan()})
	assert.Greater(t, len(loaded_market_ships.Stars), 0, "expected finding some elements")
	assert.Greater(t, len(loaded_market_ships.GlowsByNick), 0, "expected finding some elements")
}

const (
	FILENAME utils_types.FilePath = "stararch.ini"
)
