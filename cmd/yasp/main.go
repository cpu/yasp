package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cpu/yasp"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	title = "Y A S P"
)

var (
	configFile   = flag.String("config", "test/config.yml", "YAML config file path")
	windowConfig = pixelgl.WindowConfig{}
)

func run() {
	win, err := pixelgl.NewWindow(windowConfig)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Greenyellow)
	for !win.Closed() {
		win.Update()
	}
}

func main() {
	flag.Parse()

	errExit := func(msg string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, msg, args...)
		os.Exit(666)
	}
	ifErrExit := func(err error, msg string, args ...interface{}) {
		if err != nil {
			errExit(msg, args...)
		}
	}

	fmt.Println(title)
	c, err := yasp.LoadConfigFile(*configFile)
	ifErrExit(err, "failed to load YASP config from %q: %v\n", *configFile, err)

	windowConfig = pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, float64(c.WinWidth), float64(c.WinHeight)),
		VSync:  c.VSync,
	}

	pixelgl.Run(run)
	fmt.Println("... goodbye for now")
}
