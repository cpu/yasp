// package game provides the overall game elements and structure.
package game

import (
	"fmt"

	"github.com/cpu/yasp/dungeon"
)

type Player struct {
	x int
	y int
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

// Clamp restricts the Player's position to [0, maxX) and [0, maxY). It returns
// the player's updated position after clamping.
func (p *Player) Clamp(maxX, maxY int) (int, int) {
	if p.x < 0 {
		p.x = 0
	}
	if p.y < 0 {
		p.y = 0
	}
	if p.x > maxX-1 {
		p.x = maxX - 1
	}
	if p.y > maxY-1 {
		p.y = maxY - 1
	}
	return p.x, p.y
}

func (p *Player) String() string {
	return fmt.Sprintf("Player: x=%4d y=%4d", p.x, p.y)
}

type State struct {
	Debug bool
	P     Player
	Map   dungeon.Map
}

func NewGame() State {
	return State{
		P: Player{
			x: 10,
			y: 10,
		},
		Map: dungeon.One,
	}
}

func (state *State) PrintDebug() {
	/*
		if state.Debug {
			pX, pY := state.P.Pos()
			dbg := fmt.Sprintf("p x: %d y: %d", pX, pY)
			_, maxY := 100, 100
			state.Display.PrintFixed(0, maxY-1, view.DefaultStyle, dbg)
		}
	*/
}

/*
func (state *State) HandleInput(ev view.InputEvent) {
	maxX, maxY := 100, 100
	pX, pY := state.P.Pos()

	switch {
	case ev == view.InputDebug:
		state.Debug = !state.Debug
	case ev == view.InputKeyRight:
		if pX+1 < maxX {
			state.P.MoveTo(pX+1, pY)
		}
	case ev == view.InputKeyLeft:
		if pX-1 >= 0 {
			state.P.MoveTo(pX-1, pY)
		}
	case ev == view.InputKeyUp:
		if pY-1 >= 0 {
			state.P.MoveTo(pX, pY-1)
		}
	case ev == view.InputKeyDown:
		if pY+1 < maxY {
			state.P.MoveTo(pX, pY+1)
		}
	}
}

func (state *State) Tick() {
	for y := 0; y < state.Map.Height; y++ {
		for x := 0; x < state.Map.Width; x++ {
			index := x + (y * state.Map.Width)
			tile := dungeon.LookupTile(state.Map.Tiles[index])
			state.Display.PrintFixed(x, y, tile.Style, tile.String())
		}
	}

	state.PrintDebug()

	pX, pY := state.P.Pos()
	state.Display.PrintFixed(pX, pY, state.P.Style, "@")
}
*/
