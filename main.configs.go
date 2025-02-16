/*
See package `configs` for description and code examples
*/
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

// *configs_export.Exporter
func GetConfigsExport() *configs_export.Exporter {
	timer_mapping := timeit.NewTimerMain(timeit.WithMsg("read mapping"))
	freelancer_folder := configs_settings.Env.FreelancerFolder
	mapped := configs_mapped.NewMappedConfigs()
	logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))
	mapped.Read(freelancer_folder)
	timer_mapping.Close()

	timer_export := timeit.NewTimerMain(timeit.WithMsg("read mapping"))
	configs := configs_export.NewExporter(mapped)
	configs.Export(configs_export.ExportOptions{})
	timer_export.Close()
	configs.Clean()
	return configs
}

// from configs. Refactor to integrate it
// go run . configs
// go tool pprof -alloc_space -http=":8001" -nodefraction=0 http://localhost:6060/debug/pprof/heap
func main_configs() {
	for i := 0; i < 1; i++ {
		timer_total := timeit.NewTimer("total time")
		var configs *configs_export.Exporter

		func() {
			configs = GetConfigsExport()
		}()

		runtime.GC()
		_ = configs

		fmt.Println("configs are prepared")

		timer_total.Close()

		for {
			fmt.Println(configs.Bases[0])
			time.Sleep(time.Hour)
		}
	}
}
