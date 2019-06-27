// package view provides tcell UI and rendering.
package view

import (
	"fmt"
	"os"

	"github.com/cpu/yasp/dungeon"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	runewidth "github.com/mattn/go-runewidth"
)

var (
	DefaultStyle = Style{
		style: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlack),
	}
	PaleGreen = Style{
		style: tcell.StyleDefault.
			Foreground(tcell.ColorPaleGreen).
			Background(tcell.ColorBlack),
	}
	Green = Style{
		style: tcell.StyleDefault.
			Foreground(tcell.ColorGreen).
			Background(tcell.ColorBlack),
	}
	Chocolate = Style{
		style: tcell.StyleDefault.
			Foreground(tcell.ColorChocolate).
			Background(tcell.ColorBlack),
	}
	Brown = Style{
		style: tcell.StyleDefault.
			Foreground(tcell.ColorBrown).
			Background(tcell.ColorBlack),
	}
	InputKeyRight = InputEvent{
		ev: tcell.KeyRight,
	}
	InputKeyLeft = InputEvent{
		ev: tcell.KeyLeft,
	}
	InputKeyUp = InputEvent{
		ev: tcell.KeyUp,
	}
	InputKeyDown = InputEvent{
		ev: tcell.KeyDown,
	}
	InputDebug = InputEvent{
		ev: tcell.KeyCtrlD,
	}
)

type InputHandler func(ev InputEvent)

type TickHandler func()

type Style struct {
	style tcell.Style
}

type InputEvent struct {
	ev tcell.Key
}

type mainWindow struct {
	content       *views.BoxLayout
	dungeonView   *views.CellView
	dungeonModel  *dungeonModel
	questlogView  *views.CellView
	questlogModel *questlogModel
	keybar        *views.SimpleStyledText
	status        *views.SimpleStyledTextBar
	display       *Display

	views.Panel
}

type questlogModel struct {
	width  int
	height int
	items  []string
}

func (m *questlogModel) GetBounds() (int, int) {
	return m.width, m.height
}

func (m *questlogModel) MoveCursor(_, _ int) {
}

func (m *questlogModel) GetCursor() (int, int, bool, bool) {
	return 0, 0, false, true
}

func (m *questlogModel) SetCursor(_, _ int) {
}

func (m *questlogModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var ch rune
	if y >= len(m.items) {
		return ch, DefaultStyle.style, nil, 1
	}
	item := m.items[y]

	if x >= len(item) {
		return ch, DefaultStyle.style, nil, 1
	}

	return rune(item[x]), Green.style, nil, 1
}

type dungeonModel struct {
	x    int
	y    int
	mapp dungeon.Map
	hide bool
	enab bool
	loc  string
}

func (m *dungeonModel) GetBounds() (int, int) {
	return m.mapp.Dimensions()
}

func (m *dungeonModel) MoveCursor(offx, offy int) {
	m.x += offx
	m.y += offy
	m.limitCursor()
}

func (m *dungeonModel) limitCursor() {
	if m.x < 0 {
		m.x = 0
	}
	if m.x > m.mapp.Width-1 {
		m.x = m.mapp.Width - 1
	}
	if m.y < 0 {
		m.y = 0
	}
	if m.y > m.mapp.Height-1 {
		m.y = m.mapp.Height - 1
	}
	m.loc = fmt.Sprintf("Player %d,%d", m.x, m.y)
}

func (m *dungeonModel) GetCursor() (int, int, bool, bool) {
	return m.x, m.y, m.enab, !m.hide
}

func (m *dungeonModel) SetCursor(x int, y int) {
	m.x = x
	m.y = y
	m.limitCursor()
}

func (m *dungeonModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var ch rune
	if x >= m.mapp.Width || y >= m.mapp.Height {
		return ch, DefaultStyle.style, nil, 1
	}
	index := x + (y * m.mapp.Width)
	tile := dungeon.LookupTile(m.mapp.Tiles[index])
	var style tcell.Style
	switch label := tile.String(); label {
	case ".":
		style = Green.style
	case "#":
		style = Chocolate.style
	case "~":
		style = PaleGreen.style
	case "=":
		style = Brown.style
	default:
		style = DefaultStyle.style
	}
	return rune(tile.String()[0]), style, nil, 1
}

type Display struct {
	inputHandler InputHandler
	tickHandler  TickHandler

	mainWin *mainWindow
	app     *views.Application
}

// PrintFixed prints a msg to the provided tcell.Screen in the given style. The
// position is fixed at x,y and takes into account the rune width of the msg
// content.
func (d Display) PrintFixed(x, y int, style Style, msg string) {
	for _, c := range msg {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		//d.s.SetContent(x, y, c, comb, style.style)
		fmt.Printf("(%d,%d): %#v %#v %#v\n", x, y, c, comb, style.style)
		x += w
	}
}

func (win *mainWindow) HandleEvent(ev tcell.Event) bool {
	app := win.display.app
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			app.Quit()
			return true
		case tcell.KeyCtrlL:
			app.Refresh()
			return true

		case tcell.KeyRune:
			switch ev.Rune() {
			case 'Q', 'q':
				app.Quit()
				return true
			case 'S', 's':
				win.dungeonModel.hide = false
				win.updateKeys()
				return true
			case 'H', 'h':
				win.dungeonModel.hide = true
				win.updateKeys()
				return true
			case 'E', 'e':
				win.dungeonModel.enab = true
				win.updateKeys()
				return true
			case 'D', 'd':
				win.dungeonModel.enab = false
				win.updateKeys()
				return true
			}
		}
	}

	return win.Panel.HandleEvent(ev)
}

func (win *mainWindow) Draw() {
	win.status.SetLeft(win.dungeonModel.loc)
	win.Panel.Draw()
}

func (win *mainWindow) updateKeys() {
	m := win.dungeonModel
	w := "[%AQ%N] Quit"
	if !m.enab {
		w += "  [%AE%N] Enable cursor"
	} else {
		w += "  [%AD%N] Disable cursor"
		if !m.hide {
			w += "  [%AH%N] Hide cursor"
		} else {
			w += "  [%AS%N] Show cursor"
		}
	}
	app := win.display.app
	win.keybar.SetMarkup(w)
	app.Update()
}

func (d *Display) RunForever() {
	if e := d.app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}

	/*
			d.s.Clear()
			d.tickHandler()
			d.s.Show()
		d.s.Fini()
	*/
}

func New(h InputHandler, t TickHandler) (*Display, error) {
	d := &Display{
		inputHandler: h,
		tickHandler:  t,
	}

	app := &views.Application{}
	window := &mainWindow{
		display: d,
		dungeonModel: &dungeonModel{
			mapp: dungeon.One,
		},
		questlogModel: &questlogModel{
			width: dungeon.One.Width,
			items: []string{
				"You were eaten by a grue.",
				"A pox on your soul was lifted.",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"You are now a chicken",
				"and I am too",
			},
		},
	}

	title := views.NewTextBar()
	title.SetStyle(DefaultStyle.style)
	title.SetCenter("Y A S P", tcell.StyleDefault)
	title.SetRight("HP: 0", tcell.StyleDefault)

	window.keybar = views.NewSimpleStyledText()
	window.keybar.RegisterStyle('N', tcell.StyleDefault.
		Background(tcell.ColorSilver).
		Foreground(tcell.ColorBlack))
	window.keybar.RegisterStyle('A', tcell.StyleDefault.
		Background(tcell.ColorSilver).
		Foreground(tcell.ColorRed))

	window.status = views.NewSimpleStyledTextBar()
	window.status.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorYellow))
	window.status.RegisterLeftStyle('N', tcell.StyleDefault.
		Background(tcell.ColorYellow).
		Foreground(tcell.ColorBlack))

	window.status.SetLeft("Player x, y")
	window.status.SetRight("%U2019%N")
	window.status.SetCenter("%U@cpu%N")

	window.dungeonView = views.NewCellView()
	window.dungeonView.SetModel(window.dungeonModel)
	window.dungeonView.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorBlack))

	window.questlogView = views.NewCellView()
	window.questlogView.SetModel(window.questlogModel)
	window.questlogView.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack))

	window.content = views.NewBoxLayout(views.Vertical)
	window.content.AddWidget(window.dungeonView, 1.0)
	window.content.AddWidget(window.questlogView, 1.0)

	window.SetMenu(window.keybar)
	window.SetTitle(title)
	window.SetContent(window.content)
	window.SetStatus(window.status)

	app.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

	d.app = app
	d.mainWin = window

	window.updateKeys()

	app.SetRootWidget(window)
	return d, nil

	/*
		tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
		s, err := tcell.NewScreen()
		if err != nil {
			return nil, err
		}

		if err := s.Init(); err != nil {
			return nil, err
		}

		s.SetStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlack))
		s.Clear()
	*/

}
