package algorithms

import (
	"github.com/bionoren/mazes/grid"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
	"strconv"
)

type Dijkstra struct {
	reference grid.Cell
	distances []int
	grid      grid.Grid
}

func NewDijkstra(g grid.Grid) Dijkstra {
	return Dijkstra{
		distances: make([]int, g.Rows()*g.Cols()),
		grid:      g,
	}
}

func (d Dijkstra) Reference() grid.Cell {
	return d.reference
}

func (d Dijkstra) Distances() []int {
	return d.distances
}

func (d *Dijkstra) Init(start grid.Cell) {
	for i := range d.distances {
		d.distances[i] = 0
	}
	d.reference = start

	queue := make([]*grid.Cell, 1, len(d.distances))
	visited := make([]bool, len(d.distances)) // we consider a node "visited" when it is *added* to the queue, not when it is actually visited
	visited[start.Index()] = true
	queue[0] = &d.reference

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]
		dist := d.distances[cell.Index()] + 1

		for dir := grid.NORTH; dir <= grid.WEST; dir++ {
			if next := cell.Neighbor(dir); next != nil && cell.Connected(dir) && !visited[next.Index()] {
				d.distances[next.Index()] = dist
				queue = append(queue, next)
				visited[next.Index()] = true
			}
		}
	}
}

// ShortestPath returns the shortest path from the cell this was initialized with to the specified end cell
// the return value is a slice of cell indexes
func (d Dijkstra) ShortestPath(end grid.Cell) []int {
	dist := d.distances[end.Index()]
	path := make([]int, dist+1)
	path[0] = end.Index()

	cell := &end
	path[0] = cell.Index()
	for i := 1; i < len(path); i++ {
		for dir := grid.NORTH; dir <= grid.WEST; dir++ {
			if next := cell.Neighbor(dir); next != nil && cell.Connected(dir) && d.distances[next.Index()] < dist {
				path[i] = next.Index()
				cell = next
				dist = d.distances[cell.Index()]
				break
			}
		}
	}

	return path
}

// ShortestPath returns the longest path from the cell this was initialized with
// the return value is a slice of cell indexes
func (d Dijkstra) LongestPath() []int {
	var maxDistance int
	var maxIndex int
	for i, dist := range d.distances {
		if dist > maxDistance {
			maxIndex = i
			maxDistance = dist
		}
	}

	return d.ShortestPath(d.grid.CellForIndex(maxIndex))
}

func (d Dijkstra) Draw(window pixel.Target, size pixel.Rect, thickness float64, floodFill bool) {
	if floodFill {
		d.Fill(window, size, thickness)
		return
	}

	target := imdraw.New(nil)

	cellWidth := (size.W() - thickness) / float64(d.grid.Cols())
	cellHeight := (size.H() - thickness) / float64(d.grid.Rows())

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	labelWriter := text.New(pixel.ZV, basicAtlas)
	labelWriter.Color = color.White

	for r := 0; r < d.grid.Rows(); r++ {
		y := float64(d.grid.Rows()-r)*cellHeight + thickness*2 // top left
		for c := 0; c < d.grid.Cols(); c++ {
			x := float64(c)*cellWidth + thickness // top left
			cell := d.grid.Cell(r, c)

			labelWriter.Dot = pixel.V(x+thickness, y-cellHeight/2)
			_, _ = labelWriter.WriteString(strconv.Itoa(d.distances[cell.Index()]))
		}
	}

	target.Draw(window)
	labelWriter.Draw(window, pixel.IM)
}

func (d Dijkstra) DrawShortestPath(end grid.Cell, window pixel.Target, size pixel.Rect, thickness float64) {
	d.drawPath(d.ShortestPath(end), window, size, thickness/2, color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	})
}

func (d Dijkstra) DrawLongestPath(window pixel.Target, size pixel.Rect, thickness float64) {
	d.drawPath(d.LongestPath(), window, size, thickness/2, color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 255,
	})
}

func (d Dijkstra) Fill(window pixel.Target, size pixel.Rect, thickness float64) {
	target := imdraw.New(nil)

	cellWidth := (size.W() - thickness) / float64(d.grid.Cols())
	cellHeight := (size.H() - thickness) / float64(d.grid.Rows())

	var maxDistance int
	for _, dist := range d.distances {
		if dist > maxDistance {
			maxDistance = dist
		}
	}

	for r := 0; r < d.grid.Rows(); r++ {
		y := float64(d.grid.Rows()-r)*cellHeight + thickness*2 // top left
		for c := 0; c < d.grid.Cols(); c++ {
			x := float64(c)*cellWidth + thickness // top left
			cell := d.grid.Cell(r, c)

			colorVal := uint8(math.Round(255 * float64(maxDistance-d.distances[cell.Index()]) / float64(maxDistance)))
			target.Color = color.RGBA{
				R: colorVal,
				G: colorVal,
				B: colorVal,
			}

			target.Push(pixel.V(x, y), pixel.V(x+cellWidth, y-cellHeight))
			target.Rectangle(0)
		}
	}

	target.Draw(window)
}

func (d Dijkstra) drawPath(path []int, window pixel.Target, size pixel.Rect, thickness float64, pathColor color.Color) {
	target := imdraw.New(nil)
	target.Color = pathColor

	cellWidth := (size.W() - thickness) / float64(d.grid.Cols())
	cellHeight := (size.H() - thickness) / float64(d.grid.Rows())

	for _, idx := range path {
		cell := d.grid.CellForIndex(idx)
		x := float64(cell.Col())*cellWidth + thickness                  // top left
		y := float64(d.grid.Rows()-cell.Row())*cellHeight + thickness*2 // top left

		target.Push(pixel.V(x+cellWidth/2, y-cellHeight/2))
	}

	target.Line(thickness)
	target.Draw(window)
}
