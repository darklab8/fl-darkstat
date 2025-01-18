package configs_export

import (
	"fmt"
	"strings"
	"testing"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestGetShips(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)
	ids := exporter.GetTractors()
	thrusters := exporter.GetThrusters(ids)

	var TractorsByID map[cfgtype.TractorID]Tractor = make(map[cfgtype.TractorID]Tractor)
	for _, tractor := range ids {
		TractorsByID[tractor.Nickname] = tractor
	}
	items := exporter.GetShips(ids, TractorsByID, thrusters)
	assert.Greater(t, len(items), 0)

	for _, item := range items {
		if strings.Contains(item.Nickname, "medium_miner") {
			fmt.Println()
		}
	}

	for _, item := range items {
		if strings.Contains(item.Nickname, "dsy_li_cruiser") {
			fmt.Println()
		}
	}
}
