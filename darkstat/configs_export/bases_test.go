package configs_export

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
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
	bases_by_nick := make(map[string]*Base)
	for _, base := range bases {
		bases_by_nick[string(base.Nickname)] = base
	}

	aralsk_shipyard := bases_by_nick["ew09_03_base"]
	aralsk_shipyard_universe_base, ok := configs.Universe.BasesMap[universe_mapped.BaseNickname(aralsk_shipyard.Nickname)]
	assert.True(t, ok)
	_, aralsk_shipyard_trader_exists := aralsk_shipyard_universe_base.ConfigBase.RoomMapByRoomNickname["trader"]
	assert.False(t, aralsk_shipyard_trader_exists)

	return FixtureBasesOutput{
		configs: configs,
		expoter: exporter,
		bases:   bases,
	}
}

func TestExportBases(t *testing.T) {
	_ = FixtureBases(t)
}
