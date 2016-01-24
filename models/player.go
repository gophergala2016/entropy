package models

import (
	"github.com/gorilla/websocket"
)

type GamePlayerState int

const (
	StateDisconnected GamePlayerState = iota
	StateConnected
	StateFighting
	StateUnavailable
)

/* The GamePlayer is linked to a client connected to the server.
 * Maybe we could set a name at client start...
 */
type GamePlayer struct {
	Name  string          // The name of the player
	Ws    *websocket.Conn // The connection informations
	Hp    int             // Player health points
	State GamePlayerState
}

type GamePlayers map[string]*GamePlayer
