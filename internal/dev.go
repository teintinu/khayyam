package internal

import (
	"os"
	"os/exec"
	"regexp"
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
			Colors:   true,
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
	rgTestSummary := regexp.MustCompile(`^.*Tests:.+(?:(\d+)\s+failed,)?.+(\d+).+passed,.+(\d+).+total.*$`)
	rgRUNS := regexp.MustCompile(`^.*RUNS.*\.\.\..*$`)
	processConsoleOutput := func(line string, wtts *WebTermTabRoutines) {
		println("dev " + line)
		// Tests:       5 passed, 5 total
		// Tests:       1 failed, 4 passed, 5 total
		testResult := rgTestSummary.FindStringSubmatch(line)
		if len(testResult[0]) > 0 {
			failed := testResult[1]
			passed := testResult[2]
			total := testResult[3]
			if failed == "" {
				wtts.setSuccess(line, failed, passed, total)
			} else {
				wtts.setError(line, failed, passed, total)
			}
		} else if strings.Contains(line, "No tests found related to files changed since last commit.") {
			wtts.setSuccess(line)
		} else if rgRUNS.MatchString(line) {
			wtts.setBusy()
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
