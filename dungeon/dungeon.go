// package dungeon provides dungeon maps.
package dungeon

import (
	"errors"
	"fmt"
	mrand "math/rand"
	"strings"
)

type TileCode int

type Tile struct {
	repr rune
}

func (t Tile) String() string {
	return string(t.repr)
}

func (t Tile) Rune() rune {
	return t.repr
}

var (
	GroundTile = Tile{
		repr: '.',
	}
	WallTile = Tile{
		repr: '#',
	}
	MossTile = Tile{
		repr: '~',
	}
	StumpTile = Tile{
		repr: '=',
	}
	PlayerTile = Tile{
		repr: '@',
	}
	codeToTile = map[TileCode]Tile{
		0: GroundTile,
		1: WallTile,
		2: MossTile,
		3: StumpTile,
	}
	tileToCode = map[rune]TileCode{
		GroundTile.repr: 0,
		WallTile.repr:   1,
		MossTile.repr:   2,
		StumpTile.repr:  3,
	}
)

type Map struct {
	Width  int
	Height int

	tiles []TileCode
}

func (m Map) Dimensions() (int, int) {
	return m.Width, m.Height
}

func (m Map) String() string {
	var mapStr strings.Builder

	mapStr.WriteString(" 0123456789ABCDEF\n")
	for y := 0; y < m.Height; y++ {
		mapStr.WriteString(fmt.Sprintf("%d", y))
		for x := 0; x < m.Width; x++ {
			t := m.GetTile(x, y)
			mapStr.WriteString(t.String())
		}
		mapStr.WriteString("\n")
	}

	return mapStr.String()
}

func (m Map) GetTile(x int, y int) Tile {
	index := x + (y * m.Width)
	return m.GetTileIndex(index)
}

func (m Map) GetTileIndex(i int) Tile {
	c := m.tiles[i]
	return codeToTile[c]
}

type Neighbour struct {
	Tile
	X int
	Y int
}

func (n Neighbour) String() string {
	return fmt.Sprintf("(%3d, %3d) %s", n.X, n.Y, n.Tile.String())
}

func (m Map) GetNeighbours(x int, y int) []Neighbour {
	var results []Neighbour

	toIndex := func(x, y int) (int, error) {
		if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
			return 0, errors.New("out of bounds")
		}
		return x + (y * m.Width), nil
	}

	addIfNotErr := func(x, y int) {
		i, err := toIndex(x, y)
		if err != nil {
			return
		}
		results = append(results, Neighbour{
			Tile: m.GetTileIndex(i),
			X:    x,
			Y:    y,
		})
	}

	addIfNotErr(x-1, y-1)
	addIfNotErr(x, y-1)
	addIfNotErr(x+1, y-1)
	addIfNotErr(x+1, y)
	addIfNotErr(x+1, y+1)
	addIfNotErr(x, y+1)
	addIfNotErr(x-1, y+1)
	addIfNotErr(x-1, y)
	return results
}

func (m Map) GetSurroundingTiles(x int, y int) [8]Neighbour {
	var results [8]Neighbour

	toIndex := func(x, y int) (int, int, int) {
		if x < 0 {
			x = m.Width + x
		}
		if x >= m.Width {
			x = m.Width - x
		}
		if y < 0 {
			y = m.Height + y
		}
		if y >= m.Height {
			y = m.Height - y
		}
		return x + (y * m.Width), x, y
	}

	type position struct {
		index, x, y int
	}

	makePosition := func(x, y int) position {
		index, newX, newY := toIndex(x, y)
		return position{index, newX, newY}
	}

	positions := [8]position{
		makePosition(x-1, y-1),
		makePosition(x, y-1),
		makePosition(x+1, y-1),
		makePosition(x+1, y),
		makePosition(x+1, y+1),
		makePosition(x, y+1),
		makePosition(x-1, y+1),
		makePosition(x-1, y),
	}

	for i, position := range positions {
		results[i] = Neighbour{
			Tile: m.GetTileIndex(position.index),
			X:    position.x,
			Y:    position.y,
		}
	}
	return results
}

func GenerateMap(seed int64, maxX, maxY int) Map {
	maxIndex := maxX * maxY
	mrand.Seed(seed)

	m := Map{
		Width:  maxX,
		Height: maxY,
		tiles:  make([]TileCode, maxIndex),
	}

	for i := 0; i < maxIndex; i++ {
		tile := WallTile
		roll := mrand.Intn(100)
		if roll >= 45 {
			tile = GroundTile
		}
		m.tiles[i] = tileToCode[tile.Rune()]
	}

	return m
}
