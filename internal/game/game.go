package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/danmcfan/pacman/internal/maze"
	"github.com/danmcfan/pacman/internal/sprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 224
	screenHeight = 288

	tileSize   = 8
	tileWidth  = screenWidth / tileSize
	tileHeight = screenHeight / tileSize
)

type TilePosition struct {
	X, Y int
}

type Game struct {
	EditMode  bool
	EditDot   bool
	HoverTile TilePosition
}

func New() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	if (ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)) && inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.EditMode = !g.EditMode
	}

	if g.EditMode {
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.EditDot = !g.EditDot
		}

		x, y := ebiten.CursorPosition()
		g.HoverTile = TilePosition{
			X: x / tileSize,
			Y: y / tileSize,
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if g.EditDot {
				tile := maze.Get(g.HoverTile.X, g.HoverTile.Y)
				fmt.Println(tile)
				if tile != maze.Wall {
					var newTile maze.Tile
					if tile == maze.Dot {
						newTile = maze.Empty
					} else {
						newTile = maze.Dot
					}
					maze.Set(g.HoverTile.X, g.HoverTile.Y, newTile)
					maze.Set(tileWidth-g.HoverTile.X-1, g.HoverTile.Y, newTile)
				}
			} else {
				maze.Set(g.HoverTile.X, g.HoverTile.Y, maze.Wall)
				maze.Set(tileWidth-g.HoverTile.X-1, g.HoverTile.Y, maze.Wall)
			}
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			if g.EditDot {
				tile := maze.Get(g.HoverTile.X, g.HoverTile.Y)
				fmt.Println(tile)
				if tile != maze.Wall {
					var newTile maze.Tile
					if tile == maze.Power {
						newTile = maze.Empty
					} else {
						newTile = maze.Power
					}
					maze.Set(g.HoverTile.X, g.HoverTile.Y, newTile)
					maze.Set(tileWidth-g.HoverTile.X-1, g.HoverTile.Y, newTile)
				}
			} else {
				maze.Set(g.HoverTile.X, g.HoverTile.Y, maze.Empty)
				maze.Set(tileWidth-g.HoverTile.X-1, g.HoverTile.Y, maze.Empty)
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			maze.Print()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(sprite.Maze, nil)

	for y := 0; y < len(maze.Data); y++ {
		for x := 0; x < len(maze.Data[y]); x++ {
			if maze.Get(x, y) == maze.Dot {
				s := sprite.Pellet.SubImage(image.Rect(0, 0, tileSize, tileSize)).(*ebiten.Image)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
				screen.DrawImage(s, op)
			}
			if maze.Get(x, y) == maze.Power {
				s := sprite.Pellet.SubImage(image.Rect(tileSize, 0, tileSize*2, tileSize)).(*ebiten.Image)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
				screen.DrawImage(s, op)
			}
		}
	}

	s := sprite.Pacman.SubImage(image.Rect(0, 0, tileSize*2, tileSize*2)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(13*tileSize+1, 25.5*tileSize)
	screen.DrawImage(s, op)

	if g.EditMode {
		for y := 0; y < len(maze.Data); y++ {
			for x := 0; x < len(maze.Data[y]); x++ {
				c := color.RGBA{0, 0, 0, 0}
				switch maze.Get(x, y) {
				case maze.Wall:
					c = color.RGBA{255, 0, 0, 64}
				case maze.Dot:
					c = color.RGBA{255, 255, 255, 64}
				case maze.Power:
					c = color.RGBA{255, 255, 0, 64}
				}
				vector.FillRect(screen, float32(x*tileSize), float32(y*tileSize), float32(tileSize), float32(tileSize), c, false)
			}
		}

		var c color.RGBA
		if g.EditDot {
			c = color.RGBA{255, 255, 255, 255}
		} else {
			c = color.RGBA{255, 0, 0, 255}
		}
		vector.StrokeRect(
			screen,
			float32(g.HoverTile.X*tileSize),
			float32(g.HoverTile.Y*tileSize),
			float32(tileSize),
			float32(tileSize),
			1,
			c,
			false,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
