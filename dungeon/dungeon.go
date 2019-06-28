// package dungeon provides dungeon maps.
package dungeon

import (
	"strings"
)

type TileCode int

type Tile struct {
	repr string
}

func (t Tile) String() string {
	return t.repr
}

var (
	PlayerTile = Tile{
		repr: "@",
	}
)

type Map struct {
	Width  int
	Height int

	Tiles []TileCode
}

func (m Map) Dimensions() (int, int) {
	return m.Width, m.Height
}

func (m Map) String() string {
	var mapStr strings.Builder

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			code := m.Tiles[x+(y*m.Width)]
			chr := TileMap[code]
			mapStr.WriteString(chr.String())
		}
		mapStr.WriteString("\n")
	}

	return mapStr.String()
}

func LookupTile(c TileCode) Tile {
	return TileMap[c]
}
