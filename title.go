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
	w.BorderLabelFg = termui.ColorGreen
	w.BorderFg = termui.ColorWhite

	return w
}

func doIterationName() *termui.Par {
	titleColour := termui.ColorGreen
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
