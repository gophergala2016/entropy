package net

import (
	"github.com/gophergala2016/entropy/models"
)

type Message struct {
	Msg     string
	Version string
}

type Connection struct {
	Username string
}

type Disconnection struct {
	Username string
}

type GetUserList struct {
	State *models.GamePlayerState
}
