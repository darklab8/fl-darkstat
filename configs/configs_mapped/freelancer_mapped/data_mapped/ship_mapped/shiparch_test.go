package ship_mapped

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	files := tests.FixtureFileFind()

	fileref1 := files.GetFile("shiparch.ini")
	fileref2 := files.GetFile("rtc_shiparch.ini")

	config := Read([]*iniload.IniLoader{iniload.NewLoader(fileref1).Scan(), iniload.NewLoader(fileref2).Scan()})

	assert.Greater(t, len(config.Ships), 0)

	for _, commodity := range config.Ships {
		commodity.Nickname.Get()
	}
}
