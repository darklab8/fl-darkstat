package export_map

import (
	"fmt"
	"sort"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/solararch_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type SystemGraphInfo struct {
	VisibleByDefault bool                       // from manhattan as usual
	LeadsTo          map[string]*JumpConnection // edges to other systems
}

type SystemGraphs struct {
	Systems         map[string]*System         // sparse graph of systems
	ConnectionEdges map[string]*ConnectionEdge // ready for printing. Only single one per each two systems.
}

type ConnectionEdge struct {
	FirstSystem  *System
	SecondSystem *System

	FromFirstToSecondJumpable bool
	FromSecondToFirstJumpable bool

	// Const is made in values ordered by priority
	// if Alien gate is present, then it is most important
	// after that Jumpgate, then Jumphole, then Unstable, and then Unknown
	Kind JumpConnectionKind
}

func (c ConnectionEdge) IsBireDirectional() bool {
	return c.FromFirstToSecondJumpable && c.FromSecondToFirstJumpable
}

func (c ConnectionEdge) GetPosX() float64 {
	return (*c.FirstSystem.Pos.X + *c.SecondSystem.Pos.X) / 2
}
func (c ConnectionEdge) GetPosY() float64 {
	return (*c.FirstSystem.Pos.Y + *c.SecondSystem.Pos.Y) / 2
}

func (c *ConnectionEdge) SetJumpable(FromSystemNick string, ToSystemNick string) {
	if c.FirstSystem.Nickname == FromSystemNick && c.SecondSystem.Nickname == ToSystemNick {
		c.FromFirstToSecondJumpable = true
	} else if c.SecondSystem.Nickname == FromSystemNick && c.FirstSystem.Nickname == ToSystemNick {
		c.FromSecondToFirstJumpable = true
	} else {
		logus.Log.Panic(
			"Received unexpected SetJumptable values for ConnectionEdge",
			typelog.Any("FromSystem", FromSystemNick),
			typelog.Any("ToSystem", ToSystemNick),
			typelog.Any("c.FirstSystem", c.FirstSystem),
			typelog.Any("c.SecondSystem", c.SecondSystem),
		)
	}
}

func (c *ConnectionEdge) SetKind(first map[JumpConnectionKind]bool, second map[JumpConnectionKind]bool) {
	c.Kind = JumpKindUnknown

	if _, ok := first[JumpKindAlien]; ok {
		if _, ok := second[JumpKindAlien]; ok {
			c.Kind = JumpKindAlien
			return
		}
	}
	if _, ok := first[JumpKindJumpgate]; ok {
		if _, ok := second[JumpKindJumpgate]; ok {
			c.Kind = JumpKindJumpgate
			return
		}
	}
	if _, ok := first[JumpKindJumphole]; ok {
		if _, ok := second[JumpKindJumphole]; ok {
			c.Kind = JumpKindJumphole
			return
		}
	}
	if _, ok := first[JumpKindUnstableCargoLimited]; ok {
		if _, ok := second[JumpKindUnstableCargoLimited]; ok {
			c.Kind = JumpKindUnstableCargoLimited
			return
		}
	}

	if _, ok := first[JumpKindFighterOnly]; ok {
		if _, ok := second[JumpKindFighterOnly]; ok {
			c.Kind = JumpKindFighterOnly
			return
		}
	}
	if _, ok := first[JumpKindCampaignLocked]; ok {
		if _, ok := second[JumpKindCampaignLocked]; ok {
			c.Kind = JumpKindCampaignLocked
			return
		}
	}
}

func GetConnKey(from_system string, to_system string) string {
	systems := []string{from_system, to_system}
	sort.Strings(systems)
	return strings.Join(systems, "-")
}

func NewConnectionEdge(first_system *System, second_system *System) *ConnectionEdge {
	systems := []*System{first_system, second_system}

	sort.Slice(systems, func(i, j int) bool { // this sorting is not necessary technically. just nice to do for more deterministic debug
		return systems[i].Nickname > systems[j].Nickname
	})
	e := &ConnectionEdge{
		FirstSystem:  systems[0],
		SecondSystem: systems[1],
	}

	return e
}

type JumpConnectionKind int8

const (
	JumpKindUnknown JumpConnectionKind = iota
	JumpKindUnstableCargoLimited
	JumpKindFighterOnly
	JumpKindJumphole
	JumpKindJumpgate
	JumpKindAlien
	JumpKindCampaignLocked
)

func (j JumpConnectionKind) ToStr() string {
	switch j {
	case JumpKindAlien:
		return "jh_alien"
	case JumpKindJumpgate:
		return "jh_jumpgate"
	case JumpKindJumphole:
		return "jh_jumphole"
	case JumpKindUnstableCargoLimited:
		return "jh_unstable_cargo_limit"
	case JumpKindFighterOnly:
		return "jh_fighter_only"
	case JumpKindCampaignLocked:
		return "jh_campaign_locked"
	}
	return "jh_unknown"
}

type JumpConnection struct {
	Kind    map[JumpConnectionKind]bool
	*System `json:"-" swaggerignore:"true"`
}

func (g *SystemGraphs) DFSUtil(vertex *System, visited map[string]bool) {

	// if _, ok := g.Systems[vertex.Nickname]; !ok {
	// 	logus.Log.Error("found system in DFSUtil. Should be reachable but no in graph", typelog.Any("err", vertex.Nickname))
	// 	return
	// }

	visited[vertex.Nickname] = true
	vertex.VisibleByDefault = true

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

	everything_dockable := solararch_mapped.DockableOptions{
		IsDisco:                  e.Mapped.Discovery != nil,
		PlayersCanDockBerth:      true,
		PlayersCanDockMoorMedium: true,
		PlayersCanDockMoorLarge:  true,
	}
	if e.Mapped.Discovery != nil {
		everything_dockable.WithDiscoFreighterPaths = true
	}
	graph_reachable := e.BuildGraph(systems, everything_dockable)

	graph_reachable.DFS(graph_reachable.Systems["li01"])

	if e.Mapped.Discovery != nil {
		graph_reachable.DFS(graph_reachable.Systems["ew12"])
		for _, system := range graph_reachable.Systems {
			if strings.Contains(strings.ToLower(system.Name), "planet") {
				system.VisibleByDefault = false
			} else if strings.Contains(strings.ToLower(system.Name), "anomaly") {
				system.VisibleByDefault = false
			} else if strings.Contains(strings.ToLower(system.Name), "atmosphere") {
				system.VisibleByDefault = false
			}
		}
	}

	everything_dockable.ShowInitialWorldBlocked = true
	graph := e.BuildGraph(systems, everything_dockable)

	// preparing final edges for front render
	for origin_system_nick, origin_system := range graph.Systems {
		for _, target_conn := range origin_system.LeadsTo {
			key := GetConnKey(origin_system_nick, target_conn.Nickname)

			connection, ok := graph.ConnectionEdges[key]
			if !ok {
				connection = NewConnectionEdge(origin_system, target_conn.System)
				graph.ConnectionEdges[key] = connection
			}

			target_kinds := make(map[JumpConnectionKind]bool)
			if _, ok := graph.Systems[target_conn.Nickname].LeadsTo[origin_system_nick]; ok {
				target_kinds = graph.Systems[target_conn.Nickname].LeadsTo[origin_system_nick].Kind
			}

			connection.SetKind(target_conn.Kind, target_kinds)
			connection.SetJumpable(origin_system_nick, target_conn.Nickname)
		}
	}

	// Printer debugger of unstable connections
	unique_msgs := make(map[string]bool)
	for origin_system_nick, origin_system := range graph.Systems {
		for _, target_conn := range origin_system.LeadsTo {

			target_system := &JumpConnection{System: &System{}}
			if value, ok := graph.Systems[target_conn.Nickname].LeadsTo[origin_system_nick]; ok {
				target_system = value
			}

			if !origin_system.VisibleByDefault &&
				!target_system.VisibleByDefault {
				continue
			}

			key := GetConnKey(origin_system_nick, target_conn.Nickname)
			if _, ok := unique_msgs[key]; ok {
				continue
			}
			unique_msgs[key] = true

			connection, ok := graph.ConnectionEdges[key]
			if !ok {
				connection = NewConnectionEdge(origin_system, target_conn.System)
				graph.ConnectionEdges[key] = connection
			}

			target_kinds := make(map[JumpConnectionKind]bool)
			if _, ok := graph.Systems[target_conn.Nickname].LeadsTo[origin_system_nick]; ok {
				target_kinds = graph.Systems[target_conn.Nickname].LeadsTo[origin_system_nick].Kind
			}

			if connection.Kind == JumpKindUnknown {
				logus.Log.Warn("asymetric connection",
					typelog.Any("sys1-2", fmt.Sprintln(origin_system_nick, target_conn.Nickname)),
					typelog.Any("target_name", target_conn.Name),
					typelog.Any("origin_name", origin_system.Name),
					typelog.Any("target_conns", PrintFormatedConns(target_conn.Kind)),
					typelog.Any("origin_conns", PrintFormatedConns(target_kinds)))

			}
		}
	}

	return graph
}

func (e *Export) BuildGraph(systems []*System, everything_dockable solararch_mapped.DockableOptions) (graph SystemGraphs) {
	graph = SystemGraphs{
		Systems:         make(map[string]*System),
		ConnectionEdges: make(map[string]*ConnectionEdge),
	}

	for _, system := range systems {
		graph.Systems[system.Nickname] = system
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

			if conn, ok := system.LeadsTo[jh.GotoSystem.Get()]; ok {
				conn.Kind[e.GetJumpConnectionKind(jh)] = true
			} else {
				cn := &JumpConnection{
					System: target_system,
					Kind:   make(map[JumpConnectionKind]bool),
				}
				system.LeadsTo[jh.GotoSystem.Get()] = cn
				cn.Kind[e.GetJumpConnectionKind(jh)] = true
			}

		}

	}
	return graph
}

func PrintFormatedConns(target_kinds map[JumpConnectionKind]bool) map[string]bool {
	result := make(map[string]bool)

	for key, _ := range target_kinds {
		result[key.ToStr()] = true
	}
	return result
}

func (e *Export) GetJumpConnectionKind(jh *systems_mapped.Jumphole) JumpConnectionKind {

	// Check locked_gate if it is enterable.
	hash_id := flhash.HashNickname(jh.Nickname.Get())
	if _, ok := e.Mapped.InitialWorld.LockedGates[hash_id]; ok {
		return JumpKindCampaignLocked
	}

	jh_archetype := jh.Archetype.Get()

	var disco_cargo_limit *int
	if solar, ok := e.Mapped.Solararch.SolarsByNick[jh_archetype]; ok {
		if cargo_limit, ok := solar.CargoLimit.GetValue(); ok {
			disco_cargo_limit = ptr.Ptr(cargo_limit)
		}

		if len(solar.DockingSpheres) == 0 {
			return JumpKindUnknown
		}

		is_fighter_only := true
		has_jump_sphere := false
		for _, sphere := range solar.DockingSpheres {
			if sphere.Get() != "berth" {
				is_fighter_only = false
			}

			if sphere.Get() == "jump" {
				has_jump_sphere = true
			}
		}
		if is_fighter_only {
			return JumpKindFighterOnly
		}

		if !has_jump_sphere {
			return JumpKindUnknown
		}
	}
	if disco_cargo_limit != nil {
		if *disco_cargo_limit < trades.DiscoCargoLimitedThreshold {
			return JumpKindUnstableCargoLimited
		}
	}

	if strings.Contains(jh_archetype, "nomad") {
		return JumpKindAlien
	} else if strings.Contains(jh_archetype, "jumpgate") {
		return JumpKindJumpgate
	} else if strings.Contains(jh_archetype, "jumphole") {
		return JumpKindJumphole
	}

	return JumpKindUnknown
}
