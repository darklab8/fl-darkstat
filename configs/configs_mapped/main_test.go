package configs_mapped

import (
	"testing"

	"github.com/darklab8/go-utils/utils/timeit"
)

func TestSimple(t *testing.T) {
	timeit.NewTimerF(func() {
		configs := TestFixtureConfigs()
		configs.Write(IsDruRun(true))
	})
}
