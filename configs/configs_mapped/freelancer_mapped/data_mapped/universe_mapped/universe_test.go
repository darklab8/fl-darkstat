/*
parse universe.ini
*/
package universe_mapped

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/tests"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	filesystem := tests.FixtureFileFind()
	fileref := filesystem.GetFile(FILENAME)
	config := Read(iniload.NewLoader(fileref).Scan(), filesystem)

	assert.Greater(t, len(config.Bases), 0)
	assert.Greater(t, len(config.Systems), 0)
}

func TestIdentifySystemFiles(t *testing.T) {

	filesystem := tests.FixtureFileFind()
	logus.Log.Debug("filefind.FindConfigs" + fmt.Sprintf("%v", filesystem))

	universe_fileref := filesystem.GetFile(FILENAME)
	_ = Read(iniload.NewLoader(universe_fileref).Scan(), filesystem)
}
