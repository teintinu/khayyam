package internal

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"golang.org/x/net/websocket"
)

type WebTermTabAction struct {
	title  string
	icon   string
	id     string
	action func(*os.File, ...string)
}

type WebTermTab struct {
	id           string
	path         string
	title        string
	staticFolder string
	wsIn         bool
	actions      []*WebTermTabAction
	lastState    []string

	ws  *websocket.Conn
	pty *os.File

	routines *WebTermTabRoutines

	processConsoleOutput   func(line string, wtts *WebTermTabRoutines)
	lastConsoleOutputLines []string
}

type WebTermTabRoutines struct {
	sendToFrontEnd func(action string, args ...string)
	refreshState   func()
	setUnknow      func()
	setRunning     func()
	setSuccess     func(...string)
	setError       func(...string)
}

func (webterm *WebTerm) AddShell(
	path string, title string, staticFolder string, wsin bool,
	getCommand func() (*exec.Cmd, error),
	actions []*WebTermTabAction,
	processConsoleOutput func(line string, wtts *WebTermTabRoutines)) *WebTermTab {

	tabid := strings.ReplaceAll(path, "/", "_")

	routines := &WebTermTabRoutines{}
	tab := &WebTermTab{
		id:                   tabid,
		path:                 path,
		title:                title,
		actions:              actions,
		staticFolder:         staticFolder,
		wsIn:                 wsin,
		routines:             routines,
		processConsoleOutput: processConsoleOutput,
	}
	routines.sendToFrontEnd = func(action string, args ...string) {
		webterm.sendToFrontEnd(tabid, action, args...)
	}
	routines.refreshState = func() {
		webterm.sendToFrontEnd(tabid, "refreshState", tab.lastState...)
	}

	routines.setUnknow = func() {
		tab.lastState = []string{"unknow"}
		routines.refreshState()
	}
	routines.setRunning = func() {
		tab.lastState = []string{"running"}
		routines.refreshState()
	}
	routines.setSuccess = func(args ...string) {
		n := []string{"success"}
		n = append(n, args...)
		tab.lastState = n
		routines.refreshState()
	}
	routines.setError = func(args ...string) {
		n := []string{"error"}
		n = append(n, args...)
		tab.lastState = n
		routines.refreshState()
	}

	webterm.tabs = append(webterm.tabs, tab)

	if staticFolder == "" {
		tab.tabHandlerForWS(webterm)
	} else {
		tab.tabHandlerForStaticFolder(webterm, staticFolder)
	}
	go tab.Pty(webterm, getCommand)

	return tab
}

func (tab *WebTermTab) Pty(
	webterm *WebTerm,
	getCommand func() (*exec.Cmd, error),
) {

	cmd, err := getCommand()
	if err != nil {
		tab.consoleOutput(fmt.Sprintf("Error getting pty command: %s\r\n", err))
		return
	}

	tab.pty, err = pty.Start(cmd)
	if err != nil {
		tab.consoleOutput(fmt.Sprintf("Error creating pty: %s\r\n", err))
		return
	}

	cmdConsole := bufio.NewReader(tab.pty)
	for !webterm.IsClosed() {
		line, err := cmdConsole.ReadString(linefeedDelimiter)
		if err != nil {
			tab.consoleOutput(fmt.Sprintf("Error reading from pty: %s\r\n", err))
			return
		}
		tab.consoleOutput(line)
	}
}

func (tab *WebTermTab) consoleOutput(line string) {
	if tab.ws != nil {
		if _, err := tab.ws.Write(append([]byte(line), linefeedDelimiter)); err != nil {
			tab.ws = nil
		}
	}
	tab.processConsoleOutput(line, tab.routines)
	var limit = len(tab.lastConsoleOutputLines) - 100
	if limit < 0 {
		limit = 0
	}
	b := append(tab.lastConsoleOutputLines, line)[limit:]
	tab.lastConsoleOutputLines = b
}

func (tab *WebTermTab) tabHandlerForWS(
	webterm *WebTerm,
) {

	webterm.Handle(tab.path, "websocket for "+tab.title, websocket.Handler(func(ws *websocket.Conn) {

		tab.ws = ws

		for _, line := range tab.lastConsoleOutputLines {
			ws.Write(append([]byte(line), linefeedDelimiter))
		}
		if tab.wsIn {
			io.Copy(tab.pty, ws)
		}

	}))

}

func (tab *WebTermTab) tabHandlerForStaticFolder(
	webterm *WebTerm,
	staticFolder string,
) {
	fileServer := http.FileServer(http.Dir(staticFolder))
	webterm.Handle(tab.path, "fileserver "+staticFolder,
		http.StripPrefix(tab.path, fileServer))
}
