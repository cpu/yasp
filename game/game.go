// package game provides the overall game elements and structure.
package game

import (
	"fmt"
	"sync"

	"github.com/cpu/yasp/dungeon"
	"github.com/cpu/yasp/game/events"
)

type State struct {
	Debug        bool
	EventChannel chan events.Event
	UIEvents     chan events.Event

	mapp   dungeon.Map
	log    QuestLog
	player Player

	sync.RWMutex
}

func NewGame() *State {
	return &State{
		EventChannel: make(chan events.Event, 32),

		mapp: dungeon.GenerateMap(1337, 16, 16),
		log:  QuestLog{},
		player: Player{
			x: 1,
			y: 1,
		},
	}
}

func (s *State) QuestLog() *QuestLog {
	return &s.log
}

func (s *State) GetPlayerPos() (int, int) {
	s.RLock()
	defer s.RUnlock()

	return s.player.Pos()
}

func (s *State) GetMapDimensions() (int, int) {
	s.RLock()
	defer s.RUnlock()

	return s.mapp.Dimensions()
}

func (s *State) GetMapTile(x, y int) (dungeon.Tile, error) {
	s.RLock()
	defer s.RUnlock()

	var tile dungeon.Tile
	maxX, maxY := s.GetMapDimensions()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return tile, fmt.Errorf("provided x,y (%d,%d) is outside of map bounds (%d, %d)",
			x, y, maxX, maxY)
	}
	playerX, playerY := s.GetPlayerPos()
	if x == playerX && y == playerY {
		tile = dungeon.PlayerTile
	} else {
		tile = s.mapp.GetTile(x, y)
	}
	return tile, nil
}

func (s *State) RunForever() {
	go func() {
		for e := range s.EventChannel {
			switch v := e.(type) {
			case events.Movement:
				s.Lock()
				oldX, oldY := s.player.Pos()
				s.player.Move(v.OffX, v.OffY)
				x, y := s.player.Clamp(s.mapp.Width, s.mapp.Height)
				if oldX != x || oldY != y {
					s.log.RecordPlayerMovement(x, y, oldX, oldY)
				}
				s.Unlock()
			case events.KeyPress:
				fmt.Printf("\nPressed %s\n", string(v.Key))
			}
		}
	}()
}
