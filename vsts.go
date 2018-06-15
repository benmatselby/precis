package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/gizak/termui"
)

func getVstsBuildsForBranch(defID int, branchName string) ([]vsts.Build, error) {
	client := vsts.NewClient(vstsAccount, vstsProject, vstsToken)
	buildOpts := vsts.BuildsListOptions{Definitions: strconv.Itoa(defID), Branch: "refs/heads/" + branchName, Count: 1}
	build, err := client.Builds.List(&buildOpts)
	return build, err
}

func getVstsBuilds() (*termui.Table, error) {
	client := vsts.NewClient(vstsAccount, vstsProject, vstsToken)

	buildDefOpts := vsts.BuildDefinitionsListOptions{Path: "\\" + vstsTeam}
	definitions, err := client.BuildDefinitions.List(&buildDefOpts)
	if err != nil {
		return nil, err
	}

	var builds []vsts.Build
	for _, definition := range definitions {
		for _, branchName := range strings.Split(vstsBuildBranchFilter, ",") {
			build, err := getVstsBuildsForBranch(definition.ID, branchName)
			if err != nil {
				continue
			}
			if len(build) > 0 {
				builds = append(builds, build[0])
			}
		}
	}

	rows := [][]string{
		{"repo", "state", "branch", "finished"},
	}
	sadRows := []int{}
	buildingRows := []int{}

	if len(builds) < vstsBuildCount {
		vstsBuildCount = len(builds)
	}

	for index := 0; index < vstsBuildCount; index++ {
		build := builds[index]

		finish, error := time.Parse(time.RFC3339, build.FinishTime)
		finishAt := finish.Format("2006-01-02 15:04:05")
		if error != nil {
			finishAt = build.FinishTime
		}

		rows = append(rows, []string{build.Definition.Name, build.Status, build.Branch, finishAt})

		if build.Result == "failed" {
			sadRows = append(sadRows, len(rows)-1)
		}

		if build.Status == "inProgress" {
			buildingRows = append(buildingRows, len(rows)-1)
		}
	}

	w := termui.NewTable()
	w.Rows = rows
	w.FgColor = termui.ColorWhite
	w.BgColor = termui.ColorDefault
	w.TextAlign = termui.AlignLeft
	w.Border = true
	w.BorderLabelFg = termui.ColorCyan
	w.Block.BorderLabel = "VSTS CI Builds - " + vstsTeam

	w.Analysis()
	w.SetSize()

	for _, line := range sadRows {
		w.FgColors[line] = termui.ColorRed
	}

	for _, line := range buildingRows {
		w.FgColors[line] = termui.ColorYellow
	}

	return w, nil
}

func getVstsPulls() (*termui.Table, error) {
	client := vsts.NewClient(vstsAccount, vstsProject, vstsToken)
	opts := vsts.PullRequestListOptions{State: "active"}
	pulls, count, err := client.PullRequests.List(&opts)
	if err != nil {
		return nil, err
	}

	rows := [][]string{
		{"id", "repo", "title", "created"},
	}

	for index := 0; index < count; index++ {
		pull := pulls[index]

		created, error := time.Parse(time.RFC3339, pull.Created)
		createdOn := created.Format("2006-01-02 15:04:05")
		if error != nil {
			createdOn = pull.Created
		}

		rows = append(rows, []string{fmt.Sprintf("%d", pull.ID), pull.Repo.Name, pull.Title, createdOn})
	}

	if count == 0 {
		rows = append(rows, []string{"No", "open", "pull", "requests"})
	}

	w := termui.NewTable()
	w.Rows = rows
	w.FgColor = termui.ColorWhite
	w.BgColor = termui.ColorDefault
	w.TextAlign = termui.AlignLeft
	w.Border = true
	w.BorderLabelFg = termui.ColorCyan
	w.Block.BorderLabel = "VSTS Pull Requests - " + vstsTeam

	w.Analysis()
	w.SetSize()

	return w, nil
}

func vstsWidget(body *termui.Grid) {
	if displayVsts == false {
		return
	}

	if body == nil {
		body = termui.Body
	}

	builds, err := getVstsBuilds()
	if err != nil {
		builds = getFailureDisplay("VSTS CI Builds")
	}

	pulls, err := getVstsPulls()
	if err != nil {
		pulls = getFailureDisplay("VSTS Pull Requests")
	}

	if len(pulls.Rows) > 0 {
		body.AddRows(
			termui.NewRow(
				termui.NewCol(12, 0, pulls),
			),
		)
	}

	if builds != nil {
		body.AddRows(
			termui.NewRow(
				termui.NewCol(12, 0, builds),
			),
		)
	}

	// Calculate the layout.
	body.Align()
	// Render the termui body.
	termui.Render(body)
}
