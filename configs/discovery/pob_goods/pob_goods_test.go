package pob_goods

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	test_directory := utils_os.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, "example.json")))

	config, err := Read(fileref)
	assert.Greater(t, len(config.BasesByName), 0)
	assert.Nil(t, err)
}

func TestReader2(t *testing.T) {
	test_directory := utils_os.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, "example2.json")))

	config, err := Read(fileref)
	assert.Greater(t, len(config.BasesByName), 0)
	assert.Nil(t, err)
}

func TestReaderThread(t *testing.T) {
	test_directory := utils_os.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, "example_thread.json")))

	config, err := Read(fileref)
	assert.Greater(t, len(config.BasesByName), 0)
	assert.Nil(t, err)

	assert.NotNil(t, config.BasesByName["Fortitudine"].ForumThreadUrl)
}
