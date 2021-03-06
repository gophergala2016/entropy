package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gophergala2016/entropy/models"
	"github.com/gophergala2016/entropy/net"
	"golang.org/x/net/websocket"
)

type regConn struct {
	action   string
	username string
	ws       *websocket.Conn
}

type registrationFeedback struct {
	action   string
	username string
	done     bool
}

var regFeedback = make(chan registrationFeedback, 10)

var conns = make(map[*websocket.Conn]interface{})
var connreg = make(chan regConn, 10)

func addConn(username string, ws *websocket.Conn) {
	connreg <- regConn{"+", username, ws}
}

func rmConn(username string, ws *websocket.Conn) {
	connreg <- regConn{"-", username, ws}
}

func GameServer(ws *websocket.Conn) {
	var cmsg net.Connection
	var curruser string

	err := websocket.JSON.Receive(ws, &cmsg)
	if err != nil {
		log.Println("Wrong connection message. Stop listening to those weird attempts of communication", err)
		ws.Close()
		return
	}
	addConn(cmsg.Username, ws)
	curruser = cmsg.Username
	for {
		var msg net.Message
		err = websocket.JSON.Receive(ws, &msg)
		if err == io.EOF {
			rmConn(curruser, ws)
			return
		}
		if err != nil {
			log.Println(curruser, "Bad message, ignoring.", err)
			continue
		}
		switch {
		case msg.GetUserList != nil:

			gp := make(models.GamePlayers)
			for _, p := range gamePlayers {
				log.Println("iterating to ", p, msg)
				if p.State == msg.GetUserList.State {
					gp[p.Name] = p
				}
			}
			log.Println(curruser, "Sending userlist of ", msg.GetUserList.State, gp, err)
			websocket.JSON.Send(ws, net.Message{ResponseUserList: &net.ResponseUserList{gp}})
		case msg.RequestFight != nil:
			log.Println("fighting now")
			websocket.JSON.Send(ws, net.Message{ResponseFight: &net.ResponseFight{"lol"}})
		}
	}

}

func connRegistrator() {
	go func() {
		for {
			r := <-connreg
			if r.action == "+" {
				if gp, ok := gamePlayers[r.username]; ok {
					gp.State = models.StateConnected
					gp.Ws = r.ws
				} else {
					gamePlayers[r.username] = &models.GamePlayer{r.username, r.ws, 100, 100, models.StateConnected, []models.SpellEffect{}}
				}

				fmt.Println("Client", r.username, "connected.")
			} else if r.action == "-" {
				if gp, ok := gamePlayers[r.username]; ok {
					gp.State = models.StateDisconnected
					gp.Ws = nil
				} else {
					fmt.Println("weird.")
				}

				fmt.Println("Client", r.username, "disconnected.")
			}
		}
	}()
}

var gamePlayers = make(models.GamePlayers)

func main() {
	connRegistrator()
	http.Handle("/game", websocket.Handler(GameServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe:" + err.Error())
	}
}
