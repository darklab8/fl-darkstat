package main

import (
	"flag"
	"os"

	"github.com/darklab8/fl-darkstat/helpers/patch_disco"
)

func main() {
	f := flag.String("wd", ".", "...")
	flag.Parse()

	err := patch_disco.RunAutopatcher(*f)
	if err != nil {
		patch_disco.Log.CheckError(err, "failed to run autopatcher")
		os.Exit(1)
	}
}
