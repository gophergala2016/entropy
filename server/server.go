package main

import (
	"fmt"
	"github.com/gophergala2016/entropy/models"
	//	"github.com/gophergala2016/entropy/net"
	"github.com/gorilla/websocket"
	//	"io"
	"log"
	"net/http"
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

var upgrader = websocket.Upgrader{}

func GameServerHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println("msg: ", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

/*
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
		var msg net.GetUserList //interface{}
		err = websocket.JSON.Receive(ws, &msg)
		if err == io.EOF {
			rmConn(curruser, ws)
			return
		}
		if err != nil {
			log.Println(curruser, "Bad message, ignoring.", err)
			continue
		}

		//switch m := msg.(type) {
		//case net.GetUserList:
		gp := make(models.GamePlayers)
		for _, p := range gamePlayers {
			if p.State == msg.State {
				gp[p.Name] = p
			}
		}

		websocket.JSON.Send(ws, gp)
		//default:
		//	log.Println("Unknown message:", m)
		//}

	}

}

*/

func connRegistrator() {
	go func() {
		for {
			r := <-connreg
			if r.action == "+" {
				if gp, ok := gamePlayers[r.username]; ok {
					gp.State = models.StateConnected
					gp.Ws = r.ws
				} else {
					gamePlayers[r.username] = &models.GamePlayer{r.username, r.ws, 100, models.StateConnected}
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
	http.HandleFunc("/game", GameServerHandler)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe:" + err.Error())
	}
}
