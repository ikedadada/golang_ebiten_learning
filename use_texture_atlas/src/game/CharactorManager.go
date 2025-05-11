package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type CharactorManager interface {
	Draw(screen *ebiten.Image)
	Update()
}

type gridPoint struct {
	gridX, gridY int
}

type charactorManagerImpl struct {
	dm DrawManager
	cp gridPoint
}

func NewCharactorManager(dm DrawManager, gridX, gridY int) CharactorManager {
	return &charactorManagerImpl{
		dm: dm,
		cp: gridPoint{
			gridX: gridX,
			gridY: gridY,
		},
	}
}

func (cm *charactorManagerImpl) Draw(screen *ebiten.Image) {
	// Draw the character at the specified grid position
	cm.dm.DrawImage(screen, 25, cm.cp.gridX, cm.cp.gridY)
	// Draw the character's name or other information
	cm.dm.DrawText(screen, "Character", cm.cp.gridX, cm.cp.gridY-1, text.AlignCenter)
}

func (cm *charactorManagerImpl) Update() {
	// Update the character's position or state here
	// For example, you can move the character based on user input
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		cm.cp.gridY--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		cm.cp.gridY++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		cm.cp.gridX--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		cm.cp.gridX++
	}
	if cm.cp.gridX < 0 {
		cm.cp.gridX = 0
	}
	if cm.cp.gridY < 0 {
		cm.cp.gridY = 0
	}
	if cm.cp.gridX > 9 {
		cm.cp.gridX = 9
	}
	if cm.cp.gridY > 9 {
		cm.cp.gridY = 9
	}
}
