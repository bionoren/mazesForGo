package main

import (
	"flag"
	"fmt"
	"github.com/bionoren/mazes/algorithms"
	"github.com/bionoren/mazes/grid"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type menuSettings struct {
	algorithm    int
	solve        bool
	longestPath  bool
	floodFill    bool
	showDijkstra bool

	start *grid.Cell
	end   *grid.Cell
}

func (s *menuSettings) GridReset() {
	s.end = nil
	s.solve = false
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Mazes",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		frames int64
		second = time.Tick(time.Second)
	)

	g := grid.New(16, 16)
	dj := algorithms.NewDijkstra(g)

	minWinDim := math.Min(win.Bounds().W(), win.Bounds().H()) - thickness
	maxWinDim := math.Max(win.Bounds().W(), win.Bounds().H()) - thickness
	graphRect := pixel.R(thickness, thickness, minWinDim, minWinDim)

	cellWidth := (graphRect.W() - thickness) / float64(g.Cols())
	cellHeight := (graphRect.H() - thickness) / float64(g.Rows())

	algorithm := algorithms.BinarySearch
	algorithm(g)

	var settings menuSettings
	settings.algorithm = 1

	DrawMenu(win, pixel.R(minWinDim, minWinDim, maxWinDim, thickness), settings)

	repaint := true
	var regrid bool
	var key pixelgl.Button
	for !win.Closed() {
		win.Update()

		if win.JustPressed(pixelgl.Key1) {
			algorithm = algorithms.BinarySearch
			settings.algorithm = 1

			regrid = true
			repaint = true
		}
		if win.JustPressed(pixelgl.Key2) {
			algorithm = algorithms.Sidewinder
			settings.algorithm = 2

			regrid = true
			repaint = true
		}
		if win.JustPressed(pixelgl.Key3) {
			algorithm = algorithms.AldousBroder
			settings.algorithm = 3

			regrid = true
			repaint = true
		}
		if win.JustPressed(pixelgl.Key4) {
			algorithm = algorithms.Wilsons
			settings.algorithm = 4

			regrid = true
			repaint = true
		}
		if win.JustPressed(pixelgl.Key5) {
			algorithm = algorithms.AldousBroderWilsons
			settings.algorithm = 5

			regrid = true
			repaint = true
		}
		if win.JustPressed(pixelgl.Key6) {
			algorithm = algorithms.HuntAndKill
			settings.algorithm = 6

			regrid = true
			repaint = true
		}
		if win.JustPressed(pixelgl.Key7) {
			algorithm = algorithms.RecursiveBacktracker
			settings.algorithm = 7

			regrid = true
			repaint = true
		}

		if win.JustPressed(pixelgl.KeyN) {
			regrid = true
			repaint = true
		}
		if win.JustPressed(pixelgl.KeyD) && settings.start != nil {
			dj.Init(*settings.start)
			settings.showDijkstra = !settings.showDijkstra

			repaint = true
		}
		if win.JustPressed(pixelgl.KeyS) {
			key = pixelgl.KeyS
		}
		if win.JustPressed(pixelgl.KeyE) {
			key = pixelgl.KeyE
		}
		if win.JustPressed(pixelgl.KeyF) {
			settings.floodFill = !settings.floodFill
			repaint = true
		}
		if win.JustPressed(pixelgl.KeyP) {
			settings.solve = !settings.solve
			repaint = true
		}
		if win.JustPressed(pixelgl.KeyL) {
			settings.longestPath = !settings.longestPath
			repaint = true
		}

		if regrid {
			g = grid.New(g.Rows(), g.Cols())
			algorithm(g)
			dj = algorithms.NewDijkstra(g)
			settings.GridReset()
			if settings.showDijkstra && settings.start != nil {
				cell := g.Cell(settings.start.Row(), settings.start.Col())
				settings.start = &cell
				dj.Init(*settings.start)
			}

			regrid = false
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) && key != 0 {
			mousev := win.MousePosition()

			r := int(math.Floor(float64(g.Rows()) - (mousev.Y-thickness*2)/cellHeight))
			c := int(math.Floor((mousev.X - thickness) / cellWidth))

			if r >= 0 && r < g.Rows() && c >= 0 && c < g.Rows() {
				switch key {
				case pixelgl.KeyS:
					cell := g.Cell(r, c)
					settings.start = &cell

					if settings.showDijkstra {
						dj.Init(*settings.start)
					}
				case pixelgl.KeyE:
					cell := g.Cell(r, c)
					settings.end = &cell
				}
			}

			key = 0
			repaint = true
		}

		if repaint {
			win.Clear(color.Black)
			DrawMenu(win, pixel.R(minWinDim, minWinDim, maxWinDim, thickness), settings)
			if settings.showDijkstra {
				dj.Draw(win, graphRect, thickness, settings.floodFill)
			}

			if settings.solve && settings.end != nil {
				dj.DrawShortestPath(*settings.end, win, graphRect, thickness)
			}
			if settings.longestPath {
				dj.DrawLongestPath(win, graphRect, thickness)
			}

			DrawStartEnd(g, win, graphRect, settings)
			g.Draw(win, graphRect, thickness)

			repaint = false
		}

		frames++
		select {
		case <-second:
			win.SetTitle("Mazes | FPS: " + strconv.FormatInt(frames, 10))
			frames = 0
		default:
		}
	}
}

func DrawStartEnd(g grid.Grid, target pixel.Target, bounds pixel.Rect, settings menuSettings) {
	cellWidth := (bounds.W() - thickness) / float64(g.Cols())
	cellHeight := (bounds.H() - thickness) / float64(g.Rows())

	draw := imdraw.New(nil)

	if settings.start != nil {
		x := float64(settings.start.Col())*cellWidth + thickness             // top left
		y := float64(g.Rows()-settings.start.Row())*cellHeight + thickness*2 // top left

		draw.Color = color.RGBA{
			R: 0,
			G: 200,
			B: 0,
			A: 255,
		}

		draw.Push(pixel.V(x, y), pixel.V(x+cellWidth, y-cellHeight))
		draw.Rectangle(0)
	}

	if settings.end != nil {
		x := float64(settings.end.Col())*cellWidth + thickness             // top left
		y := float64(g.Rows()-settings.end.Row())*cellHeight + thickness*2 // top left

		draw.Color = color.RGBA{
			R: 0,
			G: 200,
			B: 200,
			A: 255,
		}

		draw.Push(pixel.V(x, y), pixel.V(x+cellWidth, y-cellHeight))
		draw.Rectangle(0)
	}

	draw.Draw(target)
}

func DrawMenu(target pixel.Target, bounds pixel.Rect, settings menuSettings) {
	draw := imdraw.New(nil)
	draw.Color = color.RGBA{
		R: 0,
		G: 0,
		B: 50,
	}
	green := color.RGBA{
		R: 0,
		G: 200,
		B: 0,
	}

	draw.Push(bounds.Vertices()[0], bounds.Vertices()[1], bounds.Vertices()[2], bounds.Vertices()[3], bounds.Vertices()[0])
	draw.Rectangle(thickness)
	draw.Draw(target)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	labelWriter := text.New(bounds.Vertices()[0].Add(pixel.V(thickness, -thickness*3)), basicAtlas)
	labelWriter.Color = color.White

	labelWriter.WriteString("algorithms:\n")
	if settings.algorithm == 1 {
		labelWriter.Color = green
	} else {
		labelWriter.Color = color.White
	}
	labelWriter.WriteString("1 - binary search\n")
	if settings.algorithm == 2 {
		labelWriter.Color = green
	} else {
		labelWriter.Color = color.White
	}
	labelWriter.WriteString("2 - sidewinder\n")
	if settings.algorithm == 3 {
		labelWriter.Color = green
	} else {
		labelWriter.Color = color.White
	}
	labelWriter.WriteString("3 - Aldous-Broder\n")
	if settings.algorithm == 4 {
		labelWriter.Color = green
	} else {
		labelWriter.Color = color.White
	}
	labelWriter.WriteString("4 - Wilson's\n")
	if settings.algorithm == 5 {
		labelWriter.Color = green
	} else {
		labelWriter.Color = color.White
	}
	labelWriter.WriteString("5 - Aldous-Broder-Wilson's\n")
	if settings.algorithm == 6 {
		labelWriter.Color = green
	} else {
		labelWriter.Color = color.White
	}
	labelWriter.WriteString("6 - Hunt and Kill\n")
	if settings.algorithm == 7 {
		labelWriter.Color = green
	} else {
		labelWriter.Color = color.White
	}
	labelWriter.WriteString("7 - Recursive Backtracker\n")
	labelWriter.Color = color.White
	labelWriter.WriteRune('\n')
	labelWriter.WriteString("commands:\n")
	labelWriter.WriteString("n - generate new maze\n")
	labelWriter.WriteString("s - mark maze start with click\n")
	labelWriter.WriteString("e - mark maze end with click\n")
	labelWriter.WriteString("d - show dijkstra data\n")
	if settings.solve {
		labelWriter.Color = green
	}
	labelWriter.WriteString("p - path (solve) the maze\n")
	labelWriter.Color = color.White
	if settings.longestPath {
		labelWriter.Color = green
	}
	labelWriter.WriteString("l - show longest path from start\n")
	labelWriter.Color = color.White
	if settings.floodFill {
		labelWriter.Color = green
	}
	labelWriter.WriteString("f - flood fill\n")
	labelWriter.Color = color.White

	labelWriter.Draw(target, pixel.IM)
}

const thickness = 5

func main() {
	set := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var (
		stats bool
	)
	set.BoolVar(&stats, "stats", false, "Print maze algorithm statistics")

	if len(os.Args) > 1 {
		_ = set.Parse(os.Args[1:])

		if stats {
			algos := []func(g grid.Grid){
				algorithms.BinarySearch,
				algorithms.Sidewinder,
				algorithms.AldousBroder,
				algorithms.Wilsons,
				algorithms.AldousBroderWilsons,
				algorithms.HuntAndKill,
				algorithms.RecursiveBacktracker,
			}
			names := make([]string, len(algos))
			var longestName int
			for i, algo := range algos {
				name := filepath.Base(runtime.FuncForPC(reflect.ValueOf(algo).Pointer()).Name())
				names[i] = name[len("algorithms."):]
				if len(names[i]) > longestName {
					longestName = len(names[i])
				}
			}

			iterations := 30
			rows := 20
			cols := 20
			size := rows * cols
			fmt.Printf("average over %d runs\n", iterations)
			fmt.Println(strings.Repeat(" ", longestName)+" ends | 4way | corridors")
			for i, algo := range algos {
				var stats grid.MazeStats
				for count := 0; count < iterations; count++ {
					g := grid.New(rows, cols)
					algo(g)
					s := g.Statistics()
					stats.DeadEnds += s.DeadEnds
					stats.Corridors += s.Corridors
					stats.FourWay += s.FourWay
				}

				deadEndPercent := math.Round(float64(stats.DeadEnds) / float64(size*iterations) * 100)
				corridorPercent := math.Round(float64(stats.Corridors) / float64(size*iterations) * 100)
				fourWayPercent := float64(stats.FourWay) / float64(size*iterations) * 100
				fmt.Printf("%"+strconv.FormatInt(int64(longestName), 10)+"s: %2d%% | %.1f%% | %d%%\n", names[i], int(deadEndPercent), fourWayPercent, int(corridorPercent))
			}
			os.Exit(0)
		}
	}

	pixelgl.Run(run)
}
