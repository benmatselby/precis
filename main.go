package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/benmatselby/precis/version"
	"github.com/gizak/termui"
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

	vstsAccount           string
	vstsProject           string
	vstsTeam              string
	vstsToken             string
	vstsBuildBranchFilter string
	vstsBuildCount        int

	githubOwner string
	githubToken string

	currentIteration string
	interval         string

	displayTravis bool
	displayVsts   bool
	displayGitHub bool

	debug bool
)

func init() {
	flag.StringVar(&travisToken, "travis-token", os.Getenv("TRAVIS_CI_TOKEN"), "The Travis CI authentication token (or define env var TRAVIS_CI_TOKEN)")
	flag.StringVar(&travisOwner, "travis-owner", os.Getenv("TRAVIS_CI_OWNER"), "The Travis CI owner (or define env var TRAVIS_CI_OWNER)")

	flag.StringVar(&vstsAccount, "vsts-account", os.Getenv("VSTS_ACCOUNT"), "The Visual Studio Team Services account (or define env var VSTS_ACCOUNT)")
	flag.StringVar(&vstsProject, "vsts-project", os.Getenv("VSTS_PROJECT"), "The Visual Studio Team Services project (or define env var VSTS_PROJECT)")
	flag.StringVar(&vstsTeam, "vsts-team", os.Getenv("VSTS_TEAM"), "The Visual Studio Team Services team (or define env var VSTS_TEAM)")
	flag.StringVar(&vstsToken, "vsts-token", os.Getenv("VSTS_TOKEN"), "The Visual Studio Team Services auth token (or define env var VSTS_TOKEN)")
	flag.IntVar(&vstsBuildCount, "vsts-build-count", 10, "How many builds should we display")
	flag.StringVar(&vstsBuildBranchFilter, "vsts-build-branch", "master", "Comma separated list of branches to display")

	flag.StringVar(&githubToken, "github-token", os.Getenv("GITHUB_TOKEN"), "The GitHub CI authentication token (or define env var GITHUB_TOKEN)")
	flag.StringVar(&githubOwner, "github-owner", os.Getenv("GITHUB_OWNER"), "The GitHub CI owner (or define env var GITHUB_OWNER)")

	flag.StringVar(&currentIteration, "current-iteration", "", "What is the current iteration")
	flag.StringVar(&interval, "interval", "60s", "The refresh rate for the dashboard")

	flag.BoolVar(&displayTravis, "display-travis", true, "Do you want to show Travis CI information?")
	flag.BoolVar(&displayVsts, "display-vsts", false, "Do you want to show Visual Studio Team Services information?")
	flag.BoolVar(&displayGitHub, "display-github", true, "Do you want to show GitHub information?")

	flag.Usage = printUsage
	flag.Parse()

	if displayTravis && (travisToken == "" || travisOwner == "") {
		printUsage()
		os.Exit(1)
	}

	if displayVsts && (vstsAccount == "" || vstsProject == "" || vstsTeam == "" || vstsToken == "") {
		printUsage()
		os.Exit(1)
	}

	if displayGitHub && (githubOwner == "" || githubToken == "") {
		printUsage()
		os.Exit(1)
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

	// Initialize termui.
	if err := termui.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "initializing termui failed: %v", err)
		os.Exit(2)
	}
	defer termui.Close()

	// Create termui widgets for google analytics.
	go titleWidget(nil)
	go vstsWidget(nil)
	go travisWidget(nil)
	go githubWidget(nil)

	// Calculate the layout.
	termui.Body.Align()
	// Render the termui body.
	termui.Render(termui.Body)

	// Handle key q pressing
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		ticker.Stop()
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		// handle Ctrl + c combination
		ticker.Stop()
		termui.StopLoop()
	})

	// Handle resize
	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	// Update on an interval
	go func() {
		for range ticker.C {
			body := termui.NewGrid()
			body.X = 0
			body.Y = 0
			body.BgColor = termui.ThemeAttr("bg")
			body.Width = termui.TermWidth()

			titleWidget(body)
			vstsWidget(body)
			travisWidget(body)
			githubWidget(body)

			// Calculate the layout.
			body.Align()
			// Render the termui body.
			termui.Render(body)
		}
	}()

	// Start the loop.
	termui.Loop()
}

func getFailureDisplay(target string) *termui.Table {
	w := termui.NewTable()
	w.Rows = [][]string{
		{"Failure"},
	}
	w.FgColor = termui.ColorWhite
	w.BgColor = termui.ColorDefault
	w.TextAlign = termui.AlignLeft
	w.Border = true
	w.BorderLabelFg = termui.ColorRed
	w.Block.BorderLabel = "Failed to get data for " + target

	w.Analysis()
	w.SetSize()

	return w
}
