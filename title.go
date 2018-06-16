package main

import (
	"time"

	"github.com/gizak/termui"
)

func doDate() *termui.Par {

	w := termui.NewPar(time.Now().Local().Format("Monday, 2 January 2006 @ 15:04:05"))
	w.Height = 3
	w.PaddingLeft = 1
	w.PaddingRight = 1
	w.TextFgColor = termui.ColorWhite
	w.BorderLabel = "Today"
	w.BorderLabelFg = termui.ColorCyan
	w.BorderFg = termui.ColorWhite

	return w
}

func doIterationName() *termui.Par {
	titleColour := termui.ColorCyan
	if currentIteration == "" {
		currentIteration = "?"
		titleColour = termui.ColorRed
	}

	w := termui.NewPar(currentIteration)
	w.Height = 3
	w.PaddingLeft = 1
	w.PaddingRight = 1
	w.TextFgColor = termui.ColorWhite
	w.BorderLabel = "Sprint"
	w.BorderLabelFg = titleColour
	w.BorderFg = termui.ColorWhite

	return w
}

func titleWidget(body *termui.Grid) {
	if body == nil {
		body = termui.Body
		// It seems that if we don't pause on the first iteration
		// of this widget, we get a crash in docker.
		time.Sleep(1 * time.Second)
	}

	body.AddRows(
		termui.NewRow(
			termui.NewCol(11, 0, doDate()),
			termui.NewCol(1, 0, doIterationName()),
		),
	)

	// Calculate the layout.
	body.Align()
	// Render the termui body.
	termui.Render(body)
}
