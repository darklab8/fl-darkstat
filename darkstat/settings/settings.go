package settings

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/darklab8/fl-configs/configs/configs_settings"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

func GetFreelancerFolder() utils_types.FilePath {
	return configs_settings.GetGameLocation()
}

//go:embed version.txt
var version string

var TractorTabName string

func init() {
	if value, ok := os.LookupEnv("DARKSTAT_TRACTOR_TAB_NAME"); ok {
		TractorTabName = value
	} else {
		TractorTabName = "Tractors"
	}

	fmt.Sprintln("settings.TractorTabName=", TractorTabName)
}

func GetVersion() string {
	// cleaning up version from... debugging logs used during dev env
	lines := strings.Split(version, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "v") {
			return line
		}
	}
	return version
}
