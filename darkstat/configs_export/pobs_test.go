package configs_export

import (
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestPobGoods(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)
	if configs.Discovery == nil {
		return
	}
	items := exporter.GetPoBs()
	assert.Greater(t, len(items), 0)
}
