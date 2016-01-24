package net

import (
	"github.com/gophergala2016/entropy/models"
)

type Message struct {
	Connection                  *Connection
	Disconnection               *Disconnection
	GetUserList                 *GetUserList
	ResponseUserList            *ResponseUserList
	RequestFight                *RequestFight
	ResponseFight               *ResponseFight
	InArena                     *InArena
	SelectFightSpells           *SelectFightSpells
	StartCasting                *StartCasting
	ResponseSpellCharacteristic *ResponseSpellCharacteristic
	CastingKey                  *CastingKey
	ResponseCastingKey          *ResponseCastingKey
	ResponseSpellResult         *ResponseSpellResult
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
	Opponent models.GamePlayer
}

type InArena struct {
}

type SelectFightSpells struct {
	Spells []models.Spell
}

type StartCasting struct {
	caster models.GamePlayer
}

type ResponseSpellCharacteristic struct {
	Spell models.Spell
}

type CastingKey struct {
	Key rune
}

type ResponseCastingKey struct {
	Key  rune
	Step int
}

type ResponseSpellResult struct {
	Caster models.GamePlayer
	Target models.GamePlayer
}
