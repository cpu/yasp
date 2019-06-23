package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cpu/yasp"
	"github.com/cpu/yasp/game"
	"github.com/faiface/pixel/pixelgl"
)

const (
	title = "Y A S P"
)

var (
	configFile = flag.String("config", "test/config.yml", "YAML config file path")
)

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
	ifErrExit(err, "failed to load config from %q: %v\n", *configFile, err)

	g, err := game.New(c)
	ifErrExit(err, "failed to construct game: %q\n", err)

	pixelgl.Run(g.Run)
	fmt.Println("... goodbye for now")
}
