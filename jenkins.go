package main

import (
	"time"

	"github.com/benmatselby/frost/jenkins"
	"github.com/gizak/termui"
)

func doJenkins() (*termui.Table, error) {
	client := jenkins.New(jenkinsURL, jenkinsUsername, jenkinsPassword)

	jobs, err := client.GetJobs(jenkinsView)
	if err != nil {
		return nil, err
	}

	rows := [][]string{
		{"build", "state", "finished"},
	}
	sadRows := []int{}
	happyRows := []int{}
	buildRows := []int{}

	for _, job := range jobs {
		if job.LastBuild.Result == "" {
			// Assumption made here is that this is a folder/pipline entry
			// with no useful information to render
			continue
		}

		finishedAt := time.Unix(0, int64(time.Millisecond)*job.LastBuild.Timestamp).Format("02-01-2006 15:04")

		rows = append(rows, []string{job.DisplayName, job.LastBuild.Result, finishedAt})

		if job.LastBuild.Result == "FAILURE" {
			sadRows = append(sadRows, len(rows)-1)
		} else if job.LastBuild.Result == "RUNNING" {
			buildRows = append(buildRows, len(rows)-1)
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
	w.BorderLabelFg = termui.ColorGreen
	w.Block.BorderLabel = "Jenkins Builds - " + jenkinsURL

	w.Analysis()
	w.SetSize()

	for _, line := range sadRows {
		w.FgColors[line] = termui.ColorRed
	}

	for _, line := range buildRows {
		w.FgColors[line] = termui.ColorYellow
	}

	for _, line := range happyRows {
		w.FgColors[line] = termui.ColorDefault
	}

	return w, nil
}

func jenkinsWidget(body *termui.Grid) {
	if displayJenkins == false {
		return
	}

	if body == nil {
		body = termui.Body
	}

	jenkins, err := doJenkins()
	if err != nil {
		jenkins = getFailureDisplay("Jenkins Builds")
	}
	if jenkins != nil {
		body.AddRows(termui.NewRow(termui.NewCol(12, 0, jenkins)))

		body.Align()
		termui.Render(body)
	}
}
