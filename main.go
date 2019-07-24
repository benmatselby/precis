package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/benmatselby/precis/version"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/viper"
)

const (
	// BANNER is rendered for the help
	BANNER = `
.______   .______       _______   ______  __       _______.
|   _  \  |   _  \     |   ____| /      ||  |     /       |
|  |_)  | |  |_)  |    |  |__   |  ,----'|  |    |   (----
|   ___/  |      /     |   __|  |  |     |  |     \   \
|  |      |  |\  \----.|  |____ |   ----.|  | .----)   |
| _|      | _|  ._____||_______| \______||__| |_______/

A terminal dashboard which gives an overview of useful things

Build: %s

`
)

var (
	travisToken string
	travisOwner string

	githubOwner string
	githubToken string

	jenkinsURL      string
	jenkinsUsername string
	jenkinsPassword string
	jenkinsView     string

	interval string

	displayBuild  bool
	displayGitHub bool

	grid *ui.Grid
)

func init() {
	flag.StringVar(&travisToken, "travis-token", os.Getenv("TRAVIS_CI_TOKEN"), "The Travis CI authentication token (or define env var TRAVIS_CI_TOKEN)")
	flag.StringVar(&travisOwner, "travis-owner", os.Getenv("TRAVIS_CI_OWNER"), "The Travis CI owner (or define env var TRAVIS_CI_OWNER)")

	flag.StringVar(&githubToken, "github-token", os.Getenv("GITHUB_TOKEN"), "The GitHub CI authentication token (or define env var GITHUB_TOKEN)")
	flag.StringVar(&githubOwner, "github-owner", os.Getenv("GITHUB_OWNER"), "The GitHub CI owner (or define env var GITHUB_OWNER)")

	flag.StringVar(&jenkinsURL, "jenkins-url", os.Getenv("JENKINS_URL"), "The Jenkins URL (or define env var JENKINS_URL)")
	flag.StringVar(&jenkinsUsername, "jenkins-username", os.Getenv("JENKINS_USERNAME"), "The Jenkins username to authenticate with (or define env var JENKINS_USERNAME)")
	flag.StringVar(&jenkinsPassword, "jenkins-password", os.Getenv("JENKINS_PASSWORD"), "The Jenkins password to authenticate with (or define env var JENKINS_PASSWORD)")
	flag.StringVar(&jenkinsView, "jenkins-view", os.Getenv("JENKINS_VIEW"), "The Jenkins view you want render, otherwise it is all (or define env var JENKINS_VIEW)")

	flag.StringVar(&interval, "interval", "60s", "The refresh rate for the dashboard")

	flag.BoolVar(&displayBuild, "display-build", true, "Do you want to show build information from TravisCI and Jenkins?")
	flag.BoolVar(&displayGitHub, "display-github", true, "Do you want to show GitHub information?")

	flag.Usage = printUsage
	flag.Parse()

	if displayBuild && (travisToken == "" || travisOwner == "") {
		printUsage()
		os.Exit(1)
	}

	if displayGitHub && (githubOwner == "" || githubToken == "") {
		printUsage()
		os.Exit(1)
	}

	if displayBuild && (jenkinsURL == "" || jenkinsUsername == "" || jenkinsPassword == "") {
		printUsage()
		os.Exit(1)
	}

	loadConfig()
}

func loadConfig() {
	viper.SetConfigName("precis")
	viper.AddConfigPath("$HOME/.benmatselby")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Fprintf(os.Stderr, "Failed to load config file: %s", err)
	}
}

func printUsage() {
	fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.GITCOMMIT))
	flag.PrintDefaults()
}

func main() {
	var ticker *time.Ticker

	// parse the duration
	dur, err := time.ParseDuration(interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing %s as duration failed: %v", interval, err)
		os.Exit(2)
	}
	ticker = time.NewTicker(dur)

	// Initialize ui.
	if err := ui.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "initializing termui failed: %v", err)
		os.Exit(2)
	}
	defer ui.Close()

	displayLoading()
	displayWidgets()

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker.C:
			displayWidgets()
		}
	}
}

func getDateTime() *widgets.Paragraph {
	w := widgets.NewParagraph()
	w.Text = time.Now().Local().Format("Monday, 2 January 2006 @ 15:04:05")
	w.PaddingLeft = 1
	w.PaddingRight = 1
	w.Title = "Today"
	return w
}

func displayWidgets() {
	grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	if displayBuild && displayGitHub {
		github := getGitHub()
		jenkins := getJenkins()
		travis := getTravis()

		grid.Set(
			ui.NewRow(1.0/9,
				getDateTime(),
			),
			ui.NewRow(4.5/10,
				ui.NewCol(1.0, github),
			),
			ui.NewRow(4.5/10,
				ui.NewCol(1.0/2, travis),
				ui.NewCol(1.0/2, jenkins),
			),
		)
	} else if displayBuild && !displayGitHub {
		jenkins := getJenkins()
		travis := getTravis()
		grid.Set(
			ui.NewRow(1.0/9,
				getDateTime(),
			),
			ui.NewRow(9.0/10,
				ui.NewCol(1.0/2, travis),
				ui.NewCol(1.0/2, jenkins),
			),
		)
	} else if displayGitHub && !displayBuild {
		github := getGitHub()
		grid.Set(
			ui.NewRow(1.0/9,
				getDateTime(),
			),
			ui.NewRow(9.0/10,
				ui.NewCol(1.0, github),
			),
		)
	}

	ui.Render(grid)
}

func displayLoading() {
	w := widgets.NewParagraph()
	w.Text = "Loading all the data"
	w.PaddingLeft = 1
	w.PaddingRight = 1

	grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0, w),
	)
	ui.Render(grid)
}

func renderError(service, msg string) *widgets.Table {
	w := widgets.NewTable()
	w.Rows = [][]string{{"Failed: " + msg}}
	w.Border = true
	w.Title = service
	w.TitleStyle = ui.Style{Fg: ui.ColorRed}
	return w
}
