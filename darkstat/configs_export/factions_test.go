package configs_export

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/stretchr/testify/assert"
)

func TestFaction(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	items := exporter.GetFactions([]*Base{})
	assert.Greater(t, len(items), 0)

	exporter.GetInfocardsDict(func(infocards infocarder.Infocards) {
		for _, faction := range items {
			if faction.Nickname == "br_m_grp" {
				lines := infocards[faction.InfocardKey]
				fmt.Println(faction.Nickname, lines)
				assert.Greater(t, len(lines), 0)
				break
			}

		}
	})

}
