package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageManager interface {
	DrawImage(screen *ebiten.Image, index int, x, y int)
	DrawLayor(screen *ebiten.Image, layor [][]int, x, y int)
}

type imageManagerImpl struct {
	tilemap *tilemap
}

type tilemap struct {
	tileset  *ebiten.Image
	tilesetW int
	tilesetH int
}

func NewImageManager(tileset *ebiten.Image) ImageManager {
	tilesetW, tilesetH := tileset.Bounds().Dx(), tileset.Bounds().Dy()
	tilemap := &tilemap{
		tileset:  tileset,
		tilesetW: tilesetW,
		tilesetH: tilesetH,
	}

	return &imageManagerImpl{
		tilemap: tilemap,
	}
}

// index=0 is empty
// index=1 is tile1
func (im *imageManagerImpl) DrawImage(screen *ebiten.Image, index int, gridX int, gridY int) {
	if index <= 0 || index >= im.tilemap.tilesetW*im.tilemap.tilesetH {
		return
	}
	sx := (index - 1) % (im.tilemap.tilesetW / GRID_SIZE)
	sy := (index - 1) / (im.tilemap.tilesetW / GRID_SIZE)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(float64(gridX*GRID_SIZE), float64(gridY*GRID_SIZE))
	screen.DrawImage(im.tilemap.tileset.SubImage(image.Rect(sx*GRID_SIZE, sy*GRID_SIZE, (sx+1)*GRID_SIZE, (sy+1)*GRID_SIZE)).(*ebiten.Image), op)
}

func (im *imageManagerImpl) DrawLayor(screen *ebiten.Image, layor [][]int, x, y int) {
	for i := 0; i < len(layor); i++ {
		for j := 0; j < len(layor[i]); j++ {
			im.DrawImage(screen, layor[i][j], x+j, y+i)
		}
	}
}
