package darkhttp

import (
	"testing"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated/statproto_deprecated"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/stretchr/testify/assert"
)

var IntegrationApiUrl = "https://darkstat.dd84ai.com"

func TestGetPoBs(t *testing.T) {
	c := NewClient(IntegrationApiUrl)
	pobs, err := c.GetPobs()
	logus.Log.CheckPanic(err, "failed to query pobs")

	assert.Greater(t, len(pobs), 0)
}

func TestGetScanners(t *testing.T) {
	c := NewClient(IntegrationApiUrl)
	items, err := c.GetScanners(pb.GetEquipmentInput{
		IncludeMarketGoods: true,
	})
	logus.Log.CheckPanic(err, "failed to query pobs")

	assert.Greater(t, len(items), 0, "non zero scanners must be")

	has_market_goods := false
	for _, item := range items {
		if len(item.MarketGoods) > 0 {
			has_market_goods = true
		}
	}
	assert.True(t, has_market_goods, "scanners must have market goods at least at some item")
}
