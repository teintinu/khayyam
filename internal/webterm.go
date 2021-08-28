package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/creack/pty"
	"golang.org/x/net/websocket"
)

type WebTerm struct {
	closed bool
	server *http.Server
	mux    *http.ServeMux

	tabs []*WebTermTab
}

const tcExecuteAllTests = "a"

type WebTermTabStateEnum uint8

const (
	TabRunning WebTermTabStateEnum = iota
	TabSuccess
	TabError
)

type WebTermTab struct {
	id      string
	path    string
	title   string
	state   WebTermTabStateEnum
	actions []*WebTermTab

	ws        *websocket.Conn
	pty       *os.File
	ptyBuffer bytes.Buffer
	ptyLock   sync.Mutex
}

type WebTermTabAction struct {
	title  string
	icon   string
	action string
}

type WebTermTabState struct {
	setBusy    func()
	setSuccess func()
	setError   func(string)
	input      func(string)
	notify     func(...string)
}

func NewWebTerm() *WebTerm {
	return &WebTerm{
		mux: http.NewServeMux(),
	}
}

func (webterm *WebTerm) Close() {
	webterm.closed = true
}

func (webterm *WebTerm) Handle(pattern string, handler http.Handler) {
	webterm.mux.Handle(pattern, handler)
}

func (webterm *WebTerm) AddShell(
	path string, title string,
	getCommand func() (*exec.Cmd, error),
	processAction func(ac []string, wtts *WebTermTabState),
	processNotification func(line string, wtts *WebTermTabState)) *WebTermTab {

	tab := newWebTermTab(path, title)

	webterm.tabs = append(webterm.tabs, tab)

	tab.handle(webterm, getCommand, processAction, processNotification)

	return tab
}

func (webterm *WebTerm) Start(port int32) {

	webterm.Handle("/", home(webterm.tabs))
	webterm.Handle("/_tab", tab())

	webterm.server = &http.Server{
		Addr:           ":" + fmt.Sprint(port),
		Handler:        webterm.mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("dev server started at http://localhost:" + strconv.Itoa(DevPort))
	webterm.server.ListenAndServe()

}

func newWebTermTab(path string, title string) *WebTermTab {

	tab := &WebTermTab{
		id:    strings.ReplaceAll(path, "/", "_"),
		path:  path,
		title: title,

		ptyBuffer: *bytes.NewBuffer([]byte{}),
	}
	return tab

}

func (tab *WebTermTab) handle(webTerm *WebTerm,
	getCommand func() (*exec.Cmd, error),
	processAction func(ac []string, wtts *WebTermTabState),
	processNotification func(line string, wtts *WebTermTabState)) {

	webTerm.Handle(tab.path, websocket.Handler(func(ws *websocket.Conn) {

		tab.ws = ws
		defer func() {
			tab.ws = nil
		}()

		// cmd.SysProcAttr = &syscall.SysProcAttr{
		// 	Setsid: true,
		// }
		// syscall.Setpgid(jestCmd.Process.Pid, jestCmd.Process.Pid)

		cmd, err := getCommand()
		if err != nil {
			ws.Write([]byte(fmt.Sprintf("Error creating pty: %s\r\n", err)))
			ws.Close()
			return
		}

		// cmd.SysProcAttr = &syscall.SysProcAttr{
		// 	Setctty: true,
		// }
		tab.pty, err = pty.Start(cmd)
		if err != nil {
			ws.Write([]byte(fmt.Sprintf("Error creating pty: %s\r\n", err)))
			ws.Close()
			return
		}

		go io.Copy(ws, tab.pty)
		io.Copy(tab.pty, ws)

		// cmdWriter := bufio.NewWriter(cmdOutput)
		// go func() {
		// 	for !webterm.closed {
		// 		s := <-cmdStdIn
		// 		cmdWriter.WriteString(s)
		// 	}
		// }()

		// cmdConsole := bufio.NewReader(cmdOutput)
		// wsWriter := bufio.NewWriter(ws)
		// for !webterm.closed {
		// 	line, err := cmdConsole.ReadString(linefeedDelimiter)
		// 	if err != nil {
		// 		ws.Write([]byte(fmt.Sprintf("Error creating pty: %s\r\n", err)))
		// 		ws.Close()
		// 		return
		// 	}
		// 	cmdStdOut <- line
		// 	wsWriter.WriteString(line)
		// 	wsWriter.WriteByte(linefeedDelimiter)
		// }

	}))

}

// webterm.mux.Handle("/_testcomm", websocket.Handler(func(ws *websocket.Conn) {

// 	wsReceiveCommand := bufio.NewReader(ws)
// 	wsSendNotification := bufio.NewWriter(ws)

// 	go func() {
// 		for !webterm.closed {
// 			s, err := wsReceiveCommand.ReadString(10)
// 			if err == nil {
// 				sa := strings.Split(s, "\u2800")
// 				if !webterm.internalCommand(sa) {
// 					webterm.commandChannel <- sa
// 				}
// 			}
// 		}
// 	}()
// 	for !webterm.closed {
// 		s := <-webterm.notifyChannel
// 		wsSendNotification.WriteString(strings.Join(s, "\u2800"))
// 	}

// }))

func (tab *WebTermTab) ResizeTerminal(width int, height int) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(height),
		uint16(width),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		tab.pty.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return errno
	} else {
		return nil
	}
}

func (tab *WebTermTab) internalCommand(s []string) bool {
	if s[0] == "$resize" {
		w, errw := strconv.ParseInt(s[1], 10, 8)
		h, errh := strconv.ParseInt(s[2], 10, 8)
		if errw == nil && errh == nil {
			tab.ResizeTerminal(int(w), int(h))
		}
		return true
	}
	return false
}
