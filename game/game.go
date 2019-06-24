package game

import (
	"errors"
	"math"
	"time"

	"github.com/cpu/yasp"
	"github.com/cpu/yasp/assets"
	"github.com/cpu/yasp/dungeon"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	tilemapFile          = "assets/data/tileset.png"
	characterTimemapFile = "assets/data/characters.png"

	tileScale         = 4.0
	defaultPlayerTile = 13
)

type player struct {
	sprite *pixel.Sprite
	matrix pixel.Matrix
}

type Game struct {
	config       *yasp.Config
	windowConfig pixelgl.WindowConfig
	win          *pixelgl.Window

	tilemap          *assets.Tilemap
	characterTilemap *assets.Tilemap

	lastFrame time.Time

	camPos       pixel.Vec
	camSpeed     float64
	camZoom      float64
	camZoomSpeed float64

	p player

	d dungeon.Map

	dBatch *pixel.Batch
}

func New(c *yasp.Config) (*Game, error) {
	if c == nil {
		return nil, errors.New("nil config")
	}

	wc := pixelgl.WindowConfig{
		Title:  "Y A S P",
		Bounds: pixel.R(0, 0, float64(c.WinWidth), float64(c.WinHeight)),
		VSync:  c.VSync,
	}

	tilemap, err := assets.LoadTilemapFile("tilemap", tilemapFile)
	if err != nil {
		return nil, err
	}

	dBatch := pixel.NewBatch(&pixel.TrianglesData{}, tilemap.Picture)

	characterTM, err := assets.LoadTilemapFile("characters", characterTimemapFile)
	if err != nil {
		return nil, err
	}

	p := player{
		sprite: pixel.NewSprite(
			characterTM.Picture,
			characterTM.Tiles[defaultPlayerTile].Rect),
		matrix: pixel.IM.Scaled(pixel.ZV, tileScale),
	}

	d := dungeon.One

	return &Game{
		config:           c,
		windowConfig:     wc,
		tilemap:          tilemap,
		characterTilemap: characterTM,

		lastFrame: time.Now(),

		camPos:       pixel.ZV,
		camSpeed:     500.0,
		camZoom:      1.0,
		camZoomSpeed: 1.2,

		p:      p,
		d:      d,
		dBatch: dBatch,
	}, nil
}

func (g *Game) Run() {
	win, err := pixelgl.NewWindow(g.windowConfig)
	if err != nil {
		panic(err)
	}
	g.win = win

	wall := pixel.NewSprite(
		g.tilemap.Picture,
		g.tilemap.Tiles[0].Rect)

	for !g.win.Closed() {
		dt := time.Since(g.lastFrame).Seconds()
		g.lastFrame = time.Now()

		cam := pixel.IM.Scaled(g.camPos, g.camZoom).
			Moved(g.win.Bounds().Center().Sub(g.camPos))
		g.win.SetMatrix(cam)

		g.readButtons(dt)

		g.dBatch.Clear()

		m := pixel.IM.Scaled(pixel.ZV, tileScale)
		for y := 0; y < g.d.Height; y++ {
			for x := 0; y < g.d.Width; x++ {
				wall.Draw(g.dBatch, m)
				m = m.Moved(pixel.V(0, float64(x)*assets.TileWidth))
			}
			m = m.Moved(pixel.V(0, float64(y)*assets.TileHeight))
		}
		g.dBatch.Draw(g.win)

		//g.win.Clear(colornames.Black)

		//g.p.sprite.Draw(g.win, g.p.matrix)

		g.win.Update()
	}
}

func (g *Game) readButtons(dt float64) {
	/* Camera movement */
	if g.win.Pressed(pixelgl.KeyA) {
		g.camPos.X -= g.camSpeed * dt
	}
	if g.win.Pressed(pixelgl.KeyD) {
		g.camPos.X += g.camSpeed * dt
	}
	if g.win.Pressed(pixelgl.KeyS) {
		g.camPos.Y -= g.camSpeed * dt
	}
	if g.win.Pressed(pixelgl.KeyW) {
		g.camPos.Y += g.camSpeed * dt
	}
	g.camZoom *= math.Pow(g.camZoomSpeed, g.win.MouseScroll().Y)
}
