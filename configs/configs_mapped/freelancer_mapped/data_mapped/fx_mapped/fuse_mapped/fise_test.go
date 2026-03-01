package fuse_mapped

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile2("fx/fuse.ini")
	fileref2 := tests.FixtureFileFind().GetFile2("fx/fuse_suprise_solar.ini")

	config := Read([]*iniload.IniLoader{
		iniload.NewLoader(fileref).Scan(),
		iniload.NewLoader(fileref2).Scan(),
	})

	assert.Greater(t, len(config.Fuses), 0, "expected finding fuses")
	assert.Greater(t, len(config.FuseMap), 0)

	for _, fuse := range config.Fuses {
		if len(fuse.LootableHardpoints) > 0 {
			fmt.Println("fuse.LootableHardpoints=", len(fuse.LootableHardpoints), " drop_cargo=", fuse.DoesDropCargo)
		}
	}
}
