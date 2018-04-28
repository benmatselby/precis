package main

import (
	"time"

	"github.com/gizak/termui"
	"github.com/sirupsen/logrus"
)

func doDate() (*termui.Par, error) {

	w := termui.NewPar(time.Now().Format(" Monday, 2 January 2006 @ 15:04:05"))
	w.Height = 3
	w.Width = 50
	w.TextFgColor = termui.ColorWhite
	w.BorderLabel = "Today"
	w.BorderLabelFg = termui.ColorCyan
	w.BorderFg = termui.ColorWhite

	return w, nil
}

func dateWidget(body *termui.Grid) {
	if body == nil {
		body = termui.Body
	}

	date, err := doDate()
	if err != nil {
		logrus.Fatal(err)
	}
	if date != nil {

		body.AddRows(termui.NewRow(termui.NewCol(12, 0, date)))

		// Calculate the layout.
		body.Align()
		// Render the termui body.
		termui.Render(body)
	}
}
