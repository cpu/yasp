package view

import (
	"github.com/cpu/yasp/game"
	"github.com/cpu/yasp/game/events"
	"github.com/gdamore/tcell"
)

type dungeonModel struct {
	game            *game.State
	highlightPlayer bool
}

func (m *dungeonModel) GetBounds() (int, int) {
	return m.game.GetMapDimensions()
}

func (m *dungeonModel) MoveCursor(offX, offY int) {
	m.game.EventChannel <- events.Movement{
		OffX: offX,
		OffY: offY,
	}
}

func (m *dungeonModel) GetCursor() (int, int, bool, bool) {
	playerX, playerY := m.game.GetPlayerPos()
	return playerX, playerY, true, m.highlightPlayer
}

func (m *dungeonModel) SetCursor(x int, y int) {
	curX, curY := m.game.GetPlayerPos()
	diffX := curX - x
	diffY := curY - y
	m.MoveCursor(diffX, diffY)
}

func (m *dungeonModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	tile, err := m.game.GetMapTile(x, y)
	if err != nil {
		return ' ', DefaultStyle, nil, 1
	}
	return tile.Rune(), runeToStyle(tile.Rune()), nil, 1
}
