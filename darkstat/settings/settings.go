package settings

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/darklab8/fl-configs/configs/configs_settings"
	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

//go:embed version.txt
var version string

type DarkstatEnvVars struct {
	utils_settings.UtilsEnvs
	configs_settings.ConfEnvVars
	TractorTabName string
	SiteRoot       string
	AppHeading     string
	AppVersion     string
}

var Env DarkstatEnvVars

func init() {
	env := enverant.NewEnverant()
	Env = DarkstatEnvVars{
		UtilsEnvs:      utils_settings.GetEnvs(env),
		ConfEnvVars:    configs_settings.GetEnvs(env),
		TractorTabName: env.GetStr("DARKSTAT_TRACTOR_TAB_NAME", enverant.OrStr("Tractors")),
		SiteRoot:       env.GetStr("SITE_ROOT", enverant.OrStr("/")),
		AppHeading:     env.GetStr("FLDARKSTAT_HEADING", enverant.OrStr("")),
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
