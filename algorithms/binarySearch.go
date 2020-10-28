package algorithms

import (
	"github.com/bionoren/mazes/grid"
	"math/rand"
)

func BinarySearch(g grid.Grid) {
	for r := 0; r < g.Rows(); r++ {
		for c := 0; c < g.Cols(); c++ {
			cell := g.Cell(r, c)
			switch rand.Intn(2) {
			case 0:
				if cell.HasNeighbor(grid.WEST) {
					g.Connect(r, c, grid.WEST)
				} else if cell.HasNeighbor(grid.SOUTH) {
					g.Connect(r, c, grid.SOUTH)
				}
			case 1:
				if cell.HasNeighbor(grid.SOUTH) {
					g.Connect(r, c, grid.SOUTH)
				} else if cell.HasNeighbor(grid.WEST) {
					g.Connect(r, c, grid.WEST)
				}
			}
		}
	}
}
