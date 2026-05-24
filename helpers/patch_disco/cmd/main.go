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
		os.Exit(1)
	}
}
