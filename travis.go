package main

import (
	"sort"
	"time"

	travis "github.com/Ableton/go-travis"
	"github.com/gizak/termui"
	"github.com/spf13/viper"
)

func doTravis() (*termui.Table, error) {
	client := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, travisToken)
	opt := &travis.RepositoryListOptions{OwnerName: travisOwner, Active: true}
	repos, _, err := client.Repositories.Find(opt)
	if err != nil {
		return nil, err
	}

	sort.Slice(repos, func(i, j int) bool { return repos[i].LastBuildFinishedAt > repos[j].LastBuildFinishedAt })

	rows := [][]string{
		{"repo", "state", "finished"},
	}
	sadRows := []int{}
	happyRows := []int{}
	buildRows := []int{}

	ignoreRepos := viper.GetStringSlice("travis.ignore_repos")

	for _, repo := range repos {
		// Trying to remove the items that are not really running in Travis CI
		// Assume there is a better way to do this?
		if repo.LastBuildState == "" {
			continue
		}

		// We may want to ignore certain repos
		// Use cases:
		//  - Personal and work dashboards
		//  - Deprecated repos that may have failed a build and now abandoned
		ignore := false
		for _, i := range ignoreRepos {
			if i == repo.Slug {
				ignore = true
			}
		}
		if ignore {
			continue
		}

		branch, _, err := client.Branches.GetFromSlug(repo.Slug, "master")
		if err != nil {
			return nil, err
		}

		finish, error := time.Parse(time.RFC3339, branch.FinishedAt)
		finishAt := finish.Format("2006-01-02 15:04:05")
		if error != nil {
			finishAt = branch.FinishedAt
		}

		rows = append(rows, []string{repo.Slug, branch.State, finishAt})

		if branch.State == "failed" {
			sadRows = append(sadRows, len(rows)-1)
		} else if branch.State == "started" {
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
	w.Block.BorderLabel = "Travis CI Builds - " + travisOwner

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

func travisWidget(body *termui.Grid) {
	if displayTravis == false {
		return
	}

	if body == nil {
		body = termui.Body
	}

	travis, err := doTravis()
	if err != nil {
		travis = getFailureDisplay("Travis CI Builds")
	}
	if travis != nil {
		body.AddRows(termui.NewRow(termui.NewCol(12, 0, travis)))

		// Calculate the layout.
		body.Align()
		// Render the termui body.
		termui.Render(body)
	}
}
