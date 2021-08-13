package internal

import (
	"regexp"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
)

const DevPort = 37000

func Dev(repo *Repository) error {
	wg := &sync.WaitGroup{}

	webterm := &WebTerm{}

	wg.Add(2)
	go devTestApps(repo, wg, webterm)
	go devRunExecutables(repo, wg, webterm)

	webterm.Start(DevPort)
	wg.Wait()

	return nil
}

func devTestApps(repo *Repository, wg *sync.WaitGroup, webterm *WebTerm) {
	jestCmd := CreateJestCommand(repo, TestOptions{
		Watch:    true,
		Coverage: false,
		Colors:   true,
	})
	jestStdIn, jestStdOut := webterm.AddShell("/jest", "Tests", jestCmd)

	go func() {
		for {
			ac := <-webterm.testCommandChannel
			if ac == tcExecuteAllTests {
				jestStdIn <- "a"
			}
		}
	}()
	rgTestSummary := regexp.MustCompile(`^.*Tests:.+(?:(\d+)\s+failed,)?.+(\d+).+passed,.+(\d+).+total.*$`)
	for !webterm.closed {
		line := <-jestStdOut
		// Tests:       5 passed, 5 total
		// Tests:       1 failed, 4 passed, 5 total
		testResult := rgTestSummary.FindStringSubmatch(line)
		if len(testResult[0]) > 0 {
			failed := testResult[1]
			passed := testResult[1]
			total := testResult[1]
			webterm.testNotifyChannel <- "testResult," + failed + "," + passed + "," + total
		}
	}
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
