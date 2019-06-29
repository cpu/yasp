package view

import (
	"github.com/cpu/yasp/dungeon"
	"github.com/gdamore/tcell"
)

var (
	DefaultStyle = tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlack)
	PaleGreen = tcell.StyleDefault.
			Foreground(tcell.ColorPaleGreen).
			Background(tcell.ColorBlack)
	Green = tcell.StyleDefault.
		Foreground(tcell.ColorGreen).
		Background(tcell.ColorBlack)
	Chocolate = tcell.StyleDefault.
			Foreground(tcell.ColorChocolate).
			Background(tcell.ColorBlack)
	Brown = tcell.StyleDefault.
		Foreground(tcell.ColorBrown).
		Background(tcell.ColorBlack)
)

func runeToStyle(r rune) tcell.Style {
	var style tcell.Style
	switch r {
	case dungeon.GroundTile.Rune():
		style = Green
	case dungeon.WallTile.Rune():
		style = Chocolate
	case dungeon.MossTile.Rune():
		style = PaleGreen
	case dungeon.StumpTile.Rune():
		style = Brown
	default:
		style = DefaultStyle
	}
	return style
}
