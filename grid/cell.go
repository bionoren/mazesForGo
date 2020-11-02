package grid

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

func (d Direction) Reverse() Direction {
	switch d {
	case NORTH:
		return SOUTH
	case EAST:
		return WEST
	case SOUTH:
		return NORTH
	case WEST:
		return EAST
	}

	return NORTH
}

type Cell struct {
	index     int
	row, col  int
	neighbors [WEST + 1]*Cell
	openings  [WEST + 1]bool
}

func (c Cell) Index() int {
	return c.index
}

func (c Cell) Row() int {
	return c.row
}

func (c Cell) Col() int {
	return c.col
}

func (c *Cell) addNeighbor(neighbor *Cell, dir Direction) {
	c.neighbors[dir] = neighbor
}

func (c Cell) Neighbor(dir Direction) *Cell {
	return c.neighbors[dir]
}

func (c Cell) HasNeighbor(dir Direction) bool {
	return c.neighbors[dir] != nil
}

func (c *Cell) connect(dir Direction) {
	c.openings[dir] = true
	c.neighbors[dir].openings[dir.Reverse()] = true
}

func (c *Cell) disconnect(dir Direction) {
	c.openings[dir] = false
	c.neighbors[dir].openings[dir.Reverse()] = false
}

func (c Cell) Connected(dir Direction) bool {
	return c.openings[dir]
}
