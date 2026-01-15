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

var Banner *ebiten.Image
var Digits *ebiten.Image
var Fruit *ebiten.Image
var Maze *ebiten.Image
var Pacman *ebiten.Image
var Pellet *ebiten.Image
var Ready *ebiten.Image

func init() {
	Banner = newImage("banner.png")
	Digits = newImage("digits.png")
	Fruit = newImage("fruit.png")
	Maze = newImage("maze.png")
	Pacman = newImage("pacman.png")
	Pellet = newImage("pellet.png")
	Ready = newImage("ready.png")
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
