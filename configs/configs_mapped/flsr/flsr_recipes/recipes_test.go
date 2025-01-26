package flsr_recipes

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile(FILENAME)

	if fileref == nil {
		return
	}

	config := Read(iniload.NewLoader(fileref).Scan())
	assert.Greater(t, len(config.Products), 0, "expected finding some elements")
}
