package tetris

import (
	"errors"
	"log"
	"sync"
	"time"
)

const (
	ROW = 19
	COL = 11

	LEVELS = 6
)

const (
	STATE_ZERO = iota
	SATE_GAMEOVER
	STATE_GAMING
	STATE_PAUSED
)

var (
	ErrorPositionOutOfBounds = errors.New("position out of bounds")

	scoreTable = [SHAPE_SIZE]int{100, 300, 500, 700}
	speedTable = [LEVELS]int{1500, 1300, 1000, 800, 500, 300}
)

type Point struct {
	left  int
	top   int
	oleft int // old left
	otop  int // old top
}

var InvalidPoint = Point{-SHAPE_SIZE, -SHAPE_SIZE, -SHAPE_SIZE, -SHAPE_SIZE}

func (p *Point) isValid() bool {
	return p.left > InvalidPoint.left && p.top > InvalidPoint.top
}

func (p Point) moveLeft() (Point, error) {
	if p.left+SHAPE_SIZE-2 < 0 {
		return InvalidPoint, ErrorPositionOutOfBounds
	}
	p.oleft = p.left
	p.otop = p.top
	p.left--
	return p, nil
}

func (p Point) moveRight() (Point, error) {
	if p.left+1 > COL-1 {
		return InvalidPoint, ErrorPositionOutOfBounds
	}
	p.oleft = p.left
	p.otop = p.top
	p.left++
	return p, nil
}

func (p Point) moveDown() (Point, error) {
	if p.top+1 > ROW-1 {
		return InvalidPoint, ErrorPositionOutOfBounds
	}
	p.oleft = p.left
	p.otop = p.top
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

	pos        Point // position of current shape
	level      uint8 // starts from 0
	score      uint64
	rows       uint
	waterLevel int

	chanPos    chan Point  // position chan
	chanRedraw chan Point  // redraw area chan
	chanHiligh chan int    // highlight row chan
	chanLevel  chan uint8  // level chan
	chanScore  chan uint64 // score chan
	chanState  chan int32  // state chan
	chanNexts  chan bool   // show next shape

	m          sync.Mutex
	resumeCond *sync.Cond
}

func NewGame() *Game {
	g := &Game{
		state:      STATE_ZERO,
		currShape:  randShape(),
		nextShape:  randShape(),
		level:      0,
		waterLevel: ROW,
		score:      0,
		rows:       0,
		chanPos:    make(chan Point),
		chanRedraw: make(chan Point),
		chanHiligh: make(chan int),
		chanLevel:  make(chan uint8),
		chanScore:  make(chan uint64),
		chanState:  make(chan int32),
		chanNexts:  make(chan bool),
	}
	g.pos = landingPoint(g.currShape)
	g.resumeCond = sync.NewCond(&g.m)
	return g
}

func (g *Game) reset() {
	log.Println("reset game status")

	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			g.model[i][j] = 0
		}
	}

	g.state = STATE_ZERO
	g.currShape = randShape()
	g.nextShape = randShape()
	g.pos = landingPoint(g.currShape)
	g.level = 0
	g.waterLevel = ROW
	g.score = 0
	g.rows = 0
}

func (g *Game) start() {
	g.m.Lock()
	defer g.m.Unlock()

	if g.state > SATE_GAMEOVER {
		log.Printf("could not start as current state is %d", g.state)
		return
	}

	if g.state > STATE_ZERO {
		g.reset()
		g.chanState <- STATE_ZERO
	}

	g.pos.sendTo(g.chanPos)
	g.chanNexts <- true
	g.chanScore <- g.score
	g.chanLevel <- g.level
	changeState(STATE_GAMING, g)
	go startGame(g)
}

func startGame(g *Game) {
	log.Println("start to game")
	time.Sleep(time.Second)

	for continueGame(g) {
		time.Sleep(speed(g))

		// check if paused
		g.m.Lock()
		for STATE_PAUSED == g.state {
			log.Println("wait to continue")
			g.resumeCond.Wait()
		}
		g.m.Unlock()
	}

	log.Println("game over")
}

func continueGame(g *Game) bool {
	g.m.Lock()
	defer g.m.Unlock()

	pos, err := g.pos.moveDown()
	if !isConflict(err, pos, g) {
		changePosition(pos, g)
		return true
	}

	// if game over
	if conflictWithModel(g.pos, g.currShape, g) {
		changeState(SATE_GAMEOVER, g)
		return false
	}

	updateWaterLevel(g)
	updateModel(g)
	promote(g)

	g.currShape = g.nextShape
	g.nextShape = randShape()
	g.pos = landingPoint(g.currShape)

	g.pos.sendTo(g.chanPos)
	g.chanNexts <- true

	return true
}

func changeState(state int32, g *Game) {
	g.state = state
	g.chanState <- state
}

func changePosition(pos Point, g *Game) {
	g.pos = pos
	pos.sendTo(g.chanPos)
}

func updateWaterLevel(g *Game) {
	k := g.pos.top + g.currShape.bounds().y
	if g.waterLevel > k {
		g.waterLevel = k
	}
}

func updateModel(g *Game) {
	p := g.pos
	b := g.currShape.bounds()
	d := g.currShape.data

	for i := b.x; i <= b.x2; i++ {
		for j := b.y; j <= b.y2; j++ {
			if d[j][i] > 0 {
				g.model[p.top+j][p.left+i] = d[j][i]
			}
		}
	}
}

func promote(g *Game) {
	p := g.pos
	m := &g.model

	// erase promoted rows
	n := 0
	top := p.top
	for i := ROW - 1; i >= top; i-- { // top
		k := i
		for j := 0; j < COL; j++ { // left
			if m[i][j] == 0 {
				k = -1
				break
			}
		}
		// erase k-th row
		if k > 0 {
			hilighRow(k, g)
			eraseRow(k, g)
			n++
			top++
			i++
		}
	}

	if n == 0 {
		return
	}

	// compute rows & score
	newScore := uint64(earnScore(n, g.level))
	g.rows += uint(n)
	g.score += newScore
	g.chanScore <- g.score
	log.Printf("[promote] rows=%d(+%d) score=%d(+%d)", g.rows, n, g.score, newScore)

	// compute level
	l := uint8(g.rows / ROW)
	if l < LEVELS && l > g.level {
		log.Printf("[promote] level %d -> %d", g.level, l)
		g.level = l
		g.chanLevel <- l
	}
}

func earnScore(n int, level uint8) int {
	return scoreTable[n-1] + 100*int(level)
}

func hilighRow(k int, g *Game) {
	g.chanHiligh <- k
}

// Erase k-th row of g.model
func eraseRow(k int, g *Game) {
	m := &g.model
	top := g.waterLevel
	for i := k; i >= top && i > 1; i-- {
		m[i] = m[i-1]
	}
	g.waterLevel++

	// notify gui to redraw the area(top~k rows)
	area := Point{top: top, otop: k}
	g.chanRedraw <- area
}

// Returns an int value in miliseconds
func speed(g *Game) time.Duration {
	return time.Duration(speedTable[g.level]) * time.Millisecond
}

func landingPoint(s *Shape) Point {
	return Point{
		left:  (COL-SHAPE_SIZE)/2 + 1,
		top:   -s.bounds().y,
		oleft: -1,
		otop:  -1,
	}
}

func (g *Game) pause() {
	g.m.Lock()
	defer g.m.Unlock()

	if g.state != STATE_GAMING {
		log.Printf("could not pause as current state is %d", g.state)
		return
	}
	changeState(STATE_PAUSED, g)
	log.Println("game paused")
}

func (g *Game) resume() {
	g.m.Lock()
	defer g.m.Unlock()

	if g.state != STATE_PAUSED {
		log.Printf("could not resume as current state is %d", g.state)
		return
	}
	changeState(STATE_GAMING, g)
	g.resumeCond.Signal()
	log.Println("game resumed")
}

func (g *Game) rotate() {
	g.m.Lock()
	defer g.m.Unlock()

	newShape := g.currShape.rotate()
	if shapeOutOfBounds(g.pos, newShape) {
		return
	}
	if checkRotateConflict(newShape, g.pos, g) {
		return
	}

	p := g.pos
	pos := Point{p.left, p.top, p.left, p.top}
	g.currShape = newShape
	changePosition(pos, g)
}

func (g *Game) moveLeft() {
	g.m.Lock()
	defer g.m.Unlock()

	pos, err := g.pos.moveLeft()
	if isConflict(err, pos, g) {
		return
	}
	changePosition(pos, g)
}

func (g *Game) moveRight() {
	g.m.Lock()
	defer g.m.Unlock()

	pos, err := g.pos.moveRight()
	if isConflict(err, pos, g) {
		return
	}
	changePosition(pos, g)
}

func (g *Game) dropDown() {
	g.m.Lock()
	defer g.m.Unlock()

	oldPos := g.pos
	newPos := InvalidPoint

	for pos, err := oldPos.moveDown(); ; {
		if checkConflict(err, pos, g) {
			break
		}
		newPos = pos
		pos, err = newPos.moveDown()
	}

	if newPos.isValid() {
		newPos.oleft = oldPos.left
		newPos.otop = oldPos.top

		changePosition(newPos, g)
	}
}

// Return true if conflict
func isConflict(err error, pos Point, g *Game) bool {
	return err != nil ||
		shapeOutOfBounds(pos, g.currShape) ||
		conflictWithModel(pos, g.currShape, g)
}

func checkRotateConflict(shape *Shape, pos Point, g *Game) bool {
	return shapeOutOfBounds(pos, shape) ||
		conflictWithModel(pos, shape, g)
}

func conflictWithModel(pos Point, shape *Shape, g *Game) bool {
	b := shape.bounds()
	d := shape.data
	m := &g.model
	for i := b.x; i <= b.x2; i++ {
		for j := b.y; j <= b.y2; j++ {
			if m[pos.top+j][pos.left+i]&d[j][i] > 0 {
				return true
			}
		}
	}
	return false
}

func shapeOutOfBounds(pos Point, shape *Shape) bool {
	b := shape.bounds()
	return isOutOfBounds(pos.left+b.x, pos.top+b.y) ||
		isOutOfBounds(pos.left+b.x2, pos.top+b.y2)
}

func isOutOfBounds(left, top int) bool {
	return top > ROW-1 || left > COL-1 || top < 0 || left < 0
}
