package players

import (
	"context"
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

	ctx := context.Background()
	config, err := Read(ctx, fileref)
	assert.Greater(t, len(config.Players), 0)
	assert.Nil(t, err)
}
