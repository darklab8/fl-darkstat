package systems_mapped

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"

	"github.com/stretchr/testify/assert"
)

func TestSaveRecycleParams(t *testing.T) {
	filesystem := tests.FixtureFileFind()

	universe_ini := iniload.NewLoader(file.NewFile(filesystem.Hashmap[universe_mapped.FILENAME].GetFilepath())).Scan()
	universe_config := universe_mapped.Read(universe_ini, filesystem)

	systems := Read(universe_config, filesystem)

	system, ok := systems.SystemsMap["br01"]
	assert.True(t, ok, "system should be present")

	system.Render()
}
