package internal

import (
	"fmt"
	"os/exec"
	"regexp"
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
	processAction := func(ac []string, wtts *WebTermTabState) {
		fmt.Println(ac)
		if ac[0] == tcExecuteAllTests {
			wtts.input("a")
		}
	}
	rgTestSummary := regexp.MustCompile(`^.*Tests:.+(?:(\d+)\s+failed,)?.+(\d+).+passed,.+(\d+).+total.*$`)
	processNotification := func(line string, wtts *WebTermTabState) {
		// Tests:       5 passed, 5 total
		// Tests:       1 failed, 4 passed, 5 total
		testResult := rgTestSummary.FindStringSubmatch(line)
		if len(testResult[0]) > 0 {
			failed := testResult[1]
			passed := testResult[2]
			total := testResult[3]
			wtts.notify("testResult", failed, passed, total)
		}
	}
	webterm.AddShell("/jest", "Tests", getCommand, processAction, processNotification)
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
