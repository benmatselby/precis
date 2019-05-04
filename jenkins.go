package main

import (
	"time"

	"github.com/benmatselby/precis/jenkins"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func getJenkins() *widgets.Table {
	client := jenkins.New(jenkinsURL, jenkinsUsername, jenkinsPassword)

	jobs, err := client.GetJobs(jenkinsView)
	if err != nil {
		return renderError("Jenkins", err.Error())
	}

	rows := [][]string{
		{"build", "state", "finished"},
	}
	sadRows := []int{}
	happyRows := []int{}
	buildRows := []int{}
	waitRows := []int{}

	for _, job := range jobs {
		if job.LastBuild.Result == "" && job.LastBuild.Timestamp != 0 {
			job.LastBuild.Result = "RUNNING"
		}

		if job.LastBuild.Result == "" && job.LastBuild.Timestamp == 0 {
			job.LastBuild.Result = "WAITING"
		}

		finishedAt := time.Unix(0, int64(time.Millisecond)*job.LastBuild.Timestamp).Format("02-01-2006 15:04")

		rows = append(rows, []string{job.DisplayName, job.LastBuild.Result, finishedAt})

		if job.LastBuild.Result == "FAILURE" {
			sadRows = append(sadRows, len(rows)-1)
		} else if job.LastBuild.Result == "RUNNING" {
			buildRows = append(buildRows, len(rows)-1)
		} else if job.LastBuild.Result == "WAITING" {
			waitRows = append(waitRows, len(rows)-1)
		} else {
			happyRows = append(happyRows, len(rows)-1)
		}
	}

	w := widgets.NewTable()
	w.Rows = rows
	w.TextStyle = ui.NewStyle(ui.ColorWhite)
	w.TextAlignment = ui.AlignLeft
	w.Border = true
	w.Title = "Jenkins Builds - " + jenkinsURL
	w.TitleStyle = ui.Style{Fg: ui.ColorGreen}

	for _, line := range sadRows {
		w.RowStyles[line] = ui.Style{Fg: ui.ColorRed}
	}

	for _, line := range buildRows {
		w.RowStyles[line] = ui.Style{Fg: ui.ColorYellow}
	}

	for _, line := range waitRows {
		w.RowStyles[line] = ui.Style{Fg: ui.ColorCyan}
	}

	for _, line := range happyRows {
		w.RowStyles[line] = ui.Style{Fg: ui.ColorGreen}
	}

	return w
}
