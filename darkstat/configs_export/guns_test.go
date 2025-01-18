package configs_export

import (
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestGetGuns(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	ids := exporter.GetTractors()
	shields := exporter.GetShields(ids)
	buyable_shield_tech := exporter.GetBuyableShields(shields)
	guns := exporter.GetGuns(ids, buyable_shield_tech)
	assert.Greater(t, len(guns), 0)
	// exporter.infocards_parser.Get()
}
