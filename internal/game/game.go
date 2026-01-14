package game

import (
	"image"

	"github.com/danmcfan/pacman/internal/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 224
	screenHeight = 288

	tileSize = 8
)

type Game struct{}

func New() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(sprite.Maze, nil)

	s := sprite.Pacman.SubImage(image.Rect(0, 0, tileSize*2, tileSize*2)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(13*tileSize+1, 25.5*tileSize)
	screen.DrawImage(s, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
