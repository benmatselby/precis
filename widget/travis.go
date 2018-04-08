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
	for _, repo := range repos {

		if repo.LastBuildState == "" {
			continue
		}
		branch, _, err := client.Branches.GetFromSlug(repo.Slug, "master")
		if err != nil {
			log.Fatal(err)
		}
		rows = append(rows, []string{repo.Slug, branch.State, branch.FinishedAt})
	}

	w := termui.NewTable()

	w.Rows = rows
	w.FgColor = termui.ColorWhite
	w.BgColor = termui.ColorDefault
	w.TextAlign = termui.AlignLeft
	w.Border = true
	w.Block.BorderLabel = "Travis CI builds"

	w.Analysis()
	w.SetSize()

	return w
}
