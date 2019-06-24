// package dungeon provides dungeon maps.
package dungeon

type TileCode int

type Map struct {
	Width  int
	Height int

	Tiles []TileCode
}

var (
	One = Map{
		Width:  5,
		Height: 5,
		Tiles: []TileCode{
			1, 2, 2, 2, 1,
			3, 0, 0, 0, 3,
			3, 0, 0, 0, 3,
			3, 0, 0, 0, 3,
			1, 2, 2, 2, 1,
		},
	}
)
