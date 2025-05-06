package main

import (
	"embed"
	"image"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type game struct {
	ticks     int
	bg        *ebiten.Image
	person    *ebiten.Image
	cat       *ebiten.Image
	messageBG *ebiten.Image
	fontFace  *text.GoTextFace
}

//go:embed assets/*
var fsys embed.FS

func loadImage(path string) (*ebiten.Image, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(i), nil
}

func loadFont(path string) (*text.GoTextFaceSource, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		return nil, err
	}
	return src, nil
}

func IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(inpututil.AppendJustPressedTouchIDs(nil)) != 0 {
		return true
	}
	return false
}

func newGame() (*game, error) {
	g := &game{}
	img, err := loadImage("assets/bg.jpg")
	if err != nil {
		return nil, err
	}
	g.bg = img
	img, err = loadImage("assets/person.png")
	if err != nil {
		return nil, err
	}
	g.person = img
	img, err = loadImage("assets/cat.png")
	if err != nil {
		return nil, err
	}
	g.cat = img

	img, err = loadImage("assets/message-bg.png")
	if err != nil {
		return nil, err
	}
	g.messageBG = img

	fontSrc, err := loadFont("assets/font/NotoSansJP-VariableFont_wght.ttf")
	if err != nil {
		return nil, err
	}
	g.fontFace = &text.GoTextFace{
		Source: fontSrc,
		Size:   30,
	}
	return g, nil
}

func (g *game) Update() error {
	g.ticks++
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

	textop := &text.DrawOptions{}
	textop.ColorScale.Scale(0, 0, 0, 1)
	textop.LineSpacing = 30 * 1.5
	glyphs := text.AppendGlyphs(nil, "吾輩は猫である。名前はまだない。", g.fontFace, &textop.LayoutOptions)

	count := g.ticks / 5
	count = count % len(glyphs)
	for _, g := range glyphs[:count] {
		textop.GeoM.Reset()
		textop.GeoM.Translate(10, float64(screen.Bounds().Dy()*2/3))
		textop.GeoM.Translate(float64(g.X), float64(g.Y))
		screen.DrawImage(g.Image, &textop.DrawImageOptions)
	}

	if IsClicked() {
		g.ticks = 0
	}

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
