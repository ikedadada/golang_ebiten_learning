package main

import "github.com/eihigh/miniten"

func main() {
	miniten.Run(draw)
}

func draw() {
	miniten.DrawImage("gopher.png", 0, 0)
}
