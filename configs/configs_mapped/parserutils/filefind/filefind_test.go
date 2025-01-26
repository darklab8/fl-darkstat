package filefind

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/stretchr/testify/assert"
)

func TestDiscoverFiles(t *testing.T) {
	// Write some data example in order to remove integration flag
	test_directory := configs_settings.Env.FreelancerFolder
	filesystem := FindConfigs(test_directory)

	assert.GreaterOrEqual(t, len(filesystem.Files), 2, "expected more 2 files, fount smth else")
}
