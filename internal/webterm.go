package internal

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type WebTerm struct {
	closed bool
	server *http.Server
	mux    *http.ServeMux

	toFrontend *websocket.Conn

	tabs []*WebTermTab
}

func NewWebTerm() *WebTerm {
	return &WebTerm{
		mux: http.NewServeMux(),
	}
}

func (webterm *WebTerm) Close() {
	webterm.closed = true
}

func (webterm *WebTerm) IsClosed() bool {
	return webterm.closed
}

func (webterm *WebTerm) Handle(pattern string, handler http.Handler) {
	webterm.mux.Handle(pattern, handler)
}

func (webterm *WebTerm) Start(port int32) {

	webterm.Handle("/", webtermHome(webterm.tabs))
	webterm.Handle("/_tab", webtermTab())
	webterm.Handle("/_comm", webtermComm(webterm))

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

func (webterm *WebTerm) processCommand(command []string) {
	tabId := command[0]
	actionId := command[1]
	if tabId == "*" {
		if actionId == "exit" {
			os.Exit(0)
		}
	} else {
		for _, tab := range webterm.tabs {
			if tab.id == tabId {
				for _, ac := range tab.actions {
					if ac.id == actionId {
						ac.action(tab.pty, command[2:]...)
					}
				}
			}
		}
	}
}

func (webterm *WebTerm) sendToFrontEnd(command ...string) {
	if webterm.toFrontend != nil {
		s := strings.Join(command, "\v")
		println("toFrontend:", command[0], command[1], command[2])
		if _, err := webterm.toFrontend.Write([]byte(s + "\n")); err != nil {
			webterm.toFrontend = nil
			fmt.Println("perdeu conexao:", err)
		}
	} else {
		println("NAO CONECTADO toFrontend:", command[0], command[1], command[2])
	}
}

var RegRemoveANSI = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")
