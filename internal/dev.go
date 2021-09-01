package internal

import (
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/evanw/esbuild/pkg/api"
)

const DevPort = 37000

func Dev(repo *Repository) error {
	wg := &sync.WaitGroup{}

	webterm := NewWebTerm()

	wg.Add(1)
	devRunExecutables(repo, wg, webterm)
	devTestApps(repo, wg, webterm)

	webterm.Start(DevPort)
	wg.Wait()

	return nil
}

func devTestApps(repo *Repository, wg *sync.WaitGroup, webterm *WebTerm) {

	jestReport := ""
	total := ""
	failed := ""

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

	processConsoleOutput := func(lineWithColors string, wtts *WebTermTabRoutines) {
		if strings.Contains(lineWithColors, "\x1b[K") {
			wtts.setRunning()
			return
		}
		// println(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(lineWithColors, "\b", "BS"), "\r", "CR"), "\x1b", "ESC"))
		// println("a:", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(lineWithColors, "\b", "BS"), "\r", "CR"), "\x1b", "ESC"))
		line := RegRemoveANSI.ReplaceAllString(lineWithColors, "")
		testResult := RegExpJestRunComplete.FindStringSubmatch(line)
		//fmt.Println()
		//fmt.Println("b:", line)
		//fmt.Println()
		//fmt.Println(testResult)

		if len(testResult) > 0 {
			total = testResult[3]
			failed = testResult[1]
		} else if strings.Contains(line, "Ran all test suites") {
			if total == "0" {
				wtts.setUnknow()
			} else if failed == "" {
				wtts.setSuccess(line, total, failed)
			} else {
				wtts.setError(line, total, failed)
			}
		} else if RegExpJestReportCreated.MatchString(line) {
			wtts.sendToFrontEnd("reloadTab")
		}
	}
	webterm.AddShell("/jest", "Tests", jestReport, true, getCommand, actions, processConsoleOutput)

	go func() {
		for !webterm.IsClosed() {
			time.Sleep(time.Second)
		}
		wg.Done()
	}()
}

func devRunExecutables(repo *Repository, wg *sync.WaitGroup, webterm *WebTerm) {
	for _, pkg := range repo.Packages {
		wg.Add(1)
		BundleWithEsbuild(repo, pkg, &BuildOpts{
			Target: api.ESNext,
			Minify: false,
			Mode:   WatchAndRun,
		})
	}
	wg.Done()
}
