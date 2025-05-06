package main

import (
	"embed"
	"image"
	_ "image/jpeg"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type chara struct {
	name string
}

type game struct {
	images   map[string]*ebiten.Image
	fontFace *text.GoTextFace

	scenario    []string
	progress    int
	message     string
	leftCamera  chara
	rightCamera chara
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
	imageNames := []string{
		"bg.jpg",
		"person.png",
		"cat.png",
		"message-bg.png",
	}
	g.images = make(map[string]*ebiten.Image, len(imageNames))
	for _, name := range imageNames {
		img, err := loadImage("assets/" + name)
		if err != nil {
			return nil, err
		}
		g.images[name] = img
	}

	fontSrc, err := loadFont("assets/font/NotoSansJP-VariableFont_wght.ttf")
	if err != nil {
		return nil, err
	}
	g.fontFace = &text.GoTextFace{
		Source: fontSrc,
		Size:   30,
	}

	g.scenario = []string{
		"rightChara=cat.png",
		"吾輩は猫である。",
		"名前はまだない。",
		"leftChara=person.png",
		"吾輩はここで始めて人間というものを見た。",
	}
	return g, nil
}

func (g *game) Update() error {
	if IsClicked() {
		s := g.scenario[g.progress]
		if g.progress < len(g.scenario)-1 {
			g.progress++
		}

		before, after, found := strings.Cut(s, "=")
		if found {
			switch before {
			case "leftChara":
				g.leftCamera.name = after
			case "rightChara":
				g.rightCamera.name = after
			default:
				g.message = s
			}
		} else {
			g.message = s
		}
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screen.Bounds().Dx())/float64(g.images["bg.jpg"].Bounds().Dx()), float64(screen.Bounds().Dy())/float64(g.images["bg.jpg"].Bounds().Dy()))
	screen.DrawImage(g.images["bg.jpg"], op)

	if g.leftCamera.name != "" {
		op = &ebiten.DrawImageOptions{}
		screen.DrawImage(g.images[g.leftCamera.name], op)
	}

	if g.rightCamera.name != "" {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(screen.Bounds().Dx()), 0)
		screen.DrawImage(g.images[g.rightCamera.name], op)
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screen.Bounds().Dx())/float64(g.images["message-bg.png"].Bounds().Dx()), float64(screen.Bounds().Dy()/3)/float64(g.images["message-bg.png"].Bounds().Dy()))
	op.GeoM.Translate(0, float64(screen.Bounds().Dy()*2/3))
	op.ColorScale.ScaleAlpha(0.5)
	screen.DrawImage(g.images["message-bg.png"], op)

	textop := &text.DrawOptions{}
	textop.GeoM.Translate(10, float64(screen.Bounds().Dy()*2/3))
	textop.ColorScale.Scale(0, 0, 0, 1)
	textop.LineSpacing = 30 * 1.5
	text.Draw(screen, g.message, g.fontFace, textop)
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
