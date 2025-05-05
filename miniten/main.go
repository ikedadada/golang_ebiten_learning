package main

import (
	"math/rand/v2"

	"github.com/eihigh/miniten"
)

var xs = []int{100, 200, 300}

func main() {
	miniten.Run(draw)
}

func draw() {
	if miniten.IsClicked() {
		xs = append(xs, rand.N(640))
	}
	for i := range xs {
		xs[i] += 1
		miniten.DrawImage("gopher.png", xs[i], 0)
	}
}
