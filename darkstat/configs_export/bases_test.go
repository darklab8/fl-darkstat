package configs_export

import (
	"fmt"
	"strings"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_settings/logus"
	"github.com/darklab8/fl-configs/configs/discovery/discoprices"
	"github.com/darklab8/go-typelog/typelog"
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

	return FixtureBasesOutput{
		configs: configs,
		expoter: exporter,
		bases:   bases,
	}
}

func TestExportBases(t *testing.T) {
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

	if configs.Discovery != nil {
		aralsk_shipyard := bases_by_nick["ew09_03_base"]
		aralsk_shipyard_universe_base, ok := configs.Universe.BasesMap[universe_mapped.BaseNickname(aralsk_shipyard.Nickname)]
		assert.True(t, ok)
		_, aralsk_shipyard_trader_exists := aralsk_shipyard_universe_base.ConfigBase.RoomMapByRoomNickname["trader"]
		assert.False(t, aralsk_shipyard_trader_exists)
	}

	// toggle func TraderExists not removing market goods in order to see data here, make func Trader exists always returning true.
	for _, base := range bases {
		universe_base := configs.Universe.BasesMap[universe_mapped.BaseNickname(base.Nickname)]
		if len(base.MarketGoodsPerNick) > 0 && !universe_base.TraderExists {
			logus.Log.Warn("____",
				typelog.Any("base_nickname", base.Nickname),
				typelog.Any("base_name", base.Name),
			)
		}
	}
}

func TestServerOverrides(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	if configs.Discovery == nil {
		return
	}

	content := `
[Price]
MarketGood = li01_01_base, commodity_basic_alloys, 1150, 1550, 1
`
	memory_file := file.NewMemoryFile(strings.Split(content, "\n"))
	scanned_file := iniload.NewLoader(memory_file).Scan()
	discoPrices := discoprices.Read(scanned_file)
	// Adding to main Freelancer instance..
	configs.Discovery.Prices = discoPrices

	exporter := NewExporter(configs)

	bases := exporter.GetBases()
	commodities := exporter.GetCommodities()
	EnhanceBasesWithServerOverrides(bases, commodities)

	var targetbase *Base
	for _, base := range bases {
		if base.Nickname == "li01_01_base" {
			targetbase = base
		}
	}

	commodity_nickname := GetCommodityKey("commodity_basic_alloys", -1)
	alloy := targetbase.MarketGoodsPerNick[commodity_nickname]
	assert.Equal(t, 1550, alloy.PriceBaseSellsFor)
	assert.True(t, alloy.IsServerSideOverride)

	var targetcom *Commodity
	for _, com := range commodities {
		if com.Nickname == "commodity_basic_alloys" {
			targetcom = com
		}
	}
	market_good := targetcom.Bases["li01_01_base"]
	assert.Equal(t, 1550, market_good.PriceBaseSellsFor)
	assert.True(t, market_good.IsServerSideOverride)
}
