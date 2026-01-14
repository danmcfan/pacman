package main

import (
	"log"

	"github.com/danmcfan/pacman/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Pac-Man")
	game := game.New()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
