// package dungeon provides dungeon maps.
package dungeon

import (
	"strings"

	"github.com/cpu/yasp/view"
)

type TileCode int

type Tile struct {
	Style view.Style
	repr  string
}

func (t Tile) String() string {
	return t.repr
}

type Map struct {
	Width  int
	Height int

	Tiles []TileCode
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
