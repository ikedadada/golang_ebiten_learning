package main

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"

	"use_texture_atlas/src/game"
)

//go:embed assets/*
var fsys embed.FS

func main() {
	game, err := game.NewGame(fsys)
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
