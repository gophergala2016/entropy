package main

import (
	"fmt"
	"github.com/gophergala2016/entropy/net"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:12345/ping"

	ws, err := websocket.Dial(url, "", origin)
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		ws.Close()
		log.Println("Connection closed")
		log.Println("Program killed")
		os.Exit(0)
	}()

	defer ws.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		mesg := net.Message{os.Args[1], "1.0"}
		websocket.JSON.Send(ws, mesg)
		websocket.JSON.Receive(ws, &mesg)
		fmt.Println("msg: ", mesg)
		time.Sleep(1 * time.Second)
	}
}
