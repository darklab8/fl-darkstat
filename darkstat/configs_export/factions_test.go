package configs_export

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestFaction(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	items := exporter.GetFactions([]*Base{})
	assert.Greater(t, len(items), 0)

	infocards := exporter.Infocards
	for _, faction := range items {
		if faction.Nickname == "br_m_grp" {
			lines := infocards[faction.Infocard]
			fmt.Println(faction.Nickname, lines)
			assert.Greater(t, len(lines), 0)
			break
		}

	}
}
