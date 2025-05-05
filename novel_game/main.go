package main

import (
	_ "embed"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type game struct {
	bg        *ebiten.Image
	person    *ebiten.Image
	cat       *ebiten.Image
	messageBG *ebiten.Image
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
	img, _, err = ebitenutil.NewImageFromFile("assets/cat.png")
	if err != nil {
		return nil, err
	}
	g.cat = img

	img, _, err = ebitenutil.NewImageFromFile("assets/message-bg.png")
	if err != nil {
		return nil, err
	}
	g.messageBG = img
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
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(screen.Bounds().Dx()), 0)
	screen.DrawImage(g.person, op)

	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(g.cat, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screen.Bounds().Dx())/float64(g.messageBG.Bounds().Dx()), float64(screen.Bounds().Dy()/3)/float64(g.messageBG.Bounds().Dy()))
	op.GeoM.Translate(0, float64(screen.Bounds().Dy()*2/3))
	op.ColorScale.ScaleAlpha(0.5)
	screen.DrawImage(g.messageBG, op)

	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

//go:embed shader.kage
var shader1Src []byte
var (
	postEffectShader *ebiten.Shader
)

func (g *game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	if postEffectShader == nil {
		s, err := ebiten.NewShader(shader1Src)
		if err != nil {
			panic(err)
		}
		postEffectShader = s
	}

	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = offscreen
	op.GeoM = geoM
	b := offscreen.Bounds()
	screen.DrawRectShader(b.Dx(), b.Dy(), postEffectShader, op)
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
