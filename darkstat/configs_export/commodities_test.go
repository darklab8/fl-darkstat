package configs_export

import (
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestExportCommodities(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	items := exporter.GetCommodities()
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
	}

}
