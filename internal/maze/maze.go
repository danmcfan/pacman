package maze

import (
	"embed"
	"strings"
)

type Tile rune

const (
	Wall  Tile = '#'
	Empty Tile = ' '
	Dot   Tile = '.'
	Power Tile = 'o'
)

//go:embed data
var data embed.FS

var Data []string

func init() {
	data, err := data.ReadFile("data/maze.txt")
	if err != nil {
		panic(err)
	}
	Data = strings.Split(string(data), "\n")
}

func Get(x, y int) rune {
	return rune(Data[y][x])
}
