package widget

import (
	"context"
	"log"

	"github.com/gizak/termui"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Github will get the data from GitHub such as open pull requests etc
func Github(token string, owner string) *termui.Table {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{}
	repos, _, err := client.Repositories.List(ctx, owner, opt)
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
		prs, _, err := client.PullRequests.List(ctx, owner, repo.GetName(), opt)
		if err != nil {
			rows = append(rows, []string{repo.GetName(), err.Error()})
			log.Fatal(err)
		}

		for _, pr := range prs {
			rows = append(rows, []string{repo.GetName(), pr.GetTitle()})
		}
	}

	w := termui.NewTable()
	w.Rows = rows
	w.FgColor = termui.ColorWhite
	w.BgColor = termui.ColorDefault
	w.TextAlign = termui.AlignLeft
	w.Border = true
	w.Block.BorderLabel = "GitHub Pull Requests - " + owner

	w.Analysis()
	w.SetSize()

	return w
}
