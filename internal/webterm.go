package internal

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/kr/pty"
	"golang.org/x/net/websocket"
)

type WebTerm struct {
	closed             bool
	server             *http.Server
	mux                *http.ServeMux
	tabs               []*WebTermTab
	testNotifyChannel  chan string
	testCommandChannel chan string
}

const tcExecuteAllTests = "ExecuteAllTests"

type WebTermTab struct {
	path  string
	title string
}

func (webterm *WebTerm) Close() {
	webterm.closed = true
}

func (webterm *WebTerm) AddIFrame(path string, title string) {
	webterm.tabs = append(webterm.tabs, &WebTermTab{
		path:  path,
		title: title,
	})
}

func (webterm *WebTerm) AddShell(path string, title string, cmd *exec.Cmd) (chan string, chan string) {

	var cmdStdIn chan string
	var cmdStdOut chan string

	webterm.tabs = append(webterm.tabs, &WebTermTab{
		path:  path,
		title: title,
	})

	webterm.mux.Handle(path, websocket.Handler(func(ws *websocket.Conn) {

		defer ws.Close()
		cmdOutput, err := pty.Start(cmd)
		if err != nil {
			ws.Write([]byte(fmt.Sprintf("Error creating pty: %s\r\n", err)))
			ws.Close()
			return
		}

		go io.Copy(ws, cmdOutput)

		cmdWriter := bufio.NewWriter(cmdOutput)
		go func() {
			for !webterm.closed {
				s := <-cmdStdIn
				cmdWriter.WriteString(s)
			}
		}()

		cmdConsole := bufio.NewReader(cmdOutput)
		wsWriter := bufio.NewWriter(ws)
		for !webterm.closed {
			line, err := cmdConsole.ReadString(linefeedDelimiter)
			if err != nil {
				ws.Write([]byte(fmt.Sprintf("Error creating pty: %s\r\n", err)))
				ws.Close()
				return
			}
			cmdStdOut <- line
			wsWriter.WriteString(line)
			wsWriter.WriteByte(linefeedDelimiter)
		}

	}))
	return cmdStdIn, cmdStdOut
}

func (webterm *WebTerm) Start(port int32) error {

	webterm.mux.Handle("/", home(webterm.tabs))
	webterm.mux.Handle("/_tab", tab())

	webterm.mux.Handle("/_testcomm", websocket.Handler(func(ws *websocket.Conn) {

		wsReceiveCommand := bufio.NewReader(ws)
		wsSendNotification := bufio.NewWriter(ws)

		go func() {
			s, err := wsReceiveCommand.ReadString(10)
			if err == nil {
				webterm.testCommandChannel <- s
			}
		}()
		for !webterm.closed {
			s := <-webterm.testNotifyChannel
			wsSendNotification.WriteString(s)
		}

	}))

	webterm.server = &http.Server{
		Addr:           ":" + fmt.Sprint(port),
		Handler:        webterm.mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return webterm.server.ListenAndServe()

}

func home(tabs []*WebTermTab) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html>
<head>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.js" integrity="sha512-2PRgAav8Os8vLcOAh1gSaDoNLe1fAyq8/G3QSdyjFFD+OqNjLeHE/8q4+S4MEZgPsuo+itHopj+hJvqS8XUQ8A==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.css" integrity="sha512-iLYuqv+v/P4u9erpk+KM83Ioe/l7SEmr7wB6g+Kg1qmEit8EShDKnKtLHlv2QXUp7GGJhmqDI+1PhJYLTsfb8w==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<style>
		#tabs {
			
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
<div class="tabs">
`)
		for _, tab := range tabs {
			onClick := `openframe('` + tab.path + `')`
			fmt.Fprintln(w, `<span class="tab" onClick="`+onClick+`">`+tab.title+`</span>`)
		}
		fmt.Fprintln(w, `</div>`)
		fmt.Fprintln(w, `</iframe id='frame' src="`+tabs[0].path+`">`)
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
		#terminal {
	    background: #000000;
		  color: #ffffff;
		  display: inline-block;
		  padding: 10px;
		}
	</style>
</head>
<body>

	<pre id="terminal"></pre>

	<script>

		var elem = document.getElementById("terminal");
		elem.tabindex = 0;

		var terminal = new Terminal();
		var input = terminal.dom(elem);

		var socket = new WebSocket('ws://'+document.location.host+`+q+`, 'echo');

		socket.addEventListener("open", function () {
				input.on('data', function (evt) {
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
