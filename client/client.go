package main

import (
	"github.com/gophergala2016/entropy/models"
	"github.com/gophergala2016/entropy/net"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"os/user"
	"time"
)

func handleInterruptSignal(ws *websocket.Conn, done chan struct{}) {
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
	ws.Close()
	log.Println("Connection closed")
	log.Println("Program killed")
	os.Exit(0)
	return
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	var username string
	if len(os.Args) > 1 {
		username = os.Args[1]
	} else if utmp, err := user.Current(); err == nil {
		username = utmp.Username
	} else {
		username = "DefaultUser"
	}

	url := "ws://localhost:12345/game"

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	//	go handleInterruptSignal(ws, done)

	go func() {
		defer ws.Close()
		defer close(done)

		for {
			var msg interface{}
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
			}
			switch m := msg.(type) {
			case net.GetUserList:
				log.Println("Getuserlist", m)
			default:
				log.Println("received", msg)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	err = ws.WriteJSON(net.Connection{username})
	if err != nil {
		log.Println("Couldn't send connection message", err)
		return
	}

	for {
		select {
		case <-ticker.C:
			st := models.StateConnected
			err := ws.WriteJSON(net.GetUserList{State: &st})
			if err != nil {
				log.Println("write JSON:", err)
				return
			}
		case <-sigchan:

			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			ws.Close()
			log.Println("Connection closed")
			log.Println("Program killed")
		}
	}
	/*
		mesg := net.Connection{username}
		err = websocket.JSON.Send(ws, mesg)
		if err != nil {
			log.Println("Couldn't send connection message", err)
			return
		}
		for {
			websocket.JSON.Send(ws, net.GetUserList{State: models.StateConnected})
			var msg models.GamePlayers //interface{}
			websocket.JSON.Receive(ws, &msg)
			//switch m := msg.(type) {
			//case models.GamePlayers:
			for k, _ := range msg {
				fmt.Println("User: ", k)
			}
			//}
			time.Sleep(1 * time.Second)
		}
	*/
}
