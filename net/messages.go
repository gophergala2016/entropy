package net

import (
	"github.com/gophergala2016/entropy/models"
)

type Message struct {
	Connection       *Connection
	Disconnection    *Disconnection
	GetUserList      *GetUserList
	ResponseUserList *ResponseUserList
	RequestFight     *RequestFight
	ResponseFight    *ResponseFight
}

type Connection struct {
	Username string
}

type Disconnection struct {
	Username string
}

type GetUserList struct {
	State models.GamePlayerState
}

type ResponseUserList struct {
	GamePlayers models.GamePlayers
}

type RequestFight struct {
	Opponentname string
}
type ResponseFight struct {
	Opponentname string
}
