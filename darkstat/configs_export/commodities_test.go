package configs_export

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestExportCommodities(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	items := exporter.GetCommodities()
	assert.Greater(t, len(items), 0)

	fmt.Println(items[0])

	for _, item := range items {
		if item.Nickname == "commodity_diamonds" {
			fmt.Println()
		}
	}
}
