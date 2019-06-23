// package assets handles asset loading and base types.
package assets

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"

	// load support for png images
	_ "image/png"

	"github.com/faiface/pixel"
)

const (
	TileWidth  float64 = 16
	TileHeight float64 = 16
)

type Tile struct {
	Name string
	Rect pixel.Rect
}

type Tilemap struct {
	Name    string
	Picture pixel.Picture
	Tiles   []Tile
}

func (tm *Tilemap) LoadTiles() error {
	var count int
	for x := tm.Picture.Bounds().Min.X; x < tm.Picture.Bounds().Max.X; x += TileWidth {
		for y := tm.Picture.Bounds().Min.Y; y < tm.Picture.Bounds().Max.Y; y += TileHeight {
			tm.Tiles = append(tm.Tiles, Tile{
				Name: fmt.Sprintf("%s_%d", tm.Name, count),
				Rect: pixel.R(x, y, x+TileWidth, y+TileHeight),
			})
			count++
		}
	}
	return nil
}

func LoadTilemapFile(name string, pngPath string) (*Tilemap, error) {
	pngBytes, err := ioutil.ReadFile(pngPath)
	if err != nil {
		return nil, err
	}
	return LoadTilemap(name, pngBytes)
}

func LoadTilemap(name string, pngData []byte) (*Tilemap, error) {
	img, _, err := image.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, err
	}

	tm := &Tilemap{
		Name:    name,
		Picture: pixel.PictureDataFromImage(img),
	}

	if err := tm.LoadTiles(); err != nil {
		return nil, err
	}

	return tm, nil
}
