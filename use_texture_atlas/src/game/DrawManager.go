package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type DrawManager interface {
	DrawText(screen *ebiten.Image, message string, x, y int, align text.Align)
	DrawImage(screen *ebiten.Image, index int, x, y int)
	DrawLayor(screen *ebiten.Image, layor [][]int, x, y int)
}

type drawManagerImpl struct {
	fontFace *text.GoTextFace
	tilemap  *tilemap
}
type tilemap struct {
	tileset  *ebiten.Image
	tilesetW int
	tilesetH int
}

func NewDrawManager(fontFace *text.GoTextFace, tileset *ebiten.Image) DrawManager {
	tilesetW, tilesetH := tileset.Bounds().Dx(), tileset.Bounds().Dy()
	tilemap := &tilemap{
		tileset:  tileset,
		tilesetW: tilesetW,
		tilesetH: tilesetH,
	}

	return &drawManagerImpl{
		fontFace: fontFace,
		tilemap:  tilemap,
	}
}

func (dm *drawManagerImpl) DrawText(screen *ebiten.Image, message string, gridX int, gridY int, align text.Align) {
	op := &text.DrawOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(float64(gridX*GRID_SIZE), float64(gridY*GRID_SIZE))
	op.Filter = ebiten.FilterNearest
	op.LayoutOptions = text.LayoutOptions{
		PrimaryAlign:   align,
		SecondaryAlign: text.AlignCenter,
	}
	text.Draw(screen, message, dm.fontFace, op)
}
func (dm *drawManagerImpl) DrawImage(screen *ebiten.Image, index int, gridX int, gridY int) {
	if index <= 0 || index >= dm.tilemap.tilesetW*dm.tilemap.tilesetH {
		return
	}
	sx := (index - 1) % (dm.tilemap.tilesetW / GRID_SIZE)
	sy := (index - 1) / (dm.tilemap.tilesetW / GRID_SIZE)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(float64(gridX*GRID_SIZE), float64(gridY*GRID_SIZE))
	screen.DrawImage(dm.tilemap.tileset.SubImage(image.Rect(sx*GRID_SIZE, sy*GRID_SIZE, (sx+1)*GRID_SIZE, (sy+1)*GRID_SIZE)).(*ebiten.Image), op)
}

func (dm *drawManagerImpl) DrawLayor(screen *ebiten.Image, layor [][]int, x, y int) {
	for i := 0; i < len(layor); i++ {
		for j := 0; j < len(layor[i]); j++ {
			dm.DrawImage(screen, layor[i][j], x+j, y+i)
		}
	}
}
