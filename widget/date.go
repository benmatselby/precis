package widget

import (
	"time"

	"github.com/gizak/termui"
)

// Date renders todays date
func Date() *termui.Par {

	w := termui.NewPar(time.Now().Format(" Monday, 2 January 2006 @ 15:04:05"))
	w.Height = 3
	w.Width = 50
	w.TextFgColor = termui.ColorWhite
	w.BorderLabel = "Today"
	w.BorderLabelFg = termui.ColorCyan
	w.BorderFg = termui.ColorWhite

	return w
}
