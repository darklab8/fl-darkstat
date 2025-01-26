package overrides

import (
	"testing"

	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func TestReadingIt(t *testing.T) {
	test_directory := utils_os.GetCurrrentTestFolder()

	overrides := Read(utils_types.FilePath(utils_filepath.Join(test_directory, FILENAME)))

	assert.Equal(t, 0.33, overrides.GetSystemSpeedMultiplier("test_nickname"))
	assert.Equal(t, 1.0, overrides.GetSystemSpeedMultiplier("another_nickname"))
}
