package front

import (
	"fmt"
	"math"

	"github.com/darklab8/fl-darkstat/darkmap/export_map"
)

type VectorCoords struct {
	x float64
	y float64
}

func Vector(edge *export_map.ConnectionEdge, rotate float64) VectorCoords {
	magnitude := math.Sqrt(math.Pow(*edge.FirstSystem.Pos.X-*edge.SecondSystem.Pos.X, 2) + math.Pow(*edge.FirstSystem.Pos.Y-*edge.SecondSystem.Pos.Y, 2))
	if magnitude == 0 {
		return VectorCoords{
			x: 0,
			y: 0,
		}
	}

	x1 := (*edge.SecondSystem.Pos.X - *edge.FirstSystem.Pos.X) / magnitude
	y1 := (*edge.SecondSystem.Pos.Y - *edge.FirstSystem.Pos.Y) / magnitude

	x_new := x1*math.Cos(rotate*math.Pi/180) - y1*math.Sin(rotate*math.Pi/180)
	y_new := x1*math.Sin(rotate*math.Pi/180) + y1*math.Cos(rotate*math.Pi/180)

	if math.IsNaN(x_new) || math.IsNaN(y_new) {
		fmt.Println(x1, y1, x_new, y_new, rotate)
		panic("NAN in vector calculations")
	}

	return VectorCoords{
		x: (x_new),
		y: (y_new),
	}
}
