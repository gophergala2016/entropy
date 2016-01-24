package models

import (
	"golang.org/x/net/websocket"
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
	Name       string          // The name of the player
	Ws         *websocket.Conn // The connection informations
	Hp         int             // Player health points
	MaxHp      int             // Player max Health points
	State      GamePlayerState
	EffectList []SpellEffect
}

type GamePlayers map[string]*GamePlayer
