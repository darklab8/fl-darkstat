package diff2money

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	filesystem := tests.FixtureFileFind()
	fileref := filesystem.GetFile(FILENAME)

	loaded_market_ships := Read(iniload.NewLoader(fileref).Scan())

	assert.Greater(t, len(loaded_market_ships.DiffToMoney), 0, "expected finding some elements")
}
