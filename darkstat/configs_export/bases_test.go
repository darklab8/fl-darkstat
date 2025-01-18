package configs_export

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

type FixtureBasesOutput struct {
	configs *configs_mapped.MappedConfigs
	expoter *Exporter
	bases   []*Base
}

func FixtureBases(t *testing.T) FixtureBasesOutput {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	bases := exporter.GetBases()
	assert.Greater(t, len(bases), 0)
	assert.NotEqual(t, bases[0].Nickname, bases[1].Nickname)

	found_goods := false
	for _, base := range bases {
		if len(base.MarketGoodsPerNick) > 0 {
			found_goods = true
		}
	}
	assert.True(t, found_goods, "expected finding some goods")

	for _, base := range bases {
		if base.Nickname == "Br01_01_Base" {
			lines := base.Infocard
			fmt.Println(base.Nickname, lines)
			assert.Greater(t, len(lines), 0, "expected finding lines in infocard")
			break
		}
	}
	return FixtureBasesOutput{
		configs: configs,
		expoter: exporter,
		bases:   bases,
	}
}

func TestExportBases(t *testing.T) {
	_ = FixtureBases(t)
}
