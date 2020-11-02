package grid

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
	"strings"
)

type Grid struct {
	grid [][]Cell // grid[row][col] // origin in the top left
}

func New(rows, cols int) Grid {
	g := Grid{
		grid: make([][]Cell, rows),
	}
	gridMem := make([]Cell, rows*cols)

	for i := 0; i < rows; i++ {
		g.grid[i] = gridMem[cols*i : cols*(i+1)]
	}

	var idx int
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g.grid[r][c].index = idx
			idx++

			g.grid[r][c].row = r
			g.grid[r][c].col = c

			if r > 0 {
				g.grid[r][c].addNeighbor(&g.grid[r-1][c], NORTH)
			}
			if r+1 < rows {
				g.grid[r][c].addNeighbor(&g.grid[r+1][c], SOUTH)
			}
			if c > 0 {
				g.grid[r][c].addNeighbor(&g.grid[r][c-1], WEST)
			}
			if c+1 < cols {
				g.grid[r][c].addNeighbor(&g.grid[r][c+1], EAST)
			}
		}
	}

	return g
}

func (g Grid) Rows() int {
	return len(g.grid)
}

func (g Grid) Cols() int {
	return len(g.grid[0])
}

func (g Grid) Connect(row, col int, dir Direction) {
	g.grid[row][col].connect(dir)
}

func (g Grid) Disconnect(row, col int, dir Direction) {
	g.grid[row][col].disconnect(dir)
}

func (g Grid) Cell(row, col int) Cell {
	return g.grid[row][col]
}

func (g Grid) CellForIndex(idx int) Cell {
	return g.grid[idx/g.Rows()][idx%g.Rows()]
}

// CellDir returns the direction from a to b
func (g Grid) CellDir(a, b Cell) Direction {
	if a.col == b.col {
		if a.Row() < b.Row() {
			return SOUTH
		}
		return NORTH
	}
	if a.Col() < b.Col() {
		return EAST
	}
	return WEST
}

func (g Grid) String() string {
	var builder strings.Builder
	for r := 0; r < len(g.grid); r++ {
		for i := 0; i < 2; i++ {
			for c := 0; c < len(g.grid[r]); c++ {
				switch i {
				case 0:
					if g.grid[r][c].Connected(NORTH) {
						builder.WriteString("+   ")
					} else {
						builder.WriteString("+---")
					}
				default:
					if g.grid[r][c].Connected(WEST) {
						builder.WriteString("    ")
					} else {
						builder.WriteString("|   ")
					}
				}
			}
			switch i {
			case 0:
				builder.WriteString("+\n")
			default:
				builder.WriteString("|\n")
			}
		}
	}
	for c := 0; c < len(g.grid[0]); c++ {
		builder.WriteString("+---")
	}
	builder.WriteString("+\n")

	return builder.String()
}

func (g Grid) Draw(window pixel.Target, size pixel.Rect, thickness float64) {
	target := imdraw.New(nil)
	target.Color = color.White
	target.Push(pixel.V(size.W()-thickness/2, size.H()+thickness), pixel.V(size.W()-thickness/2, 3*thickness/2), pixel.V(thickness, 3*thickness/2))
	target.Line(thickness)

	cellWidth := (size.W() - thickness) / float64(g.Cols())
	cellHeight := (size.H() - thickness) / float64(g.Rows())
	for r := 0; r < g.Rows(); r++ {
		y := float64(g.Rows()-r)*cellHeight + thickness*2 // top left
		for c := 0; c < g.Cols(); c++ {
			x := float64(c)*cellWidth + thickness // top left
			cell := g.Cell(r, c)

			if !cell.Connected(NORTH) {
				target.Push(pixel.V(x, y), pixel.V(x+cellWidth, y))
				target.Line(thickness)
			}
			if !cell.Connected(WEST) {
				target.Push(pixel.V(x, y), pixel.V(x, y-cellHeight))
				target.Line(thickness)
			}
		}
	}

	target.Draw(window)
}

type MazeStats struct {
	DeadEnds int
	FourWay int
	Corridors int
}

func (g Grid) Statistics() MazeStats {
	var stats MazeStats

	for r := 0; r < g.Rows(); r++ {
		for c := 0; c < g.Rows(); c++ {
			cell := g.Cell(r, c)

			var openings int
			for d := NORTH; d <= WEST; d++ {
				if cell.Connected(d) {
					openings++
				}
			}
			if openings == 1 {
				stats.DeadEnds++
			}
			if openings == 2 {
				stats.Corridors++
			}
			if openings == 4 {
				stats.FourWay++
			}
		}
	}

	return stats
}
