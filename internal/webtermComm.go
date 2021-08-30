package internal

import (
	"bufio"
	"strings"

	"golang.org/x/net/websocket"
)

func webtermComm(webterm *WebTerm) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {

		webterm.toFrontend = ws
		wsFromFrontend := bufio.NewReader(ws)

		for _, tab := range webterm.tabs {
			tab.routines.refreshState()
		}

		for !webterm.IsClosed() {
			s, err := wsFromFrontend.ReadString(10)
			if err == nil {
				sa := strings.Split(s, "\v")
				go webterm.processCommand(sa)
			}
		}

	})
}
