package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/gizak/termui"
)

func getAzureDevOpsBuildsForBranch(defID int, branchName string) ([]azuredevops.Build, error) {
	client := azuredevops.NewClient(azureDevOpsAccount, azureDevOpsProject, azureDevOpsToken)
	buildOpts := azuredevops.BuildsListOptions{Definitions: strconv.Itoa(defID), Branch: "refs/heads/" + branchName, Count: 1}
	build, err := client.Builds.List(&buildOpts)
	return build, err
}

func getAzureDevOpsBuilds() (*termui.Table, error) {
	client := azuredevops.NewClient(azureDevOpsAccount, azureDevOpsProject, azureDevOpsToken)

	buildDefOpts := azuredevops.BuildDefinitionsListOptions{Path: "\\" + azureDevOpsTeam}
	definitions, err := client.BuildDefinitions.List(&buildDefOpts)
	if err != nil {
		return nil, err
	}

	var builds []azuredevops.Build
	for _, definition := range definitions {
		for _, branchName := range strings.Split(azureDevOpsBuildBranchFilter, ",") {
			build, err := getAzureDevOpsBuildsForBranch(definition.ID, branchName)
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

	if len(builds) < azureDevOpsBuildCount {
		azureDevOpsBuildCount = len(builds)
	}

	for index := 0; index < azureDevOpsBuildCount; index++ {
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
	w.BorderLabelFg = termui.ColorGreen
	w.Block.BorderLabel = "Azure DevOps CI Builds - " + azureDevOpsTeam

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

func getAzureDevOpsPulls() (*termui.Table, error) {
	client := azuredevops.NewClient(azureDevOpsAccount, azureDevOpsProject, azureDevOpsToken)
	opts := azuredevops.PullRequestListOptions{State: "active"}
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
	w.BorderLabelFg = termui.ColorGreen
	w.Block.BorderLabel = "Azure DevOps Pull Requests - " + azureDevOpsTeam

	w.Analysis()
	w.SetSize()

	return w, nil
}
