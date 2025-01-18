package trades

import (
	"fmt"
	"math"
	"testing"
)

// // Driver Code
func TestDijkstraAPSP(t *testing.T) {
	var vertices int = 4
	var matrix [][]int = [][]int{
		{0, 0, 2, 0},
		{4, 0, 3, 0},
		{0, 0, 0, 2},
		{0, 1, 0, 0},
	}

	// Initialization
	var graph *DijkstraAPSP = NewDijkstraApspFromMatrix(vertices, matrix)

	// Function Call
	distances, parents := graph.DijkstraApsp()
	_ = parents

	// The code fragment below outputs
	// an formatted distance matrix.
	// Its first row and first
	// column represent vertices
	fmt.Println("Distance matrix:")

	fmt.Printf("   \t")
	for i := 0; i < vertices; i++ {
		fmt.Printf("%3d\t", i)
	}

	for i := 0; i < vertices; i++ {
		fmt.Println()
		fmt.Printf("%3d\t", i)
		for j := 0; j < vertices; j++ {
			if distances[i][j] == math.MaxInt {
				fmt.Printf(" X\t")
			} else {
				fmt.Printf("%3d\t",
					distances[i][j])
			}
		}
	}
}

func TestDijkstraAPSPWithGraph(t *testing.T) {
	graph := NewGameGraph(DiscoverySpeeds.AvgTransportCruiseSpeed, WithFreighterPaths(true))
	graph.SetEdge("a", "b", 5)
	graph.SetEdge("a", "d", 10)
	graph.SetEdge("b", "c", 3)
	graph.SetEdge("c", "d", 1)
	johnson := NewDijkstraApspFromGraph(graph)
	dist, parents := johnson.DijkstraApsp()

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if dist[i][j] == math.MaxInt {
				fmt.Printf("%7s", "INF")
			} else {
				fmt.Printf("%7d", dist[i][j])
			}
		}
		fmt.Println()
	}

	fmt.Println("a -> c = ", GetTimeMs2(graph, dist, "a", "c"), "path=", graph.GetPaths(parents, dist, "a", "c"))
	fmt.Println("a -> b = ", GetTimeMs2(graph, dist, "a", "b"), "path=", GetPath(graph, parents, dist, "a", "b"))
	fmt.Println("a -> d = ", GetTimeMs2(graph, dist, "a", "d"), "path=", GetPath(graph, parents, dist, "a", "d"))
}
