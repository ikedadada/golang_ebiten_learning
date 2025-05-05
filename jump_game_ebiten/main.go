package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func main() {
	ebiten.SetWindowTitle("jump_game_ebiten")
	ebiten.SetWindowSize(640, 360)
	jumpGame := JumpGame{}
	jumpGame.Init()
	src, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		panic(err)
	}
	fontFace = &text.GoTextFace{Source: src, Size: 16}
	ebiten.RunGame(&Game{
		scene: SceneTitle,
		game:  jumpGame,
	})
}

var (
	fontFace *text.GoTextFace
)

func drawTextMain(screen *ebiten.Image, message string, offsetX, offsetY float64) {
	// Draw the text on the screen using the loaded font
	op := &text.DrawOptions{}
	op.GeoM.Translate(320+offsetX, 180+offsetY)
	op.ColorScale.Scale(1, 1, 1, 1)
	op.LayoutOptions = text.LayoutOptions{
		PrimaryAlign: text.AlignCenter,
	}
	text.Draw(screen, message, fontFace, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(320+offsetX-2, 180+offsetY-2)
	op.ColorScale.Scale(0, 0, 0, 1)
	op.LayoutOptions = text.LayoutOptions{
		PrimaryAlign: text.AlignCenter,
	}
	text.Draw(screen, message, fontFace, op)
}

func drawScore(screen *ebiten.Image, score int) {
	// Draw the score on the screen using the loaded font
	scoreText := fmt.Sprintf("Score: %d", score)
	op := &text.DrawOptions{}
	op.GeoM.Translate(0, 0)
	op.ColorScale.Scale(0, 0, 0, 1)
	text.Draw(screen, scoreText, fontFace, op)
}

//go:embed assets/*.png
var fsys embed.FS

func drawImage(screen *ebiten.Image, name string, x, y float64) error {
	f, err := fsys.Open(name)
	defer f.Close()
	if err != nil {
		return err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	ebitenImg := ebiten.NewImageFromImage(img)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(ebitenImg, op)
	return nil
}

type Scene string

const (
	SceneTitle    Scene = "title"
	SceneGame     Scene = "game"
	SceneGameOver Scene = "gameover"
)

type Game struct {
	scene         Scene
	game          JumpGame
	isPrevClicked bool
	isJustClicked bool
}

func (g *Game) Update() error {
	g.isJustClicked = !g.isPrevClicked && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	g.isPrevClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	switch g.scene {
	case SceneTitle:
		g.game.Init()
		if g.isJustClicked {
			g.scene = SceneGame
		}
	case SceneGame:
		isFinished := g.game.Update()
		if isFinished {
			g.scene = SceneGameOver
		}
	case SceneGameOver:
		if g.isJustClicked {
			g.scene = SceneTitle
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	drawImage(screen, "assets/sky.png", 0, 0)
	g.game.Draw(screen)
	switch g.scene {
	case SceneTitle:
		drawTextMain(screen, "はねるGopherくんゲーム", 0, 0)
		drawTextMain(screen, "クリックでスタート", 0, 20)
		break
	case SceneGame:
		break
	case SceneGameOver:
		drawTextMain(screen, "Game Over", 0, 0)
		drawTextMain(screen, "クリックでタイトルへ", 0, 20)
		break
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}

type JumpGame struct {
	frames     int
	gopher     Gopher
	walls      []Wall
	isFinished bool
}

func (g *JumpGame) Init() {
	g.frames = 0
	g.gopher = NewGopher()
	g.walls = []Wall{}
	g.isFinished = false
}

func (g *JumpGame) Update() (isFinished bool) {
	if g.isFinished {
		return true
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
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
			g.isFinished = true
		}
	}
	return g.isFinished
}

func (g *JumpGame) Draw(screen *ebiten.Image) {
	g.gopher.Draw(screen)
	for _, wall := range g.walls {
		wall.Draw(screen)
	}
	drawScore(screen, g.frames)
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

func (g Gopher) Draw(screen *ebiten.Image) {
	drawImage(screen, "assets/gopher.png", g.x, g.y)
}

func (g Gopher) GetCenter() Point {
	return Point{
		x: g.x + 30/2,
		y: g.y + 38/2,
	}
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

func (w *Wall) Draw(screen *ebiten.Image) {
	drawImage(screen, "assets/wall.png", w.topPipe.x, w.topPipe.y)
	drawImage(screen, "assets/wall.png", w.bottomPipe.x, w.bottomPipe.y)
}

func (w *Wall) Update() {
	w.topPipe.x -= 3
	w.bottomPipe.x -= 3
}
