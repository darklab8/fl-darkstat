package darkflag

import (
	"flag"
	"fmt"
	"os"
)

var (
	ArgFreelancerFolder   = flag.String("freelancer-folder", "", "path to Freelancer folder root for data parsing. By default grabs current workdir")
	ArgFreelancerFallback = flag.String("freelancer-fallback", "", "if some configs aren't defined in first freelancer folder, grab from this one. Useful for FLSR usage in CI")
)

var (
	ArgWebPort           = flag.Int("web-port", 8000, "Main web port")
	ArgEnableUnixSockets = flag.Bool("unix-sockets", false, "Enable unix sockets for connections")
	ArgPassword          = flag.String("password", "", "protect access to web interface of darkstat with ?password=query_param")
)

var (
	ArgMapRoot = flag.String("map-site-root", "/", "map site root")
)

var (
	TradeDealsEnabled       = flag.Bool("stat-deals", true, "flag to show or not best trade deals in stat service. PERFORMANCE HEAVY. disable if not needed")
	StatSiteRoot            = flag.String("stat-site-root", "/", "useful if wishing serving darkstat from github pages sub urls. Makes sure correct link addresses")
	TradeDealsDetailedLanes = flag.Bool("stat-trade-detailed-lanes", false, "experimental option that allows to recieve more precise graph calculations by treating trade lane segments separately. Performance heavy.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Run `darkstat help` or `go run . help` for all env var and cli args possible to use\n")
	}

	if flag.Lookup("test.v") == nil {
		flag.Parse()
	} else {
		fmt.Println("run under go test")
	}

}
