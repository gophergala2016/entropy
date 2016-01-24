package main

import (
	"bufio"
	"fmt"
	"github.com/gophergala2016/entropy/models"
	"github.com/gophergala2016/entropy/net"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"strings"
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

var username string

func main() {
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
			displayUserList(&msg.ResponseUserList.GamePlayers)
		}
		time.Sleep(1 * time.Second)
	}
}

func displayUserList(gp *models.GamePlayers) {
	fmt.Println("Connected Users:")
	i := 0
	userlist := make([]string, 0, len(*gp))
	for k, _ := range *gp {
		if k == username {
			continue
		}
		fmt.Printf("%1d) %s\n", i, k)
		userlist = append(userlist, k)
		i++
	}
	correctInput := false
	for !correctInput {
		fmt.Println("\nSelect your opponent (number ID): ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		userid, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			fmt.Println("Not a correct number")
		} else {
			if len(userlist) <= userid {
				fmt.Println("Bad number, or list empty")
				correctInput = true
				continue
			}
			fmt.Println("Chosen", userlist[userid], "from", userlist)
			correctInput = true
		}

	}

}

