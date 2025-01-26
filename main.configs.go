/*
See package `configs` for description and code examples
*/
package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

// from configs. Refactor to integrate it
func main2() {

	// for profiling only stuff.
	f, err := os.Create("prof.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < 1; i++ {
		timeit.NewTimerF(func() {
			var configs *configs_mapped.MappedConfigs
			timeit.NewTimerF(func() {
				freelancer_folder := configs_settings.Env.FreelancerFolder

				configs = configs_mapped.NewMappedConfigs()
				logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))
				configs.Read(freelancer_folder)
			}, timeit.WithMsg("read mapping"))
		}, timeit.WithMsg("total time"))
	}
}
