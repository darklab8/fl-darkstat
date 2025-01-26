package equipment_mapped

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile(FILENAME)

	config := Read([]*iniload.IniLoader{iniload.NewLoader(fileref).Scan()})

	assert.Greater(t, len(config.Commodities), 0, "expected finding commodities")
	assert.Greater(t, len(config.ShipHulls), 0)
	assert.Greater(t, len(config.Ships), 0)

	for _, commodity := range config.Commodities {
		commodity.Price.Get()
	}
}
