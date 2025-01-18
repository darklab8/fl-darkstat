package configs_export

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestGetOres(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)
	commodities := exporter.GetCommodities()
	mining_operations := exporter.GetOres(commodities)
	assert.Greater(t, len(mining_operations), 0)
	fmt.Println("len(mining_operations)=", len(mining_operations))
}
