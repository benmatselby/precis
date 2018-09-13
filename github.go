package main

import (
	"context"
	"fmt"
	"sync"

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

	orgs, _, err := client.Organizations.List(context.Background(), githubOwner, nil)
	if err != nil {
		return nil, err
	}

	var allRepos [][]string
	for _, org := range orgs {
		opt := &github.RepositoryListByOrgOptions{}
		repos, _, err := client.Repositories.ListByOrg(ctx, org.GetLogin(), opt)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			allRepos = append(allRepos, []string{org.GetLogin(), repo.GetName()})
		}
	}

	opt := &github.RepositoryListOptions{}
	repos, _, err := client.Repositories.List(ctx, githubOwner, opt)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if repo.GetFork() {
			continue
		}
		allRepos = append(allRepos, []string{githubOwner, repo.GetName()})
	}

	rows := [][]string{
		{"repo", "title"},
	}

	pullRequests := make(chan *github.PullRequest)
	var wg sync.WaitGroup
	wg.Add(len(allRepos))

	go func() {
		wg.Wait()
		close(pullRequests)
	}()

	for _, repo := range allRepos {
		go func(repo []string) {
			defer wg.Done()
			opt := &github.PullRequestListOptions{
				State: "open",
			}

			prs, _, err := client.PullRequests.List(ctx, repo[0], repo[1], opt)
			if err != nil {
				fmt.Printf("unable to get pull requests for %s: %v", repo[1], err)
			}

			for _, pull := range prs {
				pullRequests <- pull
			}
		}(repo)
	}

	for result := range pullRequests {
		rows = append(rows, []string{result.GetHead().GetRepo().GetFullName(), fmt.Sprintf("#%v - %s", result.GetNumber(), result.GetTitle())})
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
	if displayGitHub == false {
		return
	}

	if body == nil {
		body = termui.Body
	}

	github, err := doGitHub()
	if err != nil {
		github = getFailureDisplay(err.Error())
	}
	if github != nil {
		body.AddRows(termui.NewRow(termui.NewCol(12, 0, github)))

		// Calculate the layout.
		body.Align()
		// Render the termui body.
		termui.Render(body)
	}
}
