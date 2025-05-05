package main

import (
	"math/rand/v2"

	"github.com/eihigh/miniten"
)

type Scene string

const (
	SceneTitle    Scene = "title"
	SceneGame     Scene = "game"
	SceneGameOver Scene = "gameover"
)

var (
	scene         = SceneTitle
	game          = Game{}
	isPrevClicked bool
	isJustClicked bool
)

func main() {
	miniten.Run(draw)
}

func draw() {
	miniten.DrawImage("sky.png", 0, 0)

	isJustClicked = !isPrevClicked && miniten.IsClicked()
	isPrevClicked = miniten.IsClicked()

	switch scene {
	case SceneTitle:
		miniten.DrawImage("gopher.png", 100, 0)
		miniten.Println("はねるGopherくんゲーム")
		miniten.Println("クリックでスタート")
		if isJustClicked {
			game.Init()
			scene = SceneGame
		}
	case SceneGame:
		game.Update()
		game.Draw()
	case SceneGameOver:
		game.Draw()
		miniten.Println("Game Over")
		miniten.Println("クリックでタイトルへ")
		if isJustClicked {
			scene = SceneTitle
		}
	}
}

type Game struct {
	frames int
	gopher Gopher
	walls  []Wall
}

func (g *Game) Init() {
	g.frames = 0
	g.gopher = NewGopher()
	g.walls = []Wall{}
}

func (g *Game) Update() {

	if miniten.IsClicked() {
		g.gopher.Jump()
	}
	g.gopher.Update()
	gopherCenter := g.gopher.GetCenter()

	g.frames += 1
	if g.frames%100 == 0 {
		g.walls = append(g.walls, NewWall())
	}

	for _, wall := range g.walls {
		wall.Update()
		if HitCheck(gopherCenter, wall) {
			scene = SceneGameOver
		}
	}
}

func (g *Game) Draw() {
	g.gopher.Draw()
	for _, wall := range g.walls {
		wall.Draw()
	}
	miniten.Println("Score:", g.frames)
}

type Point struct {
	x float64
	y float64
}

type Gopher struct {
	Point
	vy float64
}

func NewGopher() Gopher {
	return Gopher{
		Point: Point{
			x: 100,
			y: 0.0,
		},
		vy: 0.0,
	}
}

func (g *Gopher) Update() {
	g.vy += 0.5
	g.y += g.vy
	if g.y > 360 {
		g.y = 360
	}
	if g.y < 0 {
		g.y = 0
	}
}

func (g *Gopher) Jump() {
	g.vy = -6
}

func (g Gopher) Draw() {
	miniten.DrawImage("gopher.png", int(g.x), int(g.y))
}

func (g Gopher) GetCenter() Point {
	return Point{
		x: g.x + 30/2,
		y: g.y + 38/2,
	}
}

func (g Gopher) Debug() {
	miniten.Println("Gopher", g.x, int(g.y), g.vy)
}

func pipeHitCheck(gopherCenter Point, pipe *Point) bool {
	left := pipe.x
	right := left + 100
	top := pipe.y
	bottom := top + 360
	if left < gopherCenter.x && gopherCenter.x < right {
		if top < gopherCenter.y && gopherCenter.y < bottom {
			return true
		}
	}
	return false
}

func HitCheck(gopherCenter Point, wall Wall) bool {
	if pipeHitCheck(gopherCenter, wall.topPipe) {
		return true
	}
	if pipeHitCheck(gopherCenter, wall.bottomPipe) {
		return true
	}
	return false
}

type Wall struct {
	topPipe    *Point
	bottomPipe *Point
}

func NewWall() Wall {

	holeY := rand.Float64() * 200
	holeHeight := float64(100)

	topPipe := Point{
		x: 640,
		y: holeY - 360,
	}
	bottomPipe := Point{
		x: 640,
		y: holeY + holeHeight,
	}

	return Wall{
		topPipe:    &topPipe,
		bottomPipe: &bottomPipe,
	}
}

func (w *Wall) Draw() {
	miniten.DrawImage("wall.png", int(w.topPipe.x), int(w.topPipe.y))
	miniten.DrawImage("wall.png", int(w.bottomPipe.x), int(w.bottomPipe.y))
}

func (w *Wall) Update() {
	w.topPipe.x -= 3
	w.bottomPipe.x -= 3
}
