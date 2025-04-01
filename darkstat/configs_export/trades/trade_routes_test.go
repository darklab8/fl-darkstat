package trades

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/stretchr/testify/assert"
)

func TestTradeRoutes(t *testing.T) {

	configs := configs_mapped.TestFixtureConfigs()

	start := time.Now()

	graph := MapConfigsToFGraph(
		configs,
		DiscoverySpeeds.AvgTransportCruiseSpeed,
		WithFreighterPaths(false),
		make(map[string][]ExtraBase),
		MappingOptions{},
	)

	edges_count := 0
	for _, edges := range graph.matrix {
		edges_count += len(edges)
	}
	fmt.Println("graph.vertixes=", len(graph.matrix), "edges_count=", edges_count)

	// for profiling only stuff.
	f, err := os.Create("dijkstra_apsp.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	timeit.NewTimerF(func() {

		dijkstra := NewDijkstraApspFromGraph(graph)
		dist, parents := dijkstra.DijkstraApsp()

		elapsed := time.Since(start)
		log.Printf("Elapsed diijkstra time %s", elapsed)

		// This version lf algorithm can provide you with distances only originating from space bases (and not proxy bases)
		// The rest of starting points were excluded for performance reasons
		fmt.Println(`GetTimeMs2(graph, dist, "li01_01_base", "li01_to_li02")=`, GetTimeMs2(graph, dist, "li01_01_base", "li01_to_li02"))
		fmt.Println(`GetTimeMs2(graph, dist, "li01_01_base", "li02_to_li01")=`, GetTimeMs2(graph, dist, "li01_01_base", "li02_to_li01"))
		fmt.Println(`GetTimeMs2(graph, dist, "li01_01_base", "li12_02_base")=`, GetTimeMs2(graph, dist, "li01_01_base", "li12_02_base"))
		dist1 := GetTimeMs2(graph, dist, "li01_01_base", "li01_02_base")
		dist2 := GetTimeMs2(graph, dist, "li01_01_base", "br01_01_base")
		dist3 := GetTimeMs2(graph, dist, "li01_01_base", "li12_02_base")
		assert.True(t, dist1 > 0 && dist1 < intgmax/2, "expected distance between bases having reasonable value")
		fmt.Println(`GetTimeMs2(graph, dist, "li01_01_base", "li01_02_base")`, dist1)
		fmt.Println(`GetTimeMs2(graph, dist, "li01_01_base", "br01_01_base")`, dist2)
		fmt.Println(`GetTimeMs2(graph, dist, "li01_01_base", "li12_02_base")`, dist3)
		assert.Greater(t, dist1, Intg(0))
		assert.Greater(t, dist2, Intg(0))
		assert.Greater(t, dist3, Intg(0))

		fmt.Println("detailed_path")
		// paths := graph.GetPaths(parents, dist, "li01_01_base", "br01_01_base")
		// paths := graph.GetPaths(parents, dist, "hi02_01_base", "ga01_02_base")
		// paths := graph.GetPaths(parents, dist, "bw11_02_base", "ew12_02_base")

		source := "hi02_01_base"
		target := "li01_01_base"

		// source := "hi02_01_base"
		// target := "br01_01_base"

		// source := "hi02_01_base"
		// target := "ga01_02_base" // New London Atmosphere/Landing Site

		dist_ := GetTimeMs2(graph, dist, source, target)
		fmt.Println("dist=", dist_)
		fmt.Println("time_total=", graph.GetTimeForDist(float64(dist_)))
		min := math.Floor(float64(graph.GetTimeForDist(float64(dist_))) / 60)
		fmt.Println("time_min=", min)
		fmt.Println("time_sec=", float64(graph.GetTimeForDist(float64(dist_)))-min*60)
		fmt.Println(source, "->", target, " path:")

		paths := graph.GetPaths(parents, dist, source, target)

		for index, path := range paths {
			if path.Dist == 0 && (index != 0 || index != len(paths)-1) {
				continue
			}
			fmt.Println(
				fmt.Sprintf("prev= %20s", path.PrevName),
				fmt.Sprintf("next= %20s", path.NextName),
				fmt.Sprintf("prev_node= %5d", path.PrevNode),
				fmt.Sprintf("next_node= %5d", path.NextNode),
				fmt.Sprintf("prev= %25s", configs.Infocards.Infonames[path.PrevIdsName]),
				fmt.Sprintf("next= %25s", configs.Infocards.Infonames[path.NextIdsName]),
				"dist=", path.Dist,
				"min=", path.TimeMinutes,
				"sec=", path.TimeSeconds,
			)
		}
	}, timeit.WithMsg("trade routes calculated"))
	fmt.Println()
}
