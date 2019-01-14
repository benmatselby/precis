package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/benmatselby/precis/version"
	"github.com/gizak/termui"
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

	azureDevOpsAccount           string
	azureDevOpsProject           string
	azureDevOpsTeam              string
	azureDevOpsToken             string
	azureDevOpsBuildBranchFilter string
	azureDevOpsBuildCount        int

	githubOwner string
	githubToken string

	jenkinsURL      string
	jenkinsUsername string
	jenkinsPassword string
	jenkinsView     string

	currentIteration string
	interval         string

	displayBuild       bool
	displayAzureDevOps bool
	displayGitHub      bool

	debug bool
)

func init() {
	flag.StringVar(&travisToken, "travis-token", os.Getenv("TRAVIS_CI_TOKEN"), "The Travis CI authentication token (or define env var TRAVIS_CI_TOKEN)")
	flag.StringVar(&travisOwner, "travis-owner", os.Getenv("TRAVIS_CI_OWNER"), "The Travis CI owner (or define env var TRAVIS_CI_OWNER)")

	flag.StringVar(&azureDevOpsAccount, "azure-devops-account", os.Getenv("AZURE_DEVOPS_ACCOUNT"), "The Visual Studio Team Services account (or define env var AZURE_DEVOPS_ACCOUNT)")
	flag.StringVar(&azureDevOpsProject, "azure-devops-project", os.Getenv("AZURE_DEVOPS_PROJECT"), "The Visual Studio Team Services project (or define env var AZURE_DEVOPS_PROJECT)")
	flag.StringVar(&azureDevOpsTeam, "azure-devops-team", os.Getenv("AZURE_DEVOPS_TEAM"), "The Visual Studio Team Services team (or define env var AZURE_DEVOPS_TEAM)")
	flag.StringVar(&azureDevOpsToken, "azure-devops-token", os.Getenv("AZURE_DEVOPS_TOKEN"), "The Visual Studio Team Services auth token (or define env var AZURE_DEVOPS_TOKEN)")
	flag.IntVar(&azureDevOpsBuildCount, "azure-devops-build-count", 10, "How many builds should we display")
	flag.StringVar(&azureDevOpsBuildBranchFilter, "azure-devops-build-branch", "master", "Comma separated list of branches to display")

	flag.StringVar(&githubToken, "github-token", os.Getenv("GITHUB_TOKEN"), "The GitHub CI authentication token (or define env var GITHUB_TOKEN)")
	flag.StringVar(&githubOwner, "github-owner", os.Getenv("GITHUB_OWNER"), "The GitHub CI owner (or define env var GITHUB_OWNER)")

	flag.StringVar(&jenkinsURL, "jenkins-url", os.Getenv("JENKINS_URL"), "The Jenkins URL (or define env var JENKINS_URL)")
	flag.StringVar(&jenkinsUsername, "jenkins-username", os.Getenv("JENKINS_USERNAME"), "The Jenkins username to authenticate with (or define env var JENKINS_USERNAME)")
	flag.StringVar(&jenkinsPassword, "jenkins-password", os.Getenv("JENKINS_PASSWORD"), "The Jenkins password to authenticate with (or define env var JENKINS_PASSWORD)")
	flag.StringVar(&jenkinsView, "jenkins-view", os.Getenv("JENKINS_VIEW"), "The Jenkins view you want render, otherwise it is all (or define env var JENKINS_VIEW)")

	flag.StringVar(&currentIteration, "current-iteration", "", "What is the current iteration")
	flag.StringVar(&interval, "interval", "60s", "The refresh rate for the dashboard")

	flag.BoolVar(&displayBuild, "display-build", true, "Do you want to show build information from TravisCI and Jenkins?")
	flag.BoolVar(&displayAzureDevOps, "display-azure-devops", false, "Do you want to show Azure DevOps information?")
	flag.BoolVar(&displayGitHub, "display-github", true, "Do you want to show GitHub information?")

	flag.Usage = printUsage
	flag.Parse()

	if displayBuild && (travisToken == "" || travisOwner == "") {
		printUsage()
		os.Exit(1)
	}

	if displayAzureDevOps && (azureDevOpsAccount == "" || azureDevOpsProject == "" || azureDevOpsTeam == "" || azureDevOpsToken == "") {
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
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.precis")

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

	// Initialize termui.
	if err := termui.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "initializing termui failed: %v", err)
		os.Exit(2)
	}
	defer termui.Close()

	// It seems that if we don't pause on the first iteration
	// of this widget, we get a crash in docker.
	time.Sleep(1 * time.Second)

	termui.Body.Align()
	termui.Render(termui.Body)

	displayLoading()
	displayWidgets()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		ticker.Stop()
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		ticker.Stop()
		termui.StopLoop()
	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		displayWidgets()
	})

	// Update on an interval
	go func() {
		for range ticker.C {
			displayWidgets()
		}
	}()

	// Start the loop.
	termui.Loop()
}

func displayLoading() {
	body := termui.NewGrid()
	body.X = 0
	body.Y = 0
	body.BgColor = termui.ThemeAttr("bg")
	body.Width = termui.TermWidth()

	w := termui.NewPar("Loading all the data")
	w.Height = 3
	w.PaddingLeft = 1
	w.PaddingRight = 1
	w.TextFgColor = termui.ColorWhite
	w.BorderLabel = "Loading..."
	w.BorderLabelFg = termui.ColorGreen
	w.BorderFg = termui.ColorWhite

	body.AddRows(
		termui.NewRow(termui.NewCol(12, 0, w)),
	)

	body.Align()
	termui.Render(body)
}

func displayWidgets() {
	body := termui.NewGrid()
	body.X = 0
	body.Y = 0
	body.BgColor = termui.ThemeAttr("bg")
	body.Width = termui.TermWidth()

	// Date Widgets
	date := doDate()
	iterationName := doIterationName()
	body.AddRows(
		termui.NewRow(termui.NewCol(11, 0, date), termui.NewCol(1, 0, iterationName)),
	)

	// GitHub
	if displayGitHub {
		github, err := doGitHub()
		if err != nil {
			stop(fmt.Sprintf("failed to get GitHub information: %v", err))
		}

		body.AddRows(
			termui.NewRow(termui.NewCol(12, 0, github)),
		)
	}

	// Azure DevOps
	if displayAzureDevOps {
		builds, err := getAzureDevOpsBuilds()
		if err != nil {
			stop(fmt.Sprintf("failed to get Azure DevOps CI Builds information: %v", err))
		}

		pulls, err := getAzureDevOpsPulls()
		if err != nil {
			stop(fmt.Sprintf("failed to get Azure DevOps Pull Requests information: %v", err))
		}

		if len(pulls.Rows) > 0 {
			body.AddRows(
				termui.NewRow(
					termui.NewCol(12, 0, pulls),
				),
			)
		}

		if builds != nil {
			body.AddRows(
				termui.NewRow(
					termui.NewCol(12, 0, builds),
				),
			)
		}
	}

	// Travis
	if displayBuild {
		travis, err := doTravis()
		if err != nil {
			stop(fmt.Sprintf("failed to get Travis information: %v", err))
		}

		jenkins, err := doJenkins()
		if err != nil {
			stop(fmt.Sprintf("failed to get Jenkins information: %v", err))
		}

		body.AddRows(
			termui.NewRow(
				termui.NewCol(6, 0, travis),
				termui.NewCol(6, 0, jenkins),
			),
		)
	}

	body.Align()
	termui.Clear()
	termui.Render(body)
}

func stop(msg string) {
	termui.StopLoop()
	termui.Close()
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(2)
}
