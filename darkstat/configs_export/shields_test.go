package configs_export

import (
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestGetShields(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)
	ids := exporter.GetTractors()
	items := exporter.GetShields(ids)
	assert.Greater(t, len(items), 0)
}
