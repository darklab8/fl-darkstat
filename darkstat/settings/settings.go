package settings

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/darklab8/fl-configs/configs/configs_settings"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

//go:embed version.txt
var version string

type Configs struct {
	TractorTabName   string
	IsDevEnv         bool
	SiteRoot         string
	AppHeading       string
	FreelancerFolder utils_types.FilePath
	AppVersion       string
}

var Conf Configs

func init() {
	Conf = Configs{
		TractorTabName:   getEnvWithDefault("DARKSTAT_TRACTOR_TAB_NAME", "Tractors"),
		IsDevEnv:         os.Getenv("DEV") == "true",
		SiteRoot:         getEnvWithDefault("SITE_ROOT", "/"),
		AppHeading:       os.Getenv("FLDARKSTAT_HEADING"),
		FreelancerFolder: getFreelancerFolder(),
		AppVersion:       getAppVersion(),
	}
	fmt.Sprintln("conf=", Conf)
}

func getFreelancerFolder() utils_types.FilePath {
	return configs_settings.GetGameLocation()
}

func getAppVersion() string {
	// cleaning up version from... debugging logs used during dev env
	lines := strings.Split(version, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "v") {
			return line
		}
	}
	return version
}

func getEnvWithDefault(key string, default_ string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		return default_
	}
}
