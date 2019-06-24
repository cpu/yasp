// package game provides the overall game elements and structure.
package game

import (
	"fmt"

	"github.com/cpu/yasp/dungeon"
	"github.com/cpu/yasp/view"
	"github.com/gdamore/tcell"
)

type Player struct {
	x     int
	y     int
	Style view.Style
}

func (p Player) String() string {
	return "@"
}

func (p Player) Pos() (int, int) {
	return p.x, p.y
}

func (p Player) X() int {
	return p.x
}

func (p Player) Y() int {
	return p.y
}

// MoveTo changes the Player's position to the given coordinates. The *previous*
// Player position is returned. Note that there is no error checking done on the
// x, y coordinates with respect to some other game state (e.g. the map). It's
// assumed the provided values have been vetted.
func (p *Player) MoveTo(x, y int) (int, int) {
	oldX, oldY := p.Pos()
	p.x = x
	p.y = y
	return oldX, oldY
}

type State struct {
	Debug   bool
	P       Player
	Map     dungeon.Map
	Display *view.Display
}

func NewGame() State {
	return State{
		P: Player{
			Style: view.DefaultStyle,
			x:     10,
			y:     10,
		},
		Map: dungeon.One,
	}
}

func (state *State) SetDisplay(d *view.Display) {
	state.Display = d
}

func (state *State) PrintDebug() {
	if state.Debug {
		pX, pY := state.P.Pos()
		dbg := fmt.Sprintf("p x: %d y: %d", pX, pY)
		_, maxY := state.Display.Size()
		state.Display.PrintFixed(0, maxY-1, view.DefaultStyle, dbg)
	}
}

func (state *State) HandleInput(ev *tcell.EventKey) {
	maxX, maxY := state.Display.Size()
	pX, pY := state.P.Pos()

	switch ev.Key() {
	case tcell.KeyRight:
		if pX+1 < maxX {
			state.P.MoveTo(pX+1, pY)
		}
	case tcell.KeyLeft:
		if pX-1 >= 0 {
			state.P.MoveTo(pX-1, pY)
		}
	case tcell.KeyUp:
		if pY-1 >= 0 {
			state.P.MoveTo(pX, pY-1)
		}
	case tcell.KeyDown:
		if pY+1 < maxY {
			state.P.MoveTo(pX, pY+1)
		}
	}
}

func (state *State) Tick() {
	state.Display.Clear()

	for y := 0; y < state.Map.Height; y++ {
		for x := 0; x < state.Map.Width; x++ {
			index := x + (y * state.Map.Width)
			tile := dungeon.LookupTile(state.Map.Tiles[index])
			state.Display.PrintFixed(x, y, tile.Style, tile.String())
		}
	}

	pX, pY := state.P.Pos()
	state.Display.PrintFixed(pX, pY, state.P.Style, "@")
	state.Display.Show()
}
