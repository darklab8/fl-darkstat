package settings

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	darkcore_settings "github.com/darklab8/fl-darkstat/darkcore/settings"

	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

//go:embed version.txt
var version string

type DarkstatEnvVars struct {
	utils_settings.UtilsEnvs
	configs_settings.ConfEnvVars
	darkcore_settings.DarkcoreEnvVars

	TractorTabName    string
	SiteUrl           string
	SiteRoot          string
	SiteRootAcceptors string
	AppHeading        string
	AppVersion        string
	IsDetailed        bool

	RelayHost     string
	RelayRoot     string
	RelayLoopSecs int

	IsDisabledTradeRouting       bool
	TradeRoutesDetailedTradeLane bool

	IsCPUProfilerEnabled bool
	IsMemProfilerEnabled bool
}

var Env DarkstatEnvVars

func init() {
	env := enverant.NewEnverant()
	Env = DarkstatEnvVars{
		UtilsEnvs:         utils_settings.GetEnvs(env),
		ConfEnvVars:       configs_settings.GetEnvs(env),
		DarkcoreEnvVars:   darkcore_settings.GetEnvs(env),
		TractorTabName:    env.GetStr("DARKSTAT_TRACTOR_TAB_NAME", enverant.OrStr("Tractors")),
		SiteUrl:           env.GetStrOr("SITE_URL", env.GetStr("SITE_HOST", enverant.OrStr(""))),
		SiteRoot:          env.GetStr("SITE_ROOT", enverant.OrStr("/")),
		SiteRootAcceptors: env.GetStr("SITE_ROOT_ACCEPTORS", enverant.OrStr("")),
		AppHeading:        env.GetStr("FLDARKSTAT_HEADING", enverant.OrStr("")),
		AppVersion:        getAppVersion(),
		IsDetailed:        env.GetBoolOr("DARKSTAT_DETAILED", false),

		RelayHost:     env.GetStr("RELAY_HOST", enverant.OrStr("")),
		RelayRoot:     env.GetStr("RELAY_ROOT", enverant.OrStr("/")),
		RelayLoopSecs: env.GetIntOr("RELAY_LOOP_SECS", 30),

		TradeRoutesDetailedTradeLane: env.GetBoolOr("DARKSTAT_TRADE_ROUTES_DETAILED_TRADE_LANE", false),
		IsDisabledTradeRouting:       env.GetBoolOr("CONFIGS_DISABLE_TRADE_ROUTES", false), // BROKEN. DO NOT TURN THIS FEATURE ON.

		IsCPUProfilerEnabled: env.GetBoolOr("IS_CPU_PROFILER_ENABLED", false),
		IsMemProfilerEnabled: env.GetBoolOr("IS_MEM_PROFILER_ENABLED", false),
	}

	fmt.Sprintln("conf=", Env)
}

func (e DarkstatEnvVars) GetSiteRootAcceptors() []string {
	if e.SiteRootAcceptors == "" {
		return []string{}
	}

	return strings.Split(e.SiteRootAcceptors, ",")
}

func IsRelayActive(configs *configs_mapped.MappedConfigs) bool {
	if configs.Discovery != nil {
		fmt.Println("discovery always wishes to see pobs")
		return true
	}

	fmt.Println("relay is disabled")
	return false
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
