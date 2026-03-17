package export_front

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/solararch_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type SystemGraphInfo struct {
	Reachable bool                       // from manhattan as usual
	LeadsTo   map[string]*JumpConnection // edges to other systems
}

type SystemGraphs struct {
	Systems map[string]*System // sparse graph of systems

}

type JumpConnectionKind int8

const (
	JumpKindUnknown JumpConnectionKind = iota
	JumpKindJumpgate
	JumpKindJumphole
	JumpKindUnstable
	JumpKindAlien
)

type JumpConnection struct {
	Kind JumpConnectionKind
	*System
}

func (g *SystemGraphs) DFSUtil(vertex *System, visited map[string]bool) {

	// if _, ok := g.Systems[vertex.Nickname]; !ok {
	// 	logus.Log.Error("found system in DFSUtil. Should be reachable but no in graph", typelog.Any("err", vertex.Nickname))
	// 	return
	// }

	visited[vertex.Nickname] = true
	vertex.Reachable = true

	for _, v := range g.Systems[vertex.Nickname].LeadsTo {
		if !visited[v.Nickname] {
			g.DFSUtil(v.System, visited)
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
// - Two way Jump hole (Yellow connection)
// - One Way (Purple Connection)
// - Two Way Unstable (Orange connection)
// - Pink, unidentified

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

			target_system := graph.Systems[jh.GotoSystem.Get()]
			system.LeadsTo[jh.GotoSystem.Get()] = &JumpConnection{
				System: target_system,
				Kind:   e.GetJumpConnectionKind(jh),
			}
		}

	}

	graph.DFS(graph.Systems["li01"])

	return graph
}

func (e *Export) GetJumpConnectionKind(jh *systems_mapped.Jumphole) JumpConnectionKind {
	jh_archetype := jh.Archetype.Get()

	var disco_cargo_limit *int
	if solar, ok := e.Mapped.Solararch.SolarsByNick[jh_archetype]; ok {
		if cargo_limit, ok := solar.CargoLimit.GetValue(); ok {
			disco_cargo_limit = ptr.Ptr(cargo_limit)
		}
	}
	if e.Mapped.Discovery != nil {
		if disco_cargo_limit != nil {
			if *disco_cargo_limit < trades.DiscoCargoLimitedThreshold {
				return JumpKindUnstable
			}
		}
	}

	if strings.Contains(jh_archetype, "jumpgate") {
		return JumpKindJumpgate
	}

	if strings.Contains(jh_archetype, "jumphole") {
		return JumpKindJumphole
	}

	if strings.Contains(jh_archetype, "nomad") {
		return JumpKindAlien
	}

	return JumpKindUnknown
}
