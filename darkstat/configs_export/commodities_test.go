package configs_export

import (
	"context"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestExportCommodities(t *testing.T) {
	ctx := context.Background()
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	items := exporter.GetCommodities(ctx)
	assert.Greater(t, len(items), 0)

	if configs.Discovery != nil {
		var scidata *Commodity
		for _, item := range items {
			if item.Nickname == "commodity_sciencedata" {
				scidata = item
			}
		}
		assert.NotNil(t, scidata)
		assert.NotContains(t, scidata.Name, "\r")

		foundable, _ := exporter.findable_in_loot_cache["cr_heavy_battlerazor"]
		assert.True(t, foundable, "expected cr_heavy_battlerazor to be findable in loot")

		foundable, _ = exporter.findable_in_loot_cache["commodity_sciencedata"]
		assert.True(t, foundable, "expected commodity_sciencedata to be findable in loot")

		foundable, _ = exporter.findable_in_loot_cache["bs_heavy_w03"]
		assert.True(t, foundable, "expected bs_heavy_w03 to be findable in loot")

		foundable, _ = exporter.findable_in_loot_cache["commodity_sealed_container"]
		assert.True(t, foundable, "expected commodity_sealed_container to be findable in loot")
	}

}
