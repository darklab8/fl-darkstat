package darkflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	ArgFreelancerFolder   = flag.String("freelancer-folder", "", "path to Freelancer folder root for data parsing. By default grabs current workdir")
	ArgFreelancerFallback = flag.String("freelancer-fallback", "", "if some configs aren't defined in first freelancer folder, grab from this one. Useful for FLSR usage in CI")
)

var (
	ArgWebPort           = flag.Int("web-port", 8000, "Main web port")
	ArgEnableUnixSockets = flag.Bool("unix-sockets-on", false, "Enable unix sockets for connections")
	ArgPassword          = flag.String("password", "", "protect access to web interface of darkstat with ?password=query_param")
)

var (
	ArgMapRoot = flag.String("map-site-root", "/", "map site root")
)

var (
	TradeDealsEnabled                       = flag.Bool("stat-deals-on", false, "flag to show or not best trade deals in stat service. PERFORMANCE HEAVY. by default off. disable if not needed")
	StatSiteRoot                            = flag.String("stat-site-root", "/", "useful if wishing serving darkstat from github pages sub urls. Makes sure correct link addresses")
	TradeDealsDetailedLanes                 = flag.Bool("stat-trade-detailed-lanes-on", false, "experimental option that allows to recieve more precise graph calculations by treating trade lane segments separately. Performance heavy.")
	IsExperimentalMapRunWithDarkstatEnabled = flag.Bool("experimental-map-on", false, "flag to turn on map as part of darkstat. VERY EXPERIMENTAL: may lead to drastic CPU and RAM performance issues, running `map web` separately is recommended. PERFORMANCE HEAVY. use standalone map through `map web` command if u wish it faster")
	StatDefaultTheme                        = flag.String("stat-default-theme", "white", "default shell theme for index.html redirect: white, dark, or vanilla. localStorage darkstat-theme overrides when set. flags before web or build")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Run `darkstat help` or `go run . help` for all env var and cli args possible to use\n")
	}

	is_test := false

	for _, arg := range os.Args {
		if strings.Contains(arg, "-test.") {
			is_test = true
		}
	}

	if flag.Lookup("test.v") != nil || flag.Lookup("test.testlogfile") != nil || is_test {
		fmt.Println("run under go test")
	} else {
		flag.Parse()

		args := os.Args[1:]
		if len(args) > 0 {
			index_of_last_non_flag := 0
			index_of_last_flag := 0

			for arg_index, arg := range args {
				if strings.HasPrefix(arg, "-") {
					// this is flag
					index_of_last_flag = arg_index
				} else {
					// this is not flag
					index_of_last_non_flag = arg_index
				}
			}

			if index_of_last_flag > index_of_last_non_flag {
				fmt.Println("ERROR flags like --stat-deals-on must be inputed BEFORE commands like `web` or `build`")
				os.Exit(1)
			}
		}
	}
}
