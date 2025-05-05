package main

import "github.com/eihigh/miniten"

var y = 0

func main() {
	miniten.Run(draw)
}

func draw() {
	if miniten.IsClicked() {
		y -= 6
	} else {
		y += 6
	}
	miniten.DrawImage("gopher.png", 0, y)
}
