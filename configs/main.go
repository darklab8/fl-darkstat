package configs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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
	// Start profiling
	f, err := os.Create("configs.pprof")
	if err != nil {
		fmt.Println(err)
		return nil

	}
	err = pprof.StartCPUProfile(f)
	logus.Log.CheckError(err, "failed to start cpu profiling")
	start := time.Now()

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
	configs.Mapped.Clean()

	pprof.StopCPUProfile()
	elapsed := time.Since(start)
	log.Printf("Elapsed Pprof time %s", elapsed)

	return configs
}

// from configs. Refactor to integrate it
// go run . configs
// go tool pprof -alloc_space -http=":8001" -nodefraction=0 http://localhost:6060/debug/pprof/heap
// CliConfigs currently contains just debugging run to parse info, useful to profile or integration testing
// This entrypoint will serve for all commands related to configs if there will be more of them
func CliConfigs() {
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

		// for {
		// 	fmt.Println(configs.Bases[0])
		// 	time.Sleep(time.Hour)
		// }
	}
}
