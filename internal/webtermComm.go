package internal

import (
	"bufio"
	"strings"

	"golang.org/x/net/websocket"
)

func webtermComm(webterm *WebTerm) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {

		wsFromFrontend := bufio.NewReader(ws)
		wsToFrontEnd := bufio.NewWriter(ws)

		go func() {
			for !webterm.IsClosed() {
				s, err := wsFromFrontend.ReadString(10)
				if err == nil {
					sa := strings.Split(s, "\v")
					go webterm.processCommand(sa)
				}
			}
		}()
		for !webterm.closed {
			s := <-webterm.toFrontend
			wsToFrontEnd.WriteString(strings.Join(s, "\v") + "\n")
		}

	})
}
