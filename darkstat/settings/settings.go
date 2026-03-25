package settings

import (
	"flag"
	"fmt"
	"strings"

	_ "embed"

	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/fl-darkstat/darkcore/envers/darkflag"
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

	TractorTabName string
	SiteHost       string
	SiteRoot       string
	SiteUrl        string
	SiteHtmlTitle  string

	AppHeading string
	AppVersion string

	RelayHost     string
	RelayRoot     string
	RelayLoopSecs int

	TradeDealsEnabled bool

	TradeRoutesDetailedTradeLane    bool
	TradeRoutesBestDisablePobs      bool
	TradeRoutesBestTwoWaysLimitPobs int
	TradeRoutesBestDisableLiners    bool

	IsCPUProfilerEnabled bool
	IsMemProfilerEnabled bool

	IsStaticSiteGenerator bool
	Enver                 *enverant.Enverant
}

func IsApiActive() bool {
	if Env.IsStaticSiteGenerator {
		return false
	}
	return true
}

var Env DarkstatEnvVars

// var Enverants []*enverant.Enverant

var (
	TradeDealsEnabled       = flag.Bool("stat-deals", true, "flag to show or not best trade deals in stat service. PERFORMANCE HEAVY. disable if not needed")
	StatSiteRoot            = flag.String("stat-site-root", "/", "useful if wishing serving darkstat from github pages sub urls. Makes sure correct link addresses")
	TradeDealsDetailedLanes = flag.Bool("stat-trade-detailed-lanes", false, "experimental option that allows to recieve more precise graph calculations by treating trade lane segments separately. Performance heavy.")
)

func init() {
	darkflag.Parse()

	env := enverant.NewEnverant(enverant.WithPrefix("DARKSTAT_"), enverant.WithDescription("DARKSTAT set of envs for web interface for Freelancer game data navigation"))

	site_host := env.GetStr("SITE_HOST", enverant.OrStr(""), enverant.WithDesc("to show correct Swagger url/some buttons/links. Expects values with https part"))
	site_root := env.GetStr("SITE_ROOT", enverant.OrStr(*StatSiteRoot), enverant.WithDesc("useful if wishing serving darkstat from github pages sub urls. Makes sure correct link addresses"))
	Env = DarkstatEnvVars{
		Enver:           env,
		UtilsEnvs:       utils_settings.GetEnvs(),
		ConfEnvVars:     configs_settings.GetEnvs(),
		DarkcoreEnvVars: darkcore_settings.GetEnvs(),
		TractorTabName:  env.GetStr("TRACTOR_TAB_NAME", enverant.OrStr("Tractors"), enverant.WithDesc("name of Tractors tab to show in darkstat web")),

		SiteHost:      site_host,
		SiteRoot:      site_root,
		SiteUrl:       env.GetStrOr("SITE_URL", site_host+site_root, enverant.WithDesc("combined shortcut of site_host + site_root")),
		SiteHtmlTitle: env.GetStrOr("SITE_HTML_TITLE", "darkstat", enverant.WithDesc("site html title of a page")),
		AppHeading:    env.GetStr("FLDARKSTAT_HEADING", enverant.OrStr(""), enverant.WithDesc("What to show at the top of darkstat web UI. Possible to input any html")),
		AppVersion:    getAppVersion(),

		RelayHost:     env.GetStr("RELAY_HOST", enverant.OrStr(""), enverant.WithDesc("used to define relay url like with htps included. Makes sure that u deployed darkstat as static assets, they will still lead to relay backend to serve dynamic data. Useful for Discovery related deployment")),
		RelayRoot:     env.GetStr("RELAY_ROOT", enverant.OrStr("/"), enverant.WithDesc("if u ever will need to serve relay from non root path, u could use it to make sure requests go correct path.")),
		RelayLoopSecs: env.GetIntOr("RELAY_LOOP_SECS", 30, enverant.WithDesc("How often to update backend info during active app. Used for discovery to update PoB related info on a run")),

		TradeDealsEnabled: env.GetBoolOr("TRADE_DEALS_ENABLED", *TradeDealsEnabled, enverant.WithDesc("enable calculating one way and two way best trades? PERFORMANCE HEAVY. by default off. cli args must be put before command like `web`")),

		TradeRoutesDetailedTradeLane:    env.GetBoolOr("TRADE_ROUTES_DETAILED_TRADE_LANE", *TradeDealsDetailedLanes, enverant.WithDesc("experimental option that allows to recieve more precise graph calculations by treating trade lane segments separately. Performance heavy.")),
		TradeRoutesBestDisablePobs:      env.GetBoolOr("DISABLE_POBS_FOR_BEST_TRADES", false, enverant.WithDesc("if u use discovery mod, an option to turn off pobs from best trades")),
		TradeRoutesBestTwoWaysLimitPobs: env.GetIntOr("TRADE_ROUTES_BEST_TWO_WAY_LIMIT_POBS", 99, enverant.WithDesc("Limit amount of pobs participating in 2 way routes")),
		TradeRoutesBestDisableLiners:    env.GetBoolOr("TRADE_ROUTES_DISABLE_LINERS", false, enverant.WithDesc("")),
	}

	if !Env.TradeDealsEnabled {
		fmt.Println("WARN: TRADE_DEALS_ENABLED remained off. use env var set true, or cli arg `-deals` to turn on BEST TRADE DEALS")
	}
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
