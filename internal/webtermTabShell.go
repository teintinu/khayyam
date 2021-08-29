package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"golang.org/x/net/websocket"
)

type WebTermTabStateEnum uint8

const (
	TabRunning WebTermTabStateEnum = iota
	TabSuccess
	TabError
)

type WebTermTabAction struct {
	title  string
	icon   string
	id     string
	action func(*os.File, ...string)
}

type WebTermTab struct {
	id       string
	path     string
	title    string
	readonly bool
	state    WebTermTabStateEnum
	actions  []*WebTermTabAction

	ws                     *websocket.Conn
	pty                    *os.File
	consoleOutput          chan string
	lastConsoleOutputLines []string

	routines *WebTermTabRoutines
}

type WebTermTabRoutines struct {
	setBusy    func()
	setSuccess func(...string)
	setError   func(...string)
}

func (webterm *WebTerm) AddShell(
	path string, title string, readonly bool,
	getCommand func() (*exec.Cmd, error),
	actions []*WebTermTabAction,
	processConsoleOutput func(line string, wtts *WebTermTabRoutines)) *WebTermTab {

	tabid := strings.ReplaceAll(path, "/", "_")

	routines := &WebTermTabRoutines{
		setBusy: func() {
			webterm.toFrontend <- []string{tabid, "busy"}
		},
		setSuccess: func(args ...string) {
			n := []string{tabid, "success"}
			n = append(n, args...)
			webterm.toFrontend <- n
		},
		setError: func(args ...string) {
			n := []string{tabid, "error"}
			n = append(n, args...)
			webterm.toFrontend <- n
		},
	}
	tab := &WebTermTab{
		id:       tabid,
		path:     path,
		title:    title,
		readonly: readonly,
		routines: routines,
	}

	webterm.tabs = append(webterm.tabs, tab)

	tab.tabHandler(webterm)
	go tab.Pty(webterm, getCommand)

	go func() {
		for !webterm.IsClosed() {
			var line string = <-tab.consoleOutput
			processConsoleOutput(line, tab.routines)
			var limit = len(tab.lastConsoleOutputLines)
			if limit > 100 {
				limit = 100
			}
			a := []string{line}
			b := append(a, tab.lastConsoleOutputLines...)
			tab.lastConsoleOutputLines = b[:limit]
		}
	}()
	return tab
}

func (tab *WebTermTab) Pty(
	webterm *WebTerm,
	getCommand func() (*exec.Cmd, error),
) {

	cmd, err := getCommand()
	println("cmd", cmd)
	if err != nil {
		tab.consoleOutput <- fmt.Sprintf("Error getting pty command: %s\r\n", err)
		return
	}

	tab.pty, err = pty.Start(cmd)
	if err != nil {
		tab.consoleOutput <- fmt.Sprintf("Error creating pty: %s\r\n", err)
		return
	}

	cmdConsole := bufio.NewReader(tab.pty)
	for !webterm.IsClosed() {
		line, err := cmdConsole.ReadString(linefeedDelimiter)
		println("console", line)
		if err != nil {
			tab.consoleOutput <- fmt.Sprintf("Error reading from pty: %s\r\n", err)
			return
		}
		tab.consoleOutput <- line
	}

}

func (tab *WebTermTab) tabHandler(
	webterm *WebTerm,
) {

	webterm.Handle(tab.path, websocket.Handler(func(ws *websocket.Conn) {

		tab.ws = ws
		defer func() {
			tab.ws = nil
		}()

		for _, line := range tab.lastConsoleOutputLines {
			tab.ws.Write([]byte(line))
			tab.ws.Write([]byte{linefeedDelimiter})
		}
		if !tab.readonly {
			go io.Copy(ws, tab.pty)
		}

		for !webterm.IsClosed() {
			line := <-tab.consoleOutput
			tab.ws.Write([]byte(line))
			tab.ws.Write([]byte{linefeedDelimiter})
		}

	}))

}
