package settings

import (
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
