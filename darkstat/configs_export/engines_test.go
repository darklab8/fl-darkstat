package configs_export

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
)

func TestGetEngines(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)
	ids := exporter.GetTractors()
	engines := exporter.GetEngines(ids)
	// assert.Greater(t, len(items), 0) # vanilla can have zero

	count := 0
	for _, engine := range engines {
		infocard := exporter.Infocards[InfocardKey(engine.Nickname)]
		if !strings.Contains(infocard.StringsJoin(""), strconv.Itoa(engine.CruiseSpeed)) {
			fmt.Println("engine, nick=", engine.Nickname, " not found correct cruise speed in infocard")
		}
		count++

	}
	fmt.Println("count=", count)
}
