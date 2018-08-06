package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gizak/termui"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func doGitHub() (*termui.Table, error) {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{}
	repos, _, err := client.Repositories.List(ctx, githubOwner, opt)
	if err != nil {
		log.Fatal(err)
	}

	rows := [][]string{
		{"repo", "title"},
	}

	for _, repo := range repos {
		if repo.GetFork() {
			continue
		}

		opt := &github.PullRequestListOptions{
			State: "open",
		}
		prs, _, err := client.PullRequests.List(ctx, githubOwner, repo.GetName(), opt)
		if err != nil {
			rows = append(rows, []string{repo.GetName(), err.Error()})
			log.Fatal(err)
		}

		for _, pr := range prs {
			rows = append(rows, []string{repo.GetName(), fmt.Sprintf("#%v - %s", pr.GetNumber(), pr.GetTitle())})
		}
	}

	w := termui.NewTable()
	w.Rows = rows
	w.FgColor = termui.ColorWhite
	w.BgColor = termui.ColorDefault
	w.TextAlign = termui.AlignLeft
	w.Border = true
	w.Block.BorderLabel = "GitHub Pull Requests - " + githubOwner

	w.Analysis()
	w.SetSize()

	return w, nil
}

func githubWidget(body *termui.Grid) {
	// if displayTravis == false {
	// 	return
	// }

	if body == nil {
		body = termui.Body
	}

	github, err := doGitHub()
	if err != nil {
		github = getFailureDisplay("GitHub Pull Requests")
	}
	if github != nil {
		body.AddRows(termui.NewRow(termui.NewCol(12, 0, github)))

		// Calculate the layout.
		body.Align()
		// Render the termui body.
		termui.Render(body)
	}
}
