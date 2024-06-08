package settings

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/darklab8/fl-configs/configs/configs_settings"
	"github.com/darklab8/go-utils/utils/utils_env"
)

//go:embed version.txt
var version string

type DarkstatEnvVars struct {
	configs_settings.ConfEnvVars
	TractorTabName string
	IsDevEnv       utils_env.EnvBool
	SiteRoot       string
	AppHeading     string
	AppVersion     string
}

var Env DarkstatEnvVars

func init() {
	env := utils_env.NewEnvConfig()
	Env = DarkstatEnvVars{
		ConfEnvVars:    configs_settings.Env,
		TractorTabName: env.GetEnvWithDefault("DARKSTAT_TRACTOR_TAB_NAME", "Tractors"),
		IsDevEnv:       env.GetEnvBool("DEV"),
		SiteRoot:       env.GetEnvWithDefault("SITE_ROOT", "/"),
		AppHeading:     env.GetEnv("FLDARKSTAT_HEADING"),
		AppVersion:     getAppVersion(),
	}
	fmt.Sprintln("conf=", Env)
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
