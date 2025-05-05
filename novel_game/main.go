package main

import (
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type game struct {
	bg     *ebiten.Image
	person *ebiten.Image
}

func newGame() (*game, error) {
	g := &game{}
	img, _, err := ebitenutil.NewImageFromFile("assets/bg.jpg")
	if err != nil {
		return nil, err
	}
	g.bg = img
	img, _, err = ebitenutil.NewImageFromFile("assets/person.png")
	if err != nil {
		return nil, err
	}
	g.person = img
	return g, nil
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screen.Bounds().Dx())/float64(g.bg.Bounds().Dx()), float64(screen.Bounds().Dy())/float64(g.bg.Bounds().Dy()))
	screen.DrawImage(g.bg, op)

	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(g.person, op)

	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("ノベルゲーム")
	g, err := newGame()
	if err != nil {
		panic(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
