package main

import (
	"fmt"
	"github.com/gophergala2016/entropy/models"
	"github.com/gophergala2016/entropy/net"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"os/signal"
	"os/user"
	"time"
)

func handleInterruptSignal(ws *websocket.Conn) {
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	ws.Close()
	log.Println("Connection closed")
	log.Println("Program killed")
	os.Exit(0)
}
func main() {
	var username string
	if len(os.Args) > 1 {
		username = os.Args[1]
	} else if utmp, err := user.Current(); err == nil {
		username = utmp.Username
	} else {
		username = "DefaultUser"
	}

	origin := "http://localhost/"
	url := "ws://localhost:12345/game"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	go handleInterruptSignal(ws)

	mesg := net.Connection{username}
	err = websocket.JSON.Send(ws, mesg)
	if err != nil {
		log.Println("Couldn't send connection message", err)
		return
	}
	for {
		websocket.JSON.Send(ws, net.GetUserList{State: models.StateConnected})
		var msg net.Message
		websocket.JSON.Receive(ws, &msg)
		switch {
		case msg.ResponseUserList != nil:
			for k, _ := range msg.ResponseUserList.GamePlayers {
				fmt.Println("User: ", k)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
