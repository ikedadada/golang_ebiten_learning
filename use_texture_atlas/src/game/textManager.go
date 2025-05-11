package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextManager interface {
	DrawText(screen *ebiten.Image, message string, x, y int, align text.Align)
}

type TextManagerImpl struct {
	fontFace *text.GoTextFace
}

func NewTextManager(fontFace *text.GoTextFace) TextManager {
	return &TextManagerImpl{
		fontFace: fontFace,
	}
}

func (tm *TextManagerImpl) DrawText(screen *ebiten.Image, message string, gridX int, gridY int, align text.Align) {
	op := &text.DrawOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(float64(gridX*GRID_SIZE), float64(gridY*GRID_SIZE))
	op.Filter = ebiten.FilterNearest
	op.LayoutOptions = text.LayoutOptions{
		PrimaryAlign:   align,
		SecondaryAlign: text.AlignCenter,
	}
	text.Draw(screen, message, tm.fontFace, op)
}
