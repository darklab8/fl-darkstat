package configs_export

import (
	"context"
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestGetOres(t *testing.T) {
	ctx := context.Background()
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)
	commodities := exporter.GetCommodities(ctx)
	mining_operations := exporter.GetOres(ctx, commodities)
	assert.Greater(t, len(mining_operations), 0)
	fmt.Println("len(mining_operations)=", len(mining_operations))
}
