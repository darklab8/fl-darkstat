package interface_mapped

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile(FILENAME_FL_INI)

	config := Read(iniload.NewLoader(fileref).Scan())

	assert.Greater(t, len(config.InfocardMapTable.Map), 0)
}

func TestReaderComments1(t *testing.T) {
	test_directory := utils_os.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, "infocardmap.comments.ini")))

	config := Read(iniload.NewLoader(fileref).Scan())

	assert.Greater(t, len(config.InfocardMapTable.Map), 0)
}

func TestReaderComments2(t *testing.T) {
	test_directory := utils_os.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, "infocardmap.comments2.ini")))

	config := Read(iniload.NewLoader(fileref).Scan())

	assert.Greater(t, len(config.InfocardMapTable.Map), 0)
}
