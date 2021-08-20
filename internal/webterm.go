package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

type WebTermTab struct {
	path  string
	title string

	ws        *websocket.Conn
	pty       *os.File
	ptyBuffer bytes.Buffer
	ptyLock   sync.Mutex
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

func home(tabs []*WebTermTab) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html>
<head>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.js" integrity="sha512-2PRgAav8Os8vLcOAh1gSaDoNLe1fAyq8/G3QSdyjFFD+OqNjLeHE/8q4+S4MEZgPsuo+itHopj+hJvqS8XUQ8A==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.css" integrity="sha512-iLYuqv+v/P4u9erpk+KM83Ioe/l7SEmr7wB6g+Kg1qmEit8EShDKnKtLHlv2QXUp7GGJhmqDI+1PhJYLTsfb8w==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<style>
		* {
			padding: 0px;
			margin: 0px;
		}
		document, body {
			width: 100%;
			height: 100%;
		}
  	#app {
			display: flex;			
			flex-direction: column;
			width: 100%;
			height: 100%;
		}
		#frame {
			flex-grow: 1;
		}
	</style>
	<script>
	  function openframe(path) {
			var frame=document.getElementById('frame')
			frame.src = '/_tab?q=' + path
		}
	</script>
</head>
<body>
<div id="app">
`)
		fmt.Fprintln(w, `<div class="tabs">`)
		for _, tab := range tabs {
			onClick := `openframe('` + tab.path + `')`
			fmt.Fprintln(w, `<span class="tab" onClick="`+onClick+`">`+tab.title+`</span>`)
		}
		fmt.Fprintln(w, `</div>`)
		fmt.Fprintln(w, `<iframe id='frame' src="/_tab?q=`+tabs[0].path+`" />`)
		fmt.Fprint(w, `</div></body></html>`)
	}
}

func tab() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		q := r.FormValue("q")
		fmt.Fprintln(w, `<html>
<head>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.js" integrity="sha512-2PRgAav8Os8vLcOAh1gSaDoNLe1fAyq8/G3QSdyjFFD+OqNjLeHE/8q4+S4MEZgPsuo+itHopj+hJvqS8XUQ8A==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.css" integrity="sha512-iLYuqv+v/P4u9erpk+KM83Ioe/l7SEmr7wB6g+Kg1qmEit8EShDKnKtLHlv2QXUp7GGJhmqDI+1PhJYLTsfb8w==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<style>
	  * {
			padding: 0px;
			margin: 0px;
		}
	  document, body {
      width: 100%;
      height: 100%;
  	}
		#terminal {
	    background: #000000;
		  color: #ffffff;
		  display: inline-block;
		  padding: 10px;
      width: 100%;
      height: 100%;
		}
	</style>
</head>
<body>

	<pre id="terminal"></pre>

	<script>
	
	debugger
		var elem = document.getElementById("terminal");
		elem.tabindex = 0;

		var terminal = new Terminal();
		terminal.open(elem);

		var socket = new WebSocket('ws://'+document.location.host+"`+q+`", 'echo');

		socket.addEventListener("open", function () {
			terminal.on('data', function (evt) {
						socket.send(evt);
		 		});
		 });

		socket.addEventListener("message", function (evt) {
				terminal.write(event.data);
		});
	</script>
</body>
</html>
`)
	}
}

func newWebTermTab(path string, title string) *WebTermTab {

	tab := &WebTermTab{
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
