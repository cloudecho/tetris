package tetris

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	ROW = 19
	COL = 11
)

const (
	STATE_ZERO = iota
	SATE_GAMEOVER
	STATE_GAMING
	STATE_PAUSED
)

var POSITION_OUT_OF_BOUNDS = errors.New("Position out of bounds")

type Point struct {
	left  int
	top   int
	oLeft int // old left
	oTop  int // old top
}

var INVALID_POINT = Point{-1, -1, -1, -1}

func (p Point) moveLeft() (Point, error) {
	if p.left-1 < 0 {
		return INVALID_POINT, POSITION_OUT_OF_BOUNDS
	}
	p.oLeft = p.left
	p.oTop = p.top
	p.left--
	return p, nil
}

func (p Point) moveRight() (Point, error) {
	if p.left+1 > COL-1 {
		return INVALID_POINT, POSITION_OUT_OF_BOUNDS
	}
	p.oLeft = p.left
	p.oTop = p.top
	p.left++
	return p, nil
}

func (p Point) moveDown() (Point, error) {
	if p.top+1 > ROW-1 {
		return INVALID_POINT, POSITION_OUT_OF_BOUNDS
	}
	p.oLeft = p.left
	p.oTop = p.top
	p.top++
	return p, nil
}

func (p Point) sendTo(out chan<- Point) {
	out <- p
}

type Game struct {
	state     int32
	model     [ROW][COL]uint8
	currShape *Shape
	nextShape *Shape

	pos   Point // position of current shape
	level uint8
	score uint

	p chan Point // position chan
	l chan uint8 // level chan
	s chan uint  // score chan
	t chan int32 // state chan
	n chan bool  // show next shape

	m sync.Mutex
}

func NewGame() *Game {
	return &Game{
		state:     STATE_ZERO,
		currShape: randShape(),
		nextShape: randShape(),
		pos:       landingPosition(),
		p:         make(chan Point),
		l:         make(chan uint8),
		s:         make(chan uint),
		n:         make(chan bool),
	}
}

func (g *Game) start() error {
	g.m.Lock()
	defer g.m.Unlock()

	if g.state > SATE_GAMEOVER {
		return fmt.Errorf("could not start as current state is: %d", g.state)
	}

	g.state = STATE_GAMING
	go moveForward(g)
	return nil
}

func moveForward(g *Game) {
	for {
		g.pos.sendTo(g.p)
		if g.pos.oLeft < 0 {
			g.n <- true
		}
		time.Sleep(speedDuration(g))

		// TODO do more checks
		pos, err := g.pos.moveDown()
		if err != nil || shapeOutOfBound(pos, g.currShape) {
			g.pos = landingPosition()
			g.m.Lock()
			{
				g.currShape = g.nextShape
				g.nextShape = randShape()
				g.n <- true
			}
			g.m.Unlock()
			continue
		}
		g.pos = pos
	}
}

func shapeOutOfBound(pos Point, shape *Shape) bool {
	if pos.top < 0 || pos.left < 0 {
		return false
	}

	// Show the current shape
	for i := 0; i < SHAPE_SIZE; i++ { // left
		for j := 0; j < SHAPE_SIZE; j++ { // top
			if shape.data[j][i] <= 0 {
				continue
			}
			if checkOutOfBound(pos.left+i, pos.top+j) {
				return true
			}
		}
	}

	return false
}

func checkOutOfBound(left, top int) bool {
	return top > ROW-1 || left > COL-1 || top < 0 || left < 0
}

// Returns an int value in miliseconds
func speed(g *Game) int {
	// TODO
	return 300
}

func speedDuration(g *Game) time.Duration {
	d := int64(speed(g)) * int64(time.Millisecond)
	return time.Duration(d)
}

func landingPosition() Point {
	return Point{
		left:  (COL-SHAPE_SIZE)/2 + 1,
		top:   0,
		oLeft: -1,
		oTop:  -1,
	}
}

func (g *Game) pause() {

}

func (g *Game) resume() {

}

func (g *Game) rotate() {
	g.m.Lock()
	defer g.m.Unlock()
	g.currShape = g.currShape.rotate()
}

func (g *Game) moveLeft() {

}

func (g *Game) moveRight() {

}

func (g *Game) dropDown() {

}

func (g *Game) reset() {

}
