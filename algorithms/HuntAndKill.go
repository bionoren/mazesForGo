package algorithms

import (
	"github.com/bionoren/mazes/grid"
	"math/rand"
)

func HuntAndKill(g grid.Grid) {
	size := g.Rows() * g.Cols()
	node := g.CellForIndex(rand.Intn(size))
	visited := make([]bool, size)
	visited[node.Index()] = true
	visits := 1

	directions := make([]grid.Direction, (grid.WEST+1)*2)
	for d := grid.NORTH; d <= grid.WEST; d++ {
		directions[d] = d
	}
	dirOptions := directions[grid.WEST+1:]
	directions = directions[:grid.WEST+1]

	for visits < size {
		if len(dirOptions) > 0 { // kill - keep exploring around
			i := grid.Direction(rand.Intn(len(dirOptions)))
			dir := dirOptions[i]
			if next := node.Neighbor(dir); next != nil && !visited[next.Index()] {
				visited[next.Index()] = true
				visits++
				g.Connect(node.Row(), node.Col(), dir)

				node = *next
				dirOptions = dirOptions[:cap(dirOptions)]
				copy(dirOptions, directions)
			} else { // this random direction is already visited, remove it from the current option list
				dirOptions[i] = dirOptions[len(dirOptions)-1]
				dirOptions = dirOptions[:len(dirOptions)-1]
			}
		} else { // hunt - find a new trailhead
			hunt:
			for r := 0; r < g.Rows(); r++ {
				for c := 0; c < g.Cols(); c++ {
					node = g.Cell(r, c)
					if visited[node.Index()] {
						continue
					}

					for d := grid.NORTH; d <= grid.WEST; d++ {
						if n := node.Neighbor(d); n != nil && visited[n.Index()] {
							dirOptions = append(dirOptions, d)
						}
					}
					if len(dirOptions) != 0 {
						i := grid.Direction(rand.Intn(len(dirOptions)))
						dir := dirOptions[i]
						g.Connect(r, c, dir)
						visited[node.Index()] = true
						visits++

						dirOptions = dirOptions[:cap(dirOptions)]
						copy(dirOptions, directions)

						break hunt
					}
				}
			}
		}
	}
}
