package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type game struct {
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth / 2, outsideHeight / 2
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebitengine 入門")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
