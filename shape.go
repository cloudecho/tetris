package tetris

import (
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
		{1, 1, 1, 0},
		{0, 0, 0, 0},
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
		{1, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
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
		{1, 1, 1, 1},
		{0, 0, 0, 0},
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
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
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
		{1, 1, 0, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 34, next: 33, data: shapeData{
		{0, 0, 1, 1},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},

	// 35-38
	{id: 35, next: 36, data: shapeData{
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 1},
		{0, 0, 0, 0}}},

	{id: 36, next: 37, data: shapeData{
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0}}},

	{id: 37, next: 38, data: shapeData{
		{0, 0, 0, 0},
		{1, 0, 0, 0},
		{0, 1, 1, 1},
		{0, 0, 0, 0}}},

	{id: 38, next: 35, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{1, 0, 0, 0}}},

	// 39-42
	{id: 39, next: 40, data: shapeData{
		{0, 0, 0, 0},
		{0, 1, 1, 1},
		{1, 0, 0, 0},
		{0, 0, 0, 0}}},

	{id: 40, next: 41, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0}}},

	{id: 41, next: 42, data: shapeData{
		{0, 0, 0, 0},
		{0, 0, 0, 1},
		{1, 1, 1, 0},
		{0, 0, 0, 0}}},

	{id: 42, next: 39, data: shapeData{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0}}},

	// 43
	{id: 43, next: 43, data: shapeData{
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}},
}

type shapeData [SHAPE_SIZE][SHAPE_SIZE]uint8

type Shape struct {
	id   int
	next int
	data shapeData
}

func randShape() *Shape {
	k := rand.Intn(len(shapes))
	return shapes[k]
}

func (a *Shape) rotate() *Shape {
	return shapes[a.next]
}

type ShapeBounds struct {
	x  int // left
	y  int // top
	x2 int // left+width
	y2 int // top+height
}

var shapeBoundsMap = make(map[int]ShapeBounds)

func (a *Shape) bounds() ShapeBounds {
	if len(shapeBoundsMap) == 0 {
		log.Fatalln("should init first")
	}
	return shapeBoundsMap[a.id]
}

func computeBounds(a *Shape) ShapeBounds {
	d := a.data
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

	return ShapeBounds{x, y, x2, y2}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	for id, shape := range shapes {
		shapeBoundsMap[id] = computeBounds(shape)
	}
}
