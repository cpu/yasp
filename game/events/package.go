// package events provides game events that other components (like the view) can
// consume.
package events

type Event interface{}

type Movement struct {
	OffX int
	OffY int
}

type KeyPress struct {
	Key rune
}
