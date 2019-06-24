// package view provides tcell UI and rendering.
package view

import (
	"time"

	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
)

type InputHandler func(ev *tcell.EventKey)

type TickHandler func()

type Display struct {
	inputHandler InputHandler
	tickHandler  TickHandler
	quit         chan struct{}
	s            tcell.Screen
}

// PrintFixed prints a msg to the provided tcell.Screen in the given style. The
// position is fixed at x,y and takes into account the rune width of the msg
// content.
func (d Display) PrintFixed(x, y int, style tcell.Style, msg string) {
	for _, c := range msg {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		d.s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func (d *Display) pollForever() {
	go func() {
		for {
			ev := d.s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(d.quit)
					return
				case tcell.KeyCtrlL:
					d.s.Sync()
					/*
					* Game controls
					 */
				case tcell.KeyRight:
					fallthrough
				case tcell.KeyLeft:
					fallthrough
				case tcell.KeyUp:
					fallthrough
				case tcell.KeyDown:
					d.inputHandler(ev)
				}
			case *tcell.EventResize:
				d.s.Sync()
			}
		}
	}()
}

func (d *Display) RunForever() {
	d.pollForever()

loop:
	for {
		select {
		case <-d.quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
		d.tickHandler()
	}

	d.Close()
}

func (d Display) Size() (int, int) {
	return d.s.Size()
}

func (d Display) Clear() {
	d.s.Clear()
}

func (d Display) Close() {
	d.s.Fini()
}

func (d Display) Show() {
	d.s.Show()
}

func New(h InputHandler, t TickHandler) (*Display, error) {

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

	return &Display{
		quit:         make(chan struct{}),
		inputHandler: h,
		tickHandler:  t,
		s:            s,
	}, nil
}
