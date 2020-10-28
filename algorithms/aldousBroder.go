package algorithms

import (
	"github.com/bionoren/mazes/grid"
	"math/rand"
)

func AldousBroder(g grid.Grid) {
	size := g.Rows() * g.Cols()
	node := g.CellForIndex(rand.Intn(size))
	visited := make([]bool, size)
	visited[node.Index()] = true
	visits := 1

	for visits < size {
		dir := grid.Direction(rand.Intn(int(grid.WEST) + 1))
		if node.HasNeighbor(dir) {
			if !visited[node.Neighbor(dir).Index()] {
				visited[node.Neighbor(dir).Index()] = true
				visits++
				g.Connect(node.Row(), node.Col(), dir)
			}
			node = *node.Neighbor(dir)
		}
	}
}
