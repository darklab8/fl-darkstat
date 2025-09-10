package configs_mapped

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/discovery/pob_goods"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	timeit.NewTimerF(func() {
		configs := TestFixtureConfigs()
		configs.Write(IsDruRun(true))
	})
}

func TestOnRealData(t *testing.T) {
	if true {
		return
	}
	file_public_bases := file.NewWebFile(PobDataUrl)
	config, err := pob_goods.Read(file_public_bases)

	assert.Greater(t, len(config.BasesByName), 0)
	assert.Nil(t, err)
}
