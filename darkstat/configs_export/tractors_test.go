package configs_export

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/discovery/playercntl_rephacks"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func TestGetTractors(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	items := exporter.GetTractors()
	assert.Greater(t, len(items), 0)
}

func TestReadCommentsInPlayercntlRephacks(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()

	if configs.Discovery == nil {
		return
	}

	test_directory := utils_os.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, "playercntl_rephacks.cfg")))

	config_file := iniload.NewLoader(fileref).Scan()
	config_read := playercntl_rephacks.Read(config_file)
	configs.Discovery.PlayercntlRephacks = config_read
	exporter := NewExporter(configs)

	items := exporter.GetTractors()
	assert.Greater(t, len(items), 0)
}
