package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/benmatselby/precis/widget"
	"github.com/gizak/termui"
)

const (
	// BANNER is rendered for the help
	BANNER = `
A terminal dashboard which gives an overview of useful things
`
)

var (
	travisToken string
	travisOwner string

	interval string
)

func init() {
	flag.StringVar(&travisToken, "travis-token", os.Getenv("TRAVIS_CI_TOKEN"), "The Travis CI authentication token")
	flag.StringVar(&travisOwner, "travis-owner", os.Getenv("TRAVIS_CI_OWNER"), "The Travis CI owner")

	flag.StringVar(&interval, "interval", "60s", "The refresh rate for the dashboard")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER))
		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	duration, err := time.ParseDuration(interval)
	if err != nil {
		duration, _ = time.ParseDuration("60s")
	}
	ticker := time.NewTicker(duration)

	// Keyboard message
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		ticker.Stop()
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		ticker.Stop()
		termui.StopLoop()
	})

	// Handle resize
	termui.Handle("/sys/wnd/resize", func(e termui.Event) {

	})

	go func() {
		for range ticker.C {
			exec()
		}
	}()

	exec()
	termui.Loop()
}

func exec() {
	body := termui.NewGrid()
	body.X = 0
	body.Y = 0
	body.BgColor = termui.ThemeAttr("bg")
	body.Width = termui.TermWidth()

	body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, widget.Date()),
		),
		termui.NewRow(
			termui.NewCol(12, 0, widget.Travis(travisToken, travisOwner)),
		),
	)

	body.Align()
	termui.Render(body)
}
