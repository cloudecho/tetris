package tetris

import (
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
	// score table
	scores = [SHAPE_SIZE]int{100, 300, 500, 700}

	// speed table, mapping level to speed
	speeds = [LEVELS]int{1500, 1300, 1000, 800, 500, 300}
)

type Game struct {
	state     int32
	model     [ROW][COL]uint8
	currShape *Shape
	oldShape  *Shape // for rotate
	nextShape *Shape

	pos        Point // position of current shape
	level      uint8 // starts from 0
	score      uint64
	rows       uint
	waterLevel int

	chanMoving chan *Moving // shape moving chan
	chanRedraw chan *Area   // redraw area chan
	chanHiligh chan int     // highlight row chan
	chanLevel  chan uint8   // level chan
	chanScore  chan uint64  // score chan
	chanState  chan int32   // state chan
	chanNexts  chan bool    // show next shape

	m       sync.Mutex
	stateOk *sync.Cond
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
		chanMoving: make(chan *Moving),
		chanRedraw: make(chan *Area),
		chanHiligh: make(chan int),
		chanLevel:  make(chan uint8),
		chanScore:  make(chan uint64),
		chanState:  make(chan int32),
		chanNexts:  make(chan bool),
	}
	g.stateOk = sync.NewCond(&g.m)
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
	g.level = 0
	g.waterLevel = ROW
	g.score = 0
	g.rows = 0
}

// init g.pos and notiy ui
func (g *Game) landing() {
	g.pos = Point{
		left: (COL-SHAPE_SIZE)/2 + 1,
		top:  -g.currShape.bounds().y,
	}

	g.chanMoving <- &Moving{InvalidPoint, g.pos}
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

	g.landing()
	g.chanNexts <- true
	g.chanScore <- g.score
	g.chanLevel <- g.level
	g.changeState(STATE_GAMING)
	go g.startGame()
}

func (g *Game) startGame() {
	log.Println("start to game")
	time.Sleep(time.Second)

	for g.continueGame() {
		time.Sleep(g.speed())

		// check if paused
		g.m.Lock()
		for STATE_PAUSED == g.state {
			log.Println("wait to continue")
			g.stateOk.Wait()
		}
		g.m.Unlock()
	}

	log.Println("game over")
}

func (g *Game) continueGame() bool {
	g.m.Lock()
	defer g.m.Unlock()

	mv, err := g.currShape.moveDown(g.pos)
	if g.canMove(err, mv) {
		g.moveTo(mv)
		return true
	}

	g.updateWaterLevel()

	// if game over
	if 0 == g.waterLevel {
		g.changeState(SATE_GAMEOVER)
		return false
	}

	g.updateModel()
	g.promote()

	g.currShape = g.nextShape
	g.nextShape = randShape()
	g.landing()
	g.chanNexts <- true

	return true
}

func (g *Game) changeState(state int32) {
	g.state = state
	g.chanState <- state
}

func (g *Game) moveTo(mv *Moving) {
	g.pos = mv.to
	g.chanMoving <- mv
}

func (g *Game) updateWaterLevel() {
	k := g.pos.top + g.currShape.bounds().y
	if g.waterLevel > k {
		g.waterLevel = k
	}
}

func (g *Game) updateModel() {
	p := g.pos
	b := g.currShape.bounds()
	d := &g.currShape.data

	for i := b.x; i <= b.x2; i++ {
		for j := b.y; j <= b.y2; j++ {
			if d[j][i] > 0 {
				g.model[p.top+j][p.left+i] = d[j][i]
			}
		}
	}
}

func (g *Game) promote() {
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
			g.hilighRow(k)
			g.eraseRow(k)
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
	return scores[n-1] + 100*int(level)
}

func (g *Game) hilighRow(k int) {
	g.chanHiligh <- k
}

// Erase k-th row of g.model
func (g *Game) eraseRow(k int) {
	m := &g.model
	top := g.waterLevel
	for i := k; i >= top && i > 1; i-- {
		m[i] = m[i-1]
	}
	g.waterLevel++

	// notify gui to redraw the area(top~k rows)
	area := &Area{y: top, y2: k}
	g.chanRedraw <- area
}

// Returns an int value in miliseconds
func (g *Game) speed() time.Duration {
	return time.Duration(speeds[g.level]) * time.Millisecond
}

func (g *Game) pause() {
	g.m.Lock()
	defer g.m.Unlock()

	if g.state != STATE_GAMING {
		log.Printf("could not pause as current state is %d", g.state)
		return
	}
	g.changeState(STATE_PAUSED)
	log.Println("game paused")
}

func (g *Game) resume() {
	g.m.Lock()
	defer g.m.Unlock()

	if g.state != STATE_PAUSED {
		log.Printf("could not resume as current state is %d", g.state)
		return
	}
	g.changeState(STATE_GAMING)
	g.stateOk.Signal()
	log.Println("game resumed")
}

func (g *Game) rotate() {
	g.m.Lock()
	defer g.m.Unlock()

	newShape, mv, err := g.currShape.rotate(g.pos)

	if g.canMoveShape(newShape, err, mv) {
		g.oldShape = g.currShape
		g.currShape = newShape
		g.moveTo(mv)
	}
}

func (g *Game) moveLeft() {
	g.m.Lock()
	defer g.m.Unlock()

	mv, err := g.currShape.moveLeft(g.pos)
	if g.canMove(err, mv) {
		g.moveTo(mv)
	}
}

func (g *Game) moveRight() {
	g.m.Lock()
	defer g.m.Unlock()

	mv, err := g.currShape.moveRight(g.pos)
	if g.canMove(err, mv) {
		g.moveTo(mv)
	}
}

func (g *Game) dropDown() {
	g.m.Lock()
	defer g.m.Unlock()

	from := g.pos
	to := InvalidPoint
	s := g.currShape

	for mv, err := s.moveDown(from); g.canMove(err, mv); {
		to = mv.to
		mv, err = s.moveDown(to)
	}

	if to.valid() {
		g.moveTo(&Moving{from, to})
	}
}

func (g *Game) canMove(err error, mv *Moving) bool {
	return g.canMoveShape(g.currShape, err, mv)
}

// Return true if can move
func (g *Game) canMoveShape(shape *Shape, err error, mv *Moving) bool {
	if err != nil {
		return false
	}

	pos := mv.to
	a := shape.area(pos)
	d := &shape.data
	m := &g.model

	for i := a.x; i <= a.x2; i++ {
		for j := a.y; j <= a.y2; j++ {
			if m[j][i]&d[j-pos.top][i-pos.left] > 0 { // conflict
				return false
			}
		}
	}

	return true
}
