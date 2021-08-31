package internal

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
)

const DevPort = 37000

func Dev(repo *Repository) error {
	wg := &sync.WaitGroup{}

	webterm := NewWebTerm()

	wg.Add(2)
	devTestApps(repo, wg, webterm)
	// devRunExecutables(repo, wg, webterm)

	webterm.Start(DevPort)
	wg.Wait()

	return nil
}

func devTestApps(repo *Repository, wg *sync.WaitGroup, webterm *WebTerm) {

	jestReport := path.Join(repo.RootDir, ".jest/html-report")

	getCommand := func() (*exec.Cmd, error) {
		return CreateJestCommand(repo, TestOptions{
			Watch:    true,
			Coverage: false,
			Colors:   false,
		}), nil
	}
	executeAllTests := &WebTermTabAction{
		id:    "executeAllTests",
		title: "Run All Tests",
		icon:  "play",
		action: func(pty *os.File, args ...string) {
			pty.WriteString("a")
		},
	}
	actions := []*WebTermTabAction{
		executeAllTests,
	}

	processConsoleOutput := func(lineWithColors string, wtts *WebTermTabRoutines) {
		// println(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(lineWithColors, "\b", "BS"), "\r", "CR"), "\x1b", "ESC"))
		// println("a:", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(lineWithColors, "\b", "BS"), "\r", "CR"), "\x1b", "ESC"))
		line := RegRemoveANSI.ReplaceAllString(lineWithColors, "")
		testResult := RegExpJestRunComplete.FindStringSubmatch(line)
		// fmt.Println("b:", line)
		// fmt.Println()
		// fmt.Println()
		// fmt.Println(testResult)

		if len(testResult) > 0 {
			total := testResult[1]
			failed := testResult[2]
			if total == "0" {
				wtts.setUnknow()
			} else if failed == "0" {
				wtts.setSuccess(line, total, failed)
			} else {
				wtts.setError(line, total, failed)
			}
		} else if strings.Contains(line, "No tests found related to files changed since last commit.") {
			wtts.setUnknow()
		} else if RegExpJestRunStart.MatchString(line) {
			wtts.setRunning()
		} else if RegExpJestReportCreated.MatchString(line) {
			wtts.sendToFrontEnd("reloadTab")
		}
	}
	webterm.AddShell("/jest", "Tests", jestReport, true, getCommand, actions, processConsoleOutput)
}

func devRunExecutables(repo *Repository, wg *sync.WaitGroup, webterm *WebTerm) {
	for _, pkg := range repo.Packages {
		BundleWithEsbuild(repo, pkg, &BuildOpts{
			Target: api.ESNext,
			Minify: false,
			Mode:   WatchAndRun,
		})
	}
	wg.Done()
}
