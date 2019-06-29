package view

import (
	"github.com/cpu/yasp/game"
	"github.com/gdamore/tcell"
)

type questlogModel struct {
	width int
	log   *game.QuestLog
}

func (m *questlogModel) GetBounds() (int, int) {
	return m.width, m.log.Len()
}

func (m *questlogModel) MoveCursor(_, _ int) {
	// NOP for now
}

func (m *questlogModel) GetCursor() (int, int, bool, bool) {
	// NOP for now
	return 0, 0, true, false
}

func (m *questlogModel) SetCursor(_, _ int) {
	// NOP for now
}

func (m *questlogModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	style := tcell.StyleDefault.Background(tcell.ColorGray)
	green := style.Foreground(tcell.ColorGreen)

	y = m.log.Len() - y
	var ch rune
	if y >= m.log.Len() {
		return ch, style, nil, 1
	}
	item, err := m.log.GetItem(y)
	if err != nil {
		return ch, style, nil, 1
	}

	itemStr := item.String()
	if x >= len(itemStr) {
		return ch, style, nil, 1
	}

	return rune(itemStr[x]), green, nil, 1
}
