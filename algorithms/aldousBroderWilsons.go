package algorithms

import (
	"github.com/bionoren/mazes/grid"
	"math/rand"
)

// AldousBroderWilsons runs AldousBroder until either the grid is half visited or it has run for size*4 iterations. Then it runs Wilson's algorithm until the grid is fully visited
func AldousBroderWilsons(g grid.Grid) {
	// Aldous-broder
	size := g.Rows() * g.Cols()
	node := g.CellForIndex(rand.Intn(size))
	visited := make([]bool, size)
	visits := 1

	visited[node.Index()] = true
	var i = 0
	for visits < size/2 && i < size*4 {
		dir := grid.Direction(rand.Intn(int(grid.WEST) + 1))
		if node.HasNeighbor(dir) {
			if !visited[node.Neighbor(dir).Index()] {
				visited[node.Neighbor(dir).Index()] = true
				visits++
				g.Connect(node.Row(), node.Col(), dir)
			}
			node = *node.Neighbor(dir)
		}
		i++
	}

	// Wilson's
	options := make([]int, 0, size)
	for i, v := range visited {
		if !v {
			options = append(options, i)
		}
	}

	if len(options) == 0 {
		return
	}
	node = g.CellForIndex(options[rand.Intn(len(options))])

	pathed := make([]bool, size)
	path := make([]grid.Cell, 0, size/2) // on large grids, size/2 is a reasonable initial memory guess
	pathed[node.Index()] = true
	for {
		dir := grid.Direction(rand.Intn(int(grid.WEST) + 1))
		if node.HasNeighbor(dir) {
			next := *node.Neighbor(dir)

			if visited[next.Index()] { // connect the path to the visited maze
				for i := len(path) - 1; i >= 0; i-- {
					g.Connect(next.Row(), next.Col(), g.CellDir(next, path[i]))
					next = path[i]
					visited[next.Index()] = true
				}

				filteredOpts := options[:0]
				for _, idx := range options {
					if !visited[idx] {
						filteredOpts = append(filteredOpts, idx)
					}
				}
				options = filteredOpts

				if len(options) == 0 {
					break
				}
				next = g.CellForIndex(options[rand.Intn(len(options))])
				path = path[:1]
				path[0] = next
				pathed[next.Index()] = true
			} else if pathed[next.Index()] { // loop-erase
				for i := len(path) - 1; i >= 0; i-- {
					n := path[i]
					pathed[n.Index()] = false
					path = path[:i]
					if n.Index() == next.Index() {
						break
					}
				}

				pathed[next.Index()] = true
				path = append(path, next)
			} else { // build path
				pathed[next.Index()] = true
				path = append(path, next)
			}
			node = next
		}
	}
}
