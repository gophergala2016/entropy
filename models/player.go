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
	EffectList []*SpellEffect
}

func deleteFromSlice(i int, slice []*SpellEffect) []*SpellEffect {
	switch i {
	case -1: //do nothing
	case 0:
		slice = slice[1:]
	case len(slice) - 1:
		slice = slice[:len(slice)-1]
	default:
		slice = append(slice[:i], slice[i+1:]...)
	}
	return slice
}

func (player *GamePlayer) RemoveEffect(effectID int) bool {

	var chosen_id int = -1

	for i, effect := range player.EffectList {
		if effect.ID == effectID {
			chosen_id = i
		}
	}

	if chosen_id > -1 {
		player.EffectList = deleteFromSlice(chosen_id, player.EffectList)
		return true
	} else {
		return false
	}
}

func (player *GamePlayer) IsPoisoned() bool {

	for _, effect := range player.EffectList {
		if effect.action.spell.name == "Cyanide" {
			return true
		}
	}

	return false

}

func (player *GamePlayer) RemovePoison() bool {

	for _, effect := range player.EffectList {
		if effect.action.spell.name == "Cyanide" {
			effect.forceStop = true
			player.RemoveEffect(effect.ID)
			return true
		}
	}

	return false

}
func (player *GamePlayer) IsAsleep() bool {

	for _, effect := range player.EffectList {
		if effect.action.spell.name == "Sleep" {
			return true
		}
	}

	return false

}

func (player *GamePlayer) AddEffect(se *SpellEffect) {

	player.EffectList = append(player.EffectList, se)
}

type GamePlayers map[string]*GamePlayer
