package algorithms

import (
	"github.com/bionoren/mazes/grid"
	"math/rand"
)

func RecursiveBacktracker(g grid.Grid) {
	size := g.Rows() * g.Cols()
	node := g.CellForIndex(rand.Intn(size))
	visited := make([]bool, size)
	visited[node.Index()] = true

	directions := make([]grid.Direction, (grid.WEST+1)*2)
	for d := grid.NORTH; d <= grid.WEST; d++ {
		directions[d] = d
	}
	dirOptions := directions[grid.WEST+1:]
	directions = directions[:grid.WEST+1]

	stack := make([]grid.Cell, 1, g.Cols()*g.Rows()/2)
	stack[0] = node
	for len(stack) > 0 {
		// fmt.Printf("stack size: %d, dirOptions size: %d\n", len(stack), len(dirOptions))
		if len(dirOptions) > 0 {
			i := grid.Direction(rand.Intn(len(dirOptions)))
			dir := dirOptions[i]
			if next := node.Neighbor(dir); next != nil && !visited[next.Index()] {
				visited[next.Index()] = true
				g.Connect(node.Row(), node.Col(), dir)

				node = *next
				dirOptions = dirOptions[:cap(dirOptions)]
				copy(dirOptions, directions)
				stack = append(stack, node)
			} else { // this random direction is already visited, remove it from the current option list
				dirOptions[i] = dirOptions[len(dirOptions)-1]
				dirOptions = dirOptions[:len(dirOptions)-1]
			}
		} else { // backtrack
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			dirOptions = dirOptions[:cap(dirOptions)]
			copy(dirOptions, directions)
		}
	}
}
