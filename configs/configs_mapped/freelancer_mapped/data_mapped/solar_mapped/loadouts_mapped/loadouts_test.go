package loadouts_mapped

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile(FILENAME)

	loaded_market_ships := Read([]*iniload.IniLoader{iniload.NewLoader(fileref).Scan()})

	assert.Greater(t, len(loaded_market_ships.Loadouts), 0, "expected finding some elements")
}
