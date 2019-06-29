// package view provides tcell UI and rendering.
package view

import (
	"fmt"
	"os"

	"github.com/cpu/yasp/game"
	"github.com/cpu/yasp/game/events"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

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
		return ch, DefaultStyle, nil, 1
	}
	item := m.items[y]

	if x >= len(item) {
		return ch, DefaultStyle, nil, 1
	}

	return rune(item[x]), Green, nil, 1
}

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

type Display struct {
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
			switch r := ev.Rune(); r {
			case 'Q', 'q':
				app.Quit()
				return true
			case 'H', 'h':
				win.dungeonModel.highlightPlayer = !win.dungeonModel.highlightPlayer
				win.updateKeys()
				return true
			default:
				win.dungeonModel.game.EventChannel <- events.KeyPress{
					Key: r,
				}
			}
		}
	}

	return win.Panel.HandleEvent(ev)
}

func (win *mainWindow) Draw() {
	win.status.SetLeft("TODO")
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

func New(g *game.State) (*Display, error) {
	d := &Display{}

	mapWidth, _ := g.GetMapDimensions()

	app := &views.Application{}
	window := &mainWindow{
		display:       d,
		inventoryView: &inventoryView{},
		dungeonModel: &dungeonModel{
			game: g,
		},
		questlogModel: &questlogModel{
			width: mapWidth,
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
	title.SetStyle(DefaultStyle)
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
