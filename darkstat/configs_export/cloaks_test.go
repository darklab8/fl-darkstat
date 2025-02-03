package configs_export

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestGetCloaks(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)
	ids := exporter.GetTractors()
	items := exporter.GetCloaks(ids)
	if configs.Discovery == nil {
		return
	}
	assert.Greater(t, len(items), 0)
}
