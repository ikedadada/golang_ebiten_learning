package main

import "github.com/eihigh/miniten"

var (
	y  = 0.0
	vy = 0.0
)

func main() {
	miniten.Run(draw)
}

func draw() {
	if miniten.IsClicked() {
		vy = -10
	}
	vy += 0.5
	y += vy

	if 360 < y {
		y = 360
		vy = 0
	}
	miniten.DrawImage("gopher.png", 0, int(y))
}
