// package view provides tcell UI and rendering.
package view

import (
	"fmt"
	"os"

	"github.com/cpu/yasp/dungeon"
	"github.com/cpu/yasp/game"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
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
	display *Display

	content *views.BoxLayout
	keybar  *views.SimpleStyledText
	status  *views.SimpleStyledTextBar

	topContent    *views.BoxLayout
	dungeonView   *views.CellView
	dungeonModel  *dungeonModel
	inventoryView *inventoryView

	questlogView  *views.CellView
	questlogModel *questlogModel

	views.Panel
}

type inventoryView struct {
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
	game            game.State
	highlightPlayer bool
	loc             string
}

func (m *dungeonModel) GetBounds() (int, int) {
	return m.game.Map.Dimensions()
}

func (m *dungeonModel) MoveCursor(offX, offY int) {
	curX, curY := m.game.P.Pos()
	m.game.P.MoveTo(curX+offX, curY+offY)
	m.game.P.Clamp(m.game.Map.Width, m.game.Map.Height)
	m.loc = m.game.P.String()
}

func (m *dungeonModel) GetCursor() (int, int, bool, bool) {
	playerX, playerY := m.game.P.Pos()
	return playerX, playerY, true, m.highlightPlayer
}

func (m *dungeonModel) SetCursor(x int, y int) {
	m.game.P.MoveTo(x, y)
	m.game.P.Clamp(m.game.Map.Width, m.game.Map.Height)
	m.loc = m.game.P.String()
}

func (m *dungeonModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var ch rune
	if x >= m.game.Map.Width || y >= m.game.Map.Height {
		return ch, DefaultStyle.style, nil, 1
	}
	playerX, playerY := m.game.P.Pos()
	var tile dungeon.Tile
	if x == playerX && y == playerY {
		tile = dungeon.PlayerTile
	} else {
		index := x + (y * m.game.Map.Width)
		tile = dungeon.LookupTile(m.game.Map.Tiles[index])
	}
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
			case 'H', 'h':
				win.dungeonModel.highlightPlayer = !win.dungeonModel.highlightPlayer
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
	if !m.highlightPlayer {
		w += "  [%AH%N] Highlight player"
	} else {
		w += "  [%AH%N] Un-highlight player"
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
}

func New(h InputHandler, t TickHandler) (*Display, error) {
	d := &Display{
		inputHandler: h,
		tickHandler:  t,
	}

	app := &views.Application{}
	window := &mainWindow{
		display:       d,
		inventoryView: &inventoryView{},
		dungeonModel: &dungeonModel{
			game: game.NewGame(),
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

	invTitleText := views.NewSimpleStyledTextBar()
	invTitleText.SetCenter("Inventory")
	invText := views.NewSimpleStyledText()
	invText.SetMarkup("An egg\nA sword\nLow-quality meats\n")

	window.inventoryView.SetTitle(invTitleText)
	window.inventoryView.SetContent(invText)

	window.content = views.NewBoxLayout(views.Vertical)

	window.topContent = views.NewBoxLayout(views.Horizontal)
	window.topContent.AddWidget(window.dungeonView, 1.0)
	window.topContent.AddWidget(window.inventoryView, 0.0)

	window.content.AddWidget(window.topContent, 1.0)
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
}
