package main

import (
	"flag"
	"os"

	"github.com/darklab8/fl-darkstat/helpers/patch_disco"
)

func main() {
	f := flag.String("wd", ".", "...")
	use_cache := flag.Bool("cache", false, "use cached version")
	flag.Parse()

	err := patch_disco.RunAutopatcher(*f, *use_cache)
	if err != nil {
		patch_disco.Log.CheckError(err, "failed to run patch disco")
		os.Exit(1)
	}
}
