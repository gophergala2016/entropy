package main

import (
	"fmt"
	"github.com/gophergala2016/entropy/net"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
)

type regConn struct {
	action string
	ws     *websocket.Conn
}

var conns = make(map[*websocket.Conn]interface{})
var connreg = make(chan regConn, 10)

func addConn(ws *websocket.Conn) {
	connreg <- regConn{"+", ws}
}

func rmConn(ws *websocket.Conn) {
	connreg <- regConn{"-", ws}
}

func PongServer(ws *websocket.Conn) {
	var msg net.Message
	addConn(ws)

	for {
		err := websocket.JSON.Receive(ws, &msg)
		if err == io.EOF {
			rmConn(ws)
			return
		}

		if msg.Msg != "ping" {
			websocket.JSON.Send(ws, net.Message{"no.", "1.0"})
		} else {
			websocket.JSON.Send(ws, net.Message{"pong", "1.0"})
		}
	}
}

func connRegistrator() {
	go func() {
		for {
			b := <-connreg
			if b.action == "+" {
				conns[b.ws] = nil
				fmt.Println("Client connected.")
			} else if b.action == "-" {
				delete(conns, b.ws)
				fmt.Println("Client disconnected.")
			}
		}
	}()
}

func main() {
	connRegistrator()
	http.Handle("/ping", websocket.Handler(PongServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe:" + err.Error())
	}
}
