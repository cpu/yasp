// package dungeon provides dungeon maps.
package dungeon

import (
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
	tileMap = map[TileCode]Tile{
		0: GroundTile,
		1: WallTile,
		2: MossTile,
		3: StumpTile,
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

	for y := 0; y < m.Height; y++ {
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
	c := m.tiles[index]
	return tileMap[c]
}
