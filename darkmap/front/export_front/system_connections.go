package export_front

import (
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/solararch_mapped"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
	"github.com/darklab8/go-utils/typelog"
)

type SystemGraphInfo struct {
	Reachable bool               // from manhattan as usual
	LeadsTo   map[string]*System // edges to other systems
}

type SystemGraphs struct {
	Systems map[string]*System // sparse graph of systems
}

/*
TODO: Wtf all those errors
time=2026-03-17T02:45:23.472+01:00 level=ERROR msg=" has no system file" system_nick=hlp1
time=2026-03-17T02:45:23.474+01:00 level=ERROR msg=" has no system file" system_nick=sector01
time=2026-03-17T02:45:23.474+01:00 level=ERROR msg=" has no system file" system_nick=sector02
time=2026-03-17T02:45:23.474+01:00 level=ERROR msg=" has no system file" system_nick=sector03
time=2026-03-17T02:45:23.474+01:00 level=ERROR msg=" has no system file" system_nick=sector04
time=2026-03-17T02:45:23.474+01:00 level=ERROR msg=" has no system file" system_nick=unch04b
time=2026-03-17T02:45:23.474+01:00 level=ERROR msg=" has no system file" system_nick=st02c
time=2026-03-17T02:45:23.474+01:00 level=ERROR msg="found system in DFSUtil. Reachable but has no file" err=hi10
*/
func (g *SystemGraphs) DFSUtil(vertex *System, visited map[string]bool) {

	// if _, ok := g.Systems[vertex.Nickname]; !ok {
	// 	logus.Log.Error("found system in DFSUtil. Should be reachable but no in graph", typelog.Any("err", vertex.Nickname))
	// 	return
	// }

	visited[vertex.Nickname] = true
	vertex.Reachable = true

	for _, v := range g.Systems[vertex.Nickname].LeadsTo {
		if !visited[v.Nickname] {
			g.DFSUtil(v, visited)
		}
	}
}

func (g *SystemGraphs) DFS(startVertex *System) {
	visited := make(map[string]bool)
	g.DFSUtil(startVertex, visited)
}

/*
// we have to find only one between each system
// and mark if it is
// - Two Way Jump Gate (Blue Connection)
// - Two way Jump hole (Yellow or White weak connection)
// - One Way (Purple Connection)
// - or Unstable (Orange connection)

// How to find them?
// per system go? and in hashmap... marking data about each system pair ergh?

// We need to mark if the system is reachable from manhattan by any means? to define if it should be filtered
*/
func (e *Export) GetSystemConnections(systems []*System) SystemGraphs {
	var graph SystemGraphs = SystemGraphs{
		Systems: make(map[string]*System),
	}

	for _, system := range systems {
		graph.Systems[system.Nickname] = system
	}

	everything_dockable := solararch_mapped.DockableOptions{
		IsDisco:                  e.Mapped.Discovery != nil,
		PlayersCanDockBerth:      true,
		PlayersCanDockMoorMedium: true,
		PlayersCanDockMoorLarge:  true,
	}
	if e.Mapped.Discovery != nil {
		everything_dockable.WithDiscoFreighterPaths = true
	}

	for _, system := range systems {
		system_info := e.Mapped.Systems.SystemsMap[system.Nickname]

		if system_info == nil {
			logus.Log.Error(" has no system file", typelog.Any("system_nick", system.Nickname))
			continue
		}

		jumpholes := trades.GetDockableJumpholes(
			system_info,
			e.Mapped,
			everything_dockable,
		)

		for _, jh := range jumpholes {
			if _, ok := graph.Systems[system_info.Nickname]; !ok {
				graph.Systems[system_info.Nickname] = system

			}
			system.LeadsTo[jh.GotoSystem.Get()] = graph.Systems[jh.GotoSystem.Get()]
		}

	}

	graph.DFS(graph.Systems["li01"])

	return graph
}
