package trades

import (
	"container/heap"
	"math"
)

// Written based on https://www.geeksforgeeks.org/implementation-of-johnsons-algorithm-for-all-pairs-shortest-paths/
// And then rewritten in a way that there is nothing left from Johnson almost.
// It is now DijkstraAPSP https://en.wikipedia.org/wiki/Parallel_all-pairs_shortest_path_algorithm

type Neighbour struct {
	destination int
	weight      int
}

func NewNeighbour(destination int, weight int) *Neighbour {
	return &Neighbour{destination: destination,
		weight: weight,
	}
}

type DijkstraAPSP struct {
	vertices         int
	adjacencyList    [][]*Neighbour
	allowed_base_ids map[int]bool
}

const INF = math.MaxInt
const INFthreshold = math.MaxInt / 1000

// On using the below constructor,
// edges must be added manually
// to the graph using addEdge()
func NewDijkstraApsp(vertices int) *DijkstraAPSP {
	g := &DijkstraAPSP{
		vertices:         vertices,
		allowed_base_ids: make(map[int]bool),
	}

	g.adjacencyList = make([][]*Neighbour, vertices)
	for i := 0; i < vertices; i++ {
		g.adjacencyList[i] = make([]*Neighbour, 0)
	}

	return g
}

// // On using the below constructor,
// // edges will be added automatically
// // to the graph using the adjacency matrix
func NewDijkstraApspFromMatrix(vertices int, adjacencyMatrix [][]int) *DijkstraAPSP {
	g := NewDijkstraApsp(vertices)

	for i := 0; i < vertices; i++ {
		for j := 0; j < vertices; j++ {
			if adjacencyMatrix[i][j] != 0 {
				g.addEdge(i, j, adjacencyMatrix[i][j])
			}
		}
	}
	return g
}

type DijkstraOption func(graph *DijkstraAPSP)

func WithPathDistsForAllNodes() DijkstraOption {
	return func(graph *DijkstraAPSP) {
		graph.allowed_base_ids = make(map[int]bool)
	}
}

func NewDijkstraApspFromGraph(graph *GameGraph, opts ...DijkstraOption) *DijkstraAPSP {
	vertices := len(graph.matrix)
	g := NewDijkstraApsp(vertices)

	index := 0
	for vertex, _ := range graph.matrix {
		graph.IndexByNick[vertex] = index
		graph.NicknameByIndex[index] = vertex
		index++
	}

	// TODO make path reconstructions working with it
	for base_nick, _ := range graph.AllowedVertixesForCalcs {
		g.allowed_base_ids[graph.IndexByNick[base_nick]] = true
	}

	for vertex_name, vertex := range graph.matrix {
		for vertex_target, weight := range vertex {
			i := graph.IndexByNick[vertex_name]
			j := graph.IndexByNick[vertex_target]

			g.addEdge(i, j, int(weight))
		}
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func (g *DijkstraAPSP) addEdge(source int, destination int, weight int) {
	g.adjacencyList[source] = append(g.adjacencyList[source], NewNeighbour(destination, weight))
}

func ArraysFill[T any](array []T, value T) {
	for i := 0; i < len(array); i++ {
		array[i] = value
	}
}

type Parent struct {
	node   int
	weight int
}

// // Time complexity of this
// // implementation of dijkstra is O(V^2).
func (g *DijkstraAPSP) dijkstra(source int) ([]int, []Parent) {
	var distance []int = make([]int, g.vertices)

	// this page https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
	// helped to modify algorithm so it would return reconstructed paths.
	// Parent array to store shortest
	// path tree
	var parents []Parent = make([]Parent, g.vertices)
	// The starting vertex does not
	// have a parent
	for s := 0; s < g.vertices; s++ {
		parents[s].node = NO_PARENT
	}

	pq := make(PriorityQueue, 0)
	item := &Item{
		value_weight: 0,
		priority:     source,
	}
	pq.Push(item)

	ArraysFill(distance, INF)
	distance[source] = 0

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		node := item.priority
		dist := item.value_weight

		for _, neighbour := range g.adjacencyList[node] {
			if dist+neighbour.weight < distance[neighbour.destination] {
				parents[neighbour.destination].node = node
				parents[neighbour.destination].weight = neighbour.weight
				distance[neighbour.destination] = dist + neighbour.weight
				pq.Push(&Item{
					value_weight: distance[neighbour.destination],
					priority:     neighbour.destination,
				})
			}

		}

	}

	// cleanup a bit
	for pq.Len() > 0 {
		pq.Pop()
	}

	return distance, parents
}

type DijkstraResult struct {
	source         int
	dist_result    []int
	parents_result []Parent
}

const NO_PARENT = -1

func (g *DijkstraAPSP) DijkstraApsp() ([][]int, [][]Parent) {
	var distances [][]int = make([][]int, g.vertices)
	var parents [][]Parent = make([][]Parent, g.vertices)

	// Performance optimization of the algorithm
	// By skipping heaviest calculations for all shortest paths
	// originiating from vertexes not needed.
	// As those vertex are important only as intermediate travel point.
	skip_not_allowed_vertex := func(source int) ([]int, bool) {
		if len(g.allowed_base_ids) > 0 {
			_, is_base := g.allowed_base_ids[source]
			if !is_base {
				dist := make([]int, g.vertices)
				ArraysFill(dist, INF)
				dist[source] = 0
				return dist, true
			}
		}
		return nil, false
	}

	// it is nice to keep sanity by keeping optional switch removing parallelism
	is_sequential := false

	if is_sequential {
		for s := 0; s < g.vertices; s++ {
			if dist_result, is_skipped := skip_not_allowed_vertex(s); is_skipped {
				distances[s] = dist_result
				continue
			}
			dist_result, parents_result := g.dijkstra(s)
			distances[s] = dist_result
			parents[s] = parents_result
		}
	} else {
		dijkstra_results := make(chan *DijkstraResult)
		awaited := 0
		for s := 0; s < g.vertices; s++ {
			if dist_result, is_skipped := skip_not_allowed_vertex(s); is_skipped {
				distances[s] = dist_result
				continue
			}
			awaited += 1
			go func(s int) {
				dist_result, parents_result := g.dijkstra(s)
				dijkstra_results <- &DijkstraResult{
					source:         s,
					dist_result:    dist_result,
					parents_result: parents_result,
				}
			}(s)
		}
		for s := 0; s < awaited; s++ {
			result := <-dijkstra_results
			distances[result.source] = result.dist_result
			parents[result.source] = result.parents_result
		}
	}

	return distances, parents
}
