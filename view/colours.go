package view

import "github.com/gdamore/tcell"

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
	case '.':
		style = Green
	case '#':
		style = Chocolate
	case '~':
		style = PaleGreen
	case '=':
		style = Brown
	default:
		style = DefaultStyle
	}
	return style
}
