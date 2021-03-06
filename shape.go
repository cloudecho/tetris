package tetris

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

const SHAPE_SIZE = 4

var shapes = []*Shape{
	// one point
	{id: 0, next: 0, data: shapeData{
		{0, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	// two points
	{id: 1, next: 2, data: shapeData{
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 2, next: 1, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	// 3-4
	{id: 3, next: 4, data: shapeData{
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 4, next: 3, data: shapeData{
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	// three points
	{id: 5, next: 6, data: shapeData{
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 6, next: 5, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	// 7-10
	{id: 7, next: 8, data: shapeData{
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 8, next: 9, data: shapeData{
		{0, 1, 0, 0},
		{1, 0, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 9, next: 10, data: shapeData{
		{0, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 10, next: 7, data: shapeData{
		{0, 0, 0, 0},
		{1, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	// 11-14
	{id: 11, next: 12, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 12, next: 13, data: shapeData{
		{0, 0, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 13, next: 14, data: shapeData{
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 14, next: 11, data: shapeData{
		{0, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	// four points
	{id: 15, next: 16, data: shapeData{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 16, next: 15, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0}}},

	// 17-20
	{id: 17, next: 18, data: shapeData{
		{0, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 18, next: 19, data: shapeData{
		{1, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 19, next: 20, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 20, next: 17, data: shapeData{
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0}}},

	// 21-24
	{id: 21, next: 22, data: shapeData{
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0}}},

	{id: 22, next: 23, data: shapeData{
		{0, 0, 0, 0},
		{0, 1, 1, 1},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 23, next: 24, data: shapeData{
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 0}}},

	{id: 24, next: 21, data: shapeData{
		{0, 0, 0, 1},
		{0, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	// 25-28
	{id: 25, next: 26, data: shapeData{
		{1, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 26, next: 27, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 27, next: 28, data: shapeData{
		{0, 1, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 28, next: 25, data: shapeData{
		{0, 1, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	// 29 -30
	{id: 29, next: 30, data: shapeData{
		{0, 1, 1, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 30, next: 29, data: shapeData{
		{1, 0, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	// 31 -32
	{id: 31, next: 32, data: shapeData{
		{1, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 32, next: 31, data: shapeData{
		{0, 1, 0, 0},
		{1, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 0, 0}}},

	// 33-34
	{id: 33, next: 34, data: shapeData{
		{0, 0, 0, 0},
		{1, 1, 0, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 0}}},

	{id: 34, next: 33, data: shapeData{
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0}}},

	// 35-36
	{id: 35, next: 36, data: shapeData{
		{0, 0, 0, 0},
		{0, 0, 1, 1},
		{1, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 36, next: 35, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0}}},

	// 37-40
	{id: 37, next: 38, data: shapeData{
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 1},
		{0, 0, 0, 0}}},

	{id: 38, next: 39, data: shapeData{
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0}}},

	{id: 39, next: 40, data: shapeData{
		{0, 0, 0, 0},
		{1, 0, 0, 0},
		{0, 1, 1, 1},
		{0, 0, 0, 0}}},

	{id: 40, next: 37, data: shapeData{
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 1, 0, 0}}},

	// 41-44
	{id: 41, next: 42, data: shapeData{
		{0, 0, 0, 0},
		{0, 1, 1, 1},
		{1, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 42, next: 43, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0}}},

	{id: 43, next: 44, data: shapeData{
		{0, 0, 0, 0},
		{0, 0, 0, 1},
		{1, 1, 1, 0},
		{0, 0, 0, 0}}},

	{id: 44, next: 41, data: shapeData{
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0}}},

	// 45-48 (3 points)
	{id: 45, next: 46, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 46, next: 47, data: shapeData{
		{0, 0, 0, 0},
		{1, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0}}},

	{id: 47, next: 48, data: shapeData{
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 48, next: 45, data: shapeData{
		{1, 0, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	// 49-52 (3 points)
	{id: 49, next: 50, data: shapeData{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}},

	{id: 50, next: 51, data: shapeData{
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 51, next: 52, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0}}},

	{id: 52, next: 49, data: shapeData{
		{0, 0, 1, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	// 53
	{id: 53, next: 53, data: shapeData{
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},
}

type (
	Shape struct {
		id   int
		next int
		data shapeData
	}

	shapeData [SHAPE_SIZE][SHAPE_SIZE]uint8

	shapeBounds Area
)

var shapeBoundsMap = make(map[int]shapeBounds)

func randShape() *Shape {
	k := rand.Intn(len(shapes))
	return shapes[k]
}

func (s *Shape) bounds() shapeBounds {
	if len(shapeBoundsMap) == 0 {
		log.Fatalln("should init first")
	}
	return shapeBoundsMap[s.id]
}

func (s *Shape) area(o Point) *Area {
	b := s.bounds()
	return &Area{
		x:  o.left + b.x,
		y:  o.top + b.y,
		x2: o.left + b.x2,
		y2: o.top + b.y2,
	}
}

var ErrorMoving = errors.New("error moving")

func (s *Shape) rotate(o Point) (*Shape, *Moving, error) {
	newShape := shapes[s.next]
	mv, err := checkMoving(newShape.area(o), o, o)
	return newShape, mv, err
}

func (s *Shape) moveLeft(from Point) (*Moving, error) {
	to := from // copy
	to.left--
	return checkMoving(s.area(to), from, to)
}

func (s *Shape) moveRight(from Point) (*Moving, error) {
	to := from // copy
	to.left++
	return checkMoving(s.area(to), from, to)
}

func (s *Shape) moveDown(from Point) (*Moving, error) {
	to := from // copy
	to.top++
	return checkMoving(s.area(to), from, to)
}

func checkMoving(a *Area, from, to Point) (*Moving, error) {
	if a.outOfBounds() {
		return nil, ErrorMoving
	}
	return &Moving{from, to}, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())

	for id, shape := range shapes {
		shapeBoundsMap[id] = computeBounds(shape)
	}
}

// Shape bounds at the point of zero
func computeBounds(a *Shape) shapeBounds {
	d := &a.data
	x := SHAPE_SIZE
	y := SHAPE_SIZE
	x2 := 0
	y2 := 0

	for i := 0; i < SHAPE_SIZE; i++ {
		for j := 0; j < SHAPE_SIZE; j++ {
			if d[j][i] == 0 {
				continue
			}
			if x > i {
				x = i
			}
			if y > j {
				y = j
			}
			if x2 < i {
				x2 = i
			}
			if y2 < j {
				y2 = j
			}
		}
	}

	return shapeBounds{x, y, x2, y2}
}
