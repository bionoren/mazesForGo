package algorithms

import (
	"github.com/bionoren/mazes/grid"
	"math/rand"
)

func Sidewinder(g grid.Grid) {
	var runStart int
	for r := 0; r < g.Rows(); r++ {
		for c := 0; c < g.Cols(); c++ {
			cell := g.Cell(r, c)
			switch rand.Intn(2) {
			case 0: // continue the run
				if cell.HasNeighbor(grid.EAST) {
					g.Connect(r, c, grid.EAST)
				}
			case 1:
				if cell.HasNeighbor(grid.NORTH) {
					idx := rand.Intn(c-runStart+1) + runStart
					g.Connect(r, idx, grid.NORTH)
				} else if cell.HasNeighbor(grid.EAST) {
					g.Connect(r, c, grid.EAST)
				}
				runStart = c + 1
			}
		}
		if runStart < g.Cols() && g.Cell(r, 0).HasNeighbor(grid.NORTH) {
			idx := rand.Intn(g.Cols()-runStart) + runStart
			g.Connect(r, idx, grid.NORTH)
		}
		runStart = 0
	}
}
