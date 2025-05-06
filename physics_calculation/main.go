package main

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
)

const (
	boxWidth  = 40
	boxHeight = 40
)

type Box struct {
	x, y float64
	w, h float64
	r    float64
}

type game struct {
	w, h        float64
	space       *cp.Space
	boxImage    *ebiten.Image
	emptyImage  *ebiten.Image
	lastClicked int
}

//go:embed assets/*
var fsys embed.FS

func newGame() (*game, error) {
	g := &game{
		w: 800,
		h: 600,
	}

	g.space = cp.NewSpace()
	g.space.Iterations = 1000
	g.space.SetGravity(cp.Vector{X: 0, Y: 500})

	deltaY := 20.
	bf := cp.NewStaticBody()
	bf.SetPosition(cp.Vector{X: g.w / 2, Y: g.h - deltaY})
	bf.UserData = "line"
	sf := cp.NewBox(bf, g.w, 5, 0.0)
	sf.SetFriction(1.0)
	sf.SetElasticity(0.0)
	g.space.AddBody(bf)
	g.space.AddShape(sf)

	f, err := fsys.Open("assets/box.png")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	g.boxImage = ebiten.NewImageFromImage(i)

	g.emptyImage = ebiten.NewImage(1, 1)
	g.emptyImage.Fill(color.White)
	return g, nil
}

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.lastClicked < 0 {
		g.lastClicked = 20
		x, y := ebiten.CursorPosition()
		box := Box{
			x: float64(x),
			y: float64(y),
			w: boxWidth,
			h: boxHeight,
			r: rand.Float64() * 2 * math.Pi,
		}
		body := cp.NewBody(1.0, 1)
		body.SetPosition(cp.Vector{X: box.x, Y: box.y})
		body.SetAngle(box.r)
		body.UserData = "box"
		shape := cp.NewBox(body, box.w, box.h, 0.0)
		shape.SetFriction(1.0)
		shape.SetElasticity(0.0)
		g.space.AddBody(body)
		g.space.AddShape(shape)
	} else {
		g.lastClicked--
	}

	g.space.Step(1.0 / 60)

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 30})

	op := &ebiten.DrawImageOptions{}

	g.space.EachBody(func(body *cp.Body) {
		switch body.UserData {
		case "box":
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(boxWidth)/2, -float64(boxHeight)/2) // Center the image
			op.GeoM.Rotate(body.Angle())
			op.GeoM.Translate(body.Position().X, body.Position().Y)
			screen.DrawImage(g.boxImage, op)
		case "line":
			op.GeoM.Reset()
			op.GeoM.Scale(g.w, 5)
			op.GeoM.Translate(-g.w/2, -5/2) // Center the image
			op.GeoM.Translate(body.Position().X, body.Position().Y)
			screen.DrawImage(g.emptyImage, op)
		}
	})

	ebitenutil.DebugPrint(screen, "Actual TPS: "+fmt.Sprintf("%.2f", ebiten.ActualTPS()))
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.w), int(g.h)
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Physics Calculation")

	g, err := newGame()
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
