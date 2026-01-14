package maze

import (
	"embed"
	"fmt"
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

func Get(x, y int) Tile {
	return Tile(Data[y][x])
}

func Set(x, y int, tile Tile) {
	Data[y] = Data[y][:x] + string(tile) + Data[y][x+1:]
}

func Print() {
	for _, row := range Data {
		fmt.Println(row)
	}
}
