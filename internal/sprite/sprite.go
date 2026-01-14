package sprite

import (
	"bytes"
	"embed"
	"image"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed images
var assets embed.FS

var Maze *ebiten.Image
var Pacman *ebiten.Image
var Pellet *ebiten.Image

func init() {
	Maze = newImage("maze.png")
	Pacman = newImage("pacman.png")
	Pellet = newImage("pellet.png")
}

func newImage(name string) *ebiten.Image {
	data, err := assets.ReadFile("images/" + name)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}
