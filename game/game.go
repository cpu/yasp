// package game provides the overall game elements and structure.
package game

import (
	"fmt"

	"github.com/cpu/yasp/dungeon"
	"github.com/cpu/yasp/game/events"
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
// x, y coordinates with respect to some other game state (e.g. the map). Use
// clamp with the bounding dimensions to ensure the player position is
// consistent with the game world.
func (p *Player) MoveTo(x, y int) (int, int) {
	oldX, oldY := p.Pos()
	p.x = x
	p.y = y
	return oldX, oldY
}

func (p *Player) Move(offX, offY int) (int, int) {
	oldX, oldY := p.Pos()
	p.x = oldX + offX
	p.y = oldY + offY
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

	EventChannel chan events.Event
}

func NewGame() *State {
	return &State{
		P: Player{
			x: 1,
			y: 1,
		},
		Map: dungeon.One,

		EventChannel: make(chan events.Event, 8),
	}
}

func (s *State) RunForever() {
	go func() {
		for e := range s.EventChannel {
			switch v := e.(type) {
			case events.Movement:
				s.P.Move(v.OffX, v.OffY)
				s.P.Clamp(s.Map.Width, s.Map.Height)
			}
		}
	}()
}
