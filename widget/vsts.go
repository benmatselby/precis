package widget

import (
	"log"
	"time"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/gizak/termui"
)

// VstsBuilds will get build information from Visual Studio Team Services
func VstsBuilds(token string, account string, project string, team string, count int) *termui.Table {
	client := vsts.NewClient(account, project, token)
	builds, error := client.Builds.List()
	if error != nil {
		log.Fatal(error)
	}

	rows := [][]string{
		{"repo", "state", "branch", "finished"},
	}
	sadRows := []int{}
	happyRows := []int{}

	for index := 0; index < count; index++ {
		build := builds[index]

		finish, error := time.Parse(time.RFC3339, build.FinishTime)
		finishAt := finish.Format("2006-01-02 15:04:05")
		if error != nil {
			finishAt = build.FinishTime
		}

		rows = append(rows, []string{build.Definition.Name, build.Status, build.Branch, finishAt})

		if build.Result == "failed" {
			sadRows = append(sadRows, len(rows)-1)
		} else {
			happyRows = append(happyRows, len(rows)-1)
		}
	}

	w := termui.NewTable()
	w.Rows = rows
	w.FgColor = termui.ColorWhite
	w.BgColor = termui.ColorDefault
	w.TextAlign = termui.AlignLeft
	w.Border = true
	w.Block.BorderLabel = "VSTS CI builds - " + team

	w.Analysis()
	w.SetSize()

	for _, line := range sadRows {
		w.FgColors[line] = termui.ColorRed
	}

	for _, line := range happyRows {
		w.FgColors[line] = termui.ColorDefault
	}

	return w
}
