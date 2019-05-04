package main

import (
	"sort"
	"time"

	travis "github.com/Ableton/go-travis"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/viper"
)

func getTravis() *widgets.Table {
	client := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, travisToken)
	opt := &travis.RepositoryListOptions{Member: travisOwner, Active: true}
	repos, _, err := client.Repositories.Find(opt)
	if err != nil {
		return renderError("Travis", err.Error())
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
			return renderError("Travis", err.Error())
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

	w := widgets.NewTable()
	w.Rows = rows
	w.TextStyle = ui.NewStyle(ui.ColorWhite)
	w.TextAlignment = ui.AlignLeft
	w.Border = true
	w.Title = "TravisCI Builds - " + travisOwner
	w.TitleStyle = ui.Style{Fg: ui.ColorGreen}

	for _, line := range sadRows {
		w.RowStyles[line] = ui.Style{Fg: ui.ColorRed}
	}

	for _, line := range buildRows {
		w.RowStyles[line] = ui.Style{Fg: ui.ColorYellow}
	}

	for _, line := range happyRows {
		w.RowStyles[line] = ui.Style{Fg: ui.ColorGreen}
	}

	return w
}
