package internal

import (
	"os"
	"os/exec"
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
		println("a:", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(lineWithColors, "\b", "BS"), "\r", "CR"), "\x1b", "ESC"))
		line := RegRemoveANSI.ReplaceAllString(lineWithColors, "")
		println("b:", line)
		testResult := RegExpJestSummary.FindStringSubmatch(line)

		if len(testResult) > 0 {
			failed := testResult[1]
			passed := testResult[2]
			total := testResult[3]
			if failed == "" {
				wtts.setSuccess(line, failed, passed, total)
			} else {
				wtts.setError(line, failed, passed, total)
			}
		} else if strings.Contains(line, "No tests found related to files changed since last commit.") {
			wtts.setUnknow()
		} else if strings.Contains(line, "Determining test suites to run") {
			wtts.setRunning()
		} else if RegExpJestRunning.MatchString(line) {
			wtts.setRunning()
		}
	}
	webterm.AddShell("/jest", "Tests", false, getCommand, actions, processConsoleOutput)
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
