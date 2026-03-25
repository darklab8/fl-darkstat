package darkflag

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Run `darkstat help` or `go run . help` for all env var and cli args possible to use\n")
	}
}

func Parse() {
	flag.Parse()
}
