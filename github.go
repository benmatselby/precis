package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func getGitHub() *widgets.Table {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	orgs, _, err := client.Organizations.List(context.Background(), githubOwner, nil)
	if err != nil {
		return renderError("GitHub", err.Error())
	}

	var allRepos [][]string
	for _, org := range orgs {
		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		repos, _, err := client.Repositories.ListByOrg(ctx, org.GetLogin(), opt)
		if err != nil {
			return renderError("GitHub", err.Error())
		}

		for _, repo := range repos {
			if !showRepoPr(org.GetLogin(), repo.GetName()) {
				continue
			}
			allRepos = append(allRepos, []string{org.GetLogin(), repo.GetName()})
		}
	}

	opt := &github.RepositoryListOptions{}
	repos, _, err := client.Repositories.List(ctx, githubOwner, opt)
	if err != nil {
		return renderError("GitHub", err.Error())
	}

	for _, repo := range repos {
		if repo.GetFork() {
			continue
		}
		if !showRepoPr(githubOwner, repo.GetName()) {
			continue
		}
		allRepos = append(allRepos, []string{githubOwner, repo.GetName()})
	}

	rows := [][]string{}

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

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	rows = append([][]string{{"repo", "title"}}, rows...)

	w := widgets.NewTable()
	w.Rows = rows
	w.TextStyle = ui.NewStyle(ui.ColorWhite)
	w.TextAlignment = ui.AlignLeft
	w.Border = true
	w.Title = "GitHub Pull Requests - " + githubOwner

	return w
}

// showRepoPr is going to determine if we care enough to show the detail
func showRepoPr(org, name string) bool {
	watchRepos := viper.GetStringSlice("github.pull_request_repos")

	if len(watchRepos) == 0 {
		return true
	}

	show := false
	for _, i := range watchRepos {
		s := strings.Split(i, "/")

		// If we want to watch everything for a given org
		if s[1] == "*" && org == s[0] {
			show = true
			break
		}

		// If we want to watch everything for a given repo (including forks)
		if s[0] == "*" && name == s[1] {
			show = true
			break
		}

		// Otherwise we want an exact match
		if i == fmt.Sprintf("%s/%s", org, name) {
			show = true
			break
		}
	}

	return show
}
