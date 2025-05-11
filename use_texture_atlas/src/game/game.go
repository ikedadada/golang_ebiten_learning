package game

import (
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	GRID_SIZE = 16
)

type game struct {
	dm     DrawManager
	width  int
	height int
}

func NewGame(fsys embed.FS) (ebiten.Game, error) {

	// Create a new game instance with the specified width and height
	g := &game{
		width:  160, // 16 * 10
		height: 160, // 16 * 10
	}
	ebiten.SetWindowSize(g.width*4, g.height*4)

	// Load the font
	f, err := fsys.Open("assets/misaki_gothic_2nd.ttf")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	font, err := text.NewGoTextFaceSource(f)
	if err != nil {
		return nil, err
	}
	fontFace := &text.GoTextFace{
		Source: font,
		Size:   GRID_SIZE,
	}

	// Load the images
	f, err = fsys.Open("assets/tilemap_packed.png")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	tileset := ebiten.NewImageFromImage(img)

	g.dm = NewDrawManager(fontFace, tileset)

	return g, nil
}

func (g *game) Update() error {
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	layor := [][]int{
		{1, 2, 2, 2, 2, 2, 2, 2, 2, 3},
		{28, 29, 29, 29, 29, 29, 29, 29, 29, 30},
		{28, 29, 29, 29, 29, 29, 29, 29, 29, 30},
		{28, 29, 29, 29, 29, 29, 29, 29, 29, 30},
		{28, 29, 29, 29, 29, 29, 29, 29, 29, 30},
		{28, 29, 29, 29, 29, 29, 29, 29, 29, 30},
		{28, 29, 29, 29, 29, 29, 29, 29, 29, 30},
		{28, 29, 29, 29, 29, 29, 29, 29, 29, 30},
		{28, 29, 29, 29, 29, 29, 29, 29, 29, 30},
		{55, 56, 56, 56, 56, 56, 56, 56, 56, 57},
	}
	g.dm.DrawLayor(screen, layor, 0, 0)

	g.dm.DrawText(screen, "Hello, Ebiten!", g.width/(GRID_SIZE*2), g.height/(GRID_SIZE*2), text.AlignCenter)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
