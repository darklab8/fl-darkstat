package mbases_mapped

import (
	"fmt"
	"sort"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/stretchr/testify/assert"
)

func TestGetRepHacks(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile(FILENAME)

	config := Read(iniload.NewLoader(fileref).Scan())
	assert.Greater(t, len(config.Bases), 0, "expected finding some elements")

	// configs := configs_mapped.TestFixtureConfigs()
	// exporter := configs_export.NewExporter(configs)
	// bases := exporter.GetBases(configs_export.NoNameIncluded(true))

	faction_rephacks := FactionBribes(config)

	fmt.Println("printing for br_p_grp")
	chances := make([]BaseChance, 0, len(faction_rephacks["br_p_grp"]))

	for base, chance := range faction_rephacks["br_p_grp"] {
		chances = append(chances, BaseChance{
			Base:   base,
			Chance: chance,
		})
	}
	sort.Slice(chances, func(i, j int) bool {
		return chances[i].Chance > chances[j].Chance
	})

	for _, chance := range chances {
		var name string
		fmt.Println(chance.Base, " = ", 100*chance.Chance, " ", name)
	}
}
