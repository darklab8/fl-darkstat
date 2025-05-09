package configs_export

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
)

func TestGetEngines(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)
	ids := exporter.GetTractors()
	engines := exporter.GetEngines(ids)
	// assert.Greater(t, len(items), 0) # vanilla can have zero

	count := 0
	for _, engine := range engines {
		infocard := exporter.GetInfocard(infocarder.InfocardKey(engine.Nickname))
		if !strings.Contains(infocard.StringsJoin(""), strconv.Itoa(engine.CruiseSpeed)) {
			fmt.Println("engine, nick=", engine.Nickname, " not found correct cruise speed in infocard")
		}
		count++

	}
	fmt.Println("count=", count)
}
