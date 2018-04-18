package widget

import (
	"log"

	travis "github.com/Ableton/go-travis"
	"github.com/gizak/termui"
)

// Travis will get the data from the Travis CI API
func Travis(token string, owner string) *termui.Table {
	client := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, token)
	opt := &travis.RepositoryListOptions{OwnerName: owner, Active: true}
	repos, _, err := client.Repositories.Find(opt)
	if err != nil {
		log.Fatal(err)
	}

	rows := [][]string{
		{"repo", "state", "finished"},
	}
	sadRows := []int{}
	happyRows := []int{}

	for _, repo := range repos {
		// Trying to remove the items that are not really running in Travis CI
		// Assume there is a better way to do this?
		if repo.LastBuildState == "" {
			continue
		}

		branch, _, err := client.Branches.GetFromSlug(repo.Slug, "master")
		if err != nil {
			log.Fatal(err)
		}

		rows = append(rows, []string{repo.Slug, branch.State, branch.FinishedAt})

		if branch.State == "failed" {
			sadRows = append(sadRows, len(rows)-1)
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
	w.Block.BorderLabel = "Travis CI builds - " + owner

	w.Analysis()
	w.SetSize()

	for _, line := range sadRows {
		w.FgColors[line] = termui.ColorRed
	}

	for _, line := range happyRows {
		w.FgColors[line] = termui.ColorDefault
	}

	return w
}
