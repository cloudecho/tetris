package tetris

var InvalidPoint = Point{0, -SHAPE_SIZE}

type Point struct {
	left int
	top  int
}

func (p *Point) valid() bool {
	return p.top > InvalidPoint.top
}

func (a Point) equals(b Point) bool {
	return a.left == b.left && a.top == b.top
}

// {from, to Point}
type Moving struct {
	from Point
	to   Point
}

type Area struct {
	x  int // left
	y  int // top
	x2 int // left+width
	y2 int // top+height
}

func (a *Area) outOfBounds() bool {
	return outOfBounds(a.x, a.y) || outOfBounds(a.x2, a.y2)
}

func outOfBounds(left, top int) bool {
	return left < 0 || left > COL-1 || top < 0 || top > ROW-1
}
