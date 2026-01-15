package game

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/danmcfan/pacman/internal/maze"
	"github.com/danmcfan/pacman/internal/sprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 224
	screenHeight = 288

	tileSize  = 8
	tileWidth = screenWidth / tileSize

	frameCount = 12
)

type Direction Position

var (
	DirectionUp    = Direction{X: 0, Y: -1}
	DirectionDown  = Direction{X: 0, Y: 1}
	DirectionLeft  = Direction{X: -1, Y: 0}
	DirectionRight = Direction{X: 1, Y: 0}
)

type Position struct {
	X, Y int
}

func (p Position) Screen() Position {
	return Position{
		X: p.X*tileSize + tileSize/2,
		Y: p.Y*tileSize + tileSize/2,
	}
}

func (p Position) Tile() Position {
	return Position{
		X: p.X / tileSize,
		Y: p.Y / tileSize,
	}
}

func (p Position) Add(d Direction) Position {
	return Position{
		X: p.X + d.X,
		Y: p.Y + d.Y,
	}
}

type Pacman struct {
	Position  Position
	Direction Direction
	Frame     int
}

type Game struct {
	EditMode  bool
	EditDot   bool
	HoverTile Position
	Pacman    Pacman
	Score     int
	HighScore int
}

func New() *Game {
	return &Game{
		Pacman: Pacman{
			Position:  Position{X: 14 * tileSize, Y: 26.5 * tileSize},
			Direction: DirectionLeft,
		},
	}
}

func (g *Game) Update() error {
	if (ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)) && inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.EditMode = !g.EditMode
	}

	// handle keyboard input for pacman movement
	var direction Direction
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		direction = DirectionUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		direction = DirectionDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		direction = DirectionLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		direction = DirectionRight
	}

	if direction.X != 0 || direction.Y != 0 {
		currentTile := g.Pacman.Position.Tile()
		nextTile := currentTile.Add(direction)
		if maze.Get(nextTile.X, nextTile.Y) != maze.Wall {
			if g.Pacman.Direction != direction {
				centerPosition := currentTile.Screen()
				if direction == DirectionUp || direction == DirectionDown {
					g.Pacman.Position.X = centerPosition.X
				}
				if direction == DirectionLeft || direction == DirectionRight {
					g.Pacman.Position.Y = centerPosition.Y
				}
			}
			g.Pacman.Direction = direction
		}
	}

	currentPosition := g.Pacman.Position
	currentTile := currentPosition.Tile()
	nextTile := currentTile.Add(g.Pacman.Direction)

	if maze.Get(nextTile.X, nextTile.Y) == maze.Wall {
		centerPosition := currentTile.Screen()
		if currentPosition != centerPosition {
			g.Pacman.Position = currentPosition.Add(g.Pacman.Direction)
		}
	} else {
		g.Pacman.Position = currentPosition.Add(g.Pacman.Direction)
	}

	if g.Pacman.Position.Tile().Y == 17 {
		if g.Pacman.Position.X < 0 {
			g.Pacman.Position.X = screenWidth
		}
		if g.Pacman.Position.X > screenWidth {
			g.Pacman.Position.X = 0
		}
	}

	if currentPosition != g.Pacman.Position {
		g.Pacman.Frame = (g.Pacman.Frame + 1) % frameCount
	} else {
		if g.Pacman.Frame/4 == 0 {
			g.Pacman.Frame = frameCount / 4
		}
	}

	currentPosition = g.Pacman.Position
	currentTile = g.Pacman.Position.Tile()
	centerPosition := currentTile.Screen()
	if math.Abs(float64(centerPosition.X-g.Pacman.Position.X)) <= 2 && math.Abs(float64(centerPosition.Y-g.Pacman.Position.Y)) <= 2 {
		tile := maze.Get(currentTile.X, currentTile.Y)
		if tile == maze.Dot {
			g.Score += 10
		}
		if tile == maze.Power {
			g.Score += 50
		}
		maze.Set(currentTile.X, currentTile.Y, maze.Empty)
	}

	if g.EditMode {
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.EditDot = !g.EditDot
		}

		x, y := ebiten.CursorPosition()
		cursor := Position{X: x, Y: y}
		g.HoverTile = cursor.Tile()

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

	frame := g.Pacman.Frame / (frameCount / 4)
	var yIndex int
	switch g.Pacman.Direction {
	case DirectionUp:
		yIndex = 0
	case DirectionDown:
		yIndex = 1
	case DirectionLeft:
		yIndex = 2
	case DirectionRight:
		yIndex = 3
	}

	s := sprite.Pacman.SubImage(image.Rect(tileSize*2*frame, tileSize*2*yIndex, tileSize*2*frame+tileSize*2, tileSize*2*yIndex+tileSize*2)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Pacman.Position.X-tileSize), float64(g.Pacman.Position.Y-tileSize))
	screen.DrawImage(s, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(24, 0)
	screen.DrawImage(sprite.Banner, op)

	scoreText := fmt.Sprintf("%02d", g.Score)
	for i, char := range scoreText {
		s := sprite.Digits.SubImage(image.Rect(int(char-'0')*8, 0, int(char-'0')*8+8, 8)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(5*tileSize+i*8), float64(1*tileSize))
		screen.DrawImage(s, op)
	}

	if g.HighScore == 0 {
		zero := sprite.Digits.SubImage(image.Rect(0, 0, 8, 8)).(*ebiten.Image)
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(15*tileSize, 1*tileSize)
		screen.DrawImage(zero, op)

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(16*tileSize, 1*tileSize)
		screen.DrawImage(zero, op)
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(11*tileSize, 20*tileSize)
	screen.DrawImage(sprite.Ready, op)

	for i := range 2 {
		s := sprite.Fruit.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(16+i*16), float64(screenHeight-16))
		screen.DrawImage(s, op)
	}

	s = sprite.Fruit.SubImage(image.Rect(16, 0, 32, 16)).(*ebiten.Image)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screenWidth-32), float64(screenHeight-16))
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
