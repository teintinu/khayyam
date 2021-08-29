package internal

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type WebTerm struct {
	closed bool
	server *http.Server
	mux    *http.ServeMux

	toFrontend chan []string

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
