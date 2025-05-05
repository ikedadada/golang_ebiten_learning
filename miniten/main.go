package main

import (
	"math/rand/v2"

	"github.com/eihigh/miniten"
)

var (
	y  = 0.0
	vy = 0.0

	frames = 0
	walls  = []Wall{}
)

func main() {
	miniten.Run(draw)
}

func draw() {
	miniten.DrawImage("sky.png", 0, 0)

	if miniten.IsClicked() {
		vy = -6
	}
	vy += 0.5
	y += vy

	if y > 360 {
		y = 360
	}

	miniten.DrawImage("gopher.png", 100, int(y))

	frames += 1
	if frames%100 == 0 {
		walls = append(walls, NewWall())
	}

	for _, wall := range walls {
		wall.Update()
		wall.Draw()
	}
}

type pipe struct {
	x int
	y int
}

type Wall struct {
	topPipe    *pipe
	bottomPipe *pipe
}

func NewWall() Wall {

	holeY := rand.N(200)
	holeHeight := 100

	topPipe := pipe{
		x: 640,
		y: int(holeY - 360),
	}
	bottomPipe := pipe{
		x: 640,
		y: int(holeY + holeHeight),
	}

	return Wall{
		topPipe:    &topPipe,
		bottomPipe: &bottomPipe,
	}
}

func (w *Wall) Draw() {
	miniten.DrawImage("wall.png", w.topPipe.x, w.topPipe.y)
	miniten.DrawImage("wall.png", w.bottomPipe.x, w.bottomPipe.y)
}

func (w *Wall) Update() {
	w.topPipe.x -= 3
	w.bottomPipe.x -= 3
}
