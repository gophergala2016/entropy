package main

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

/* The gameplayer is linked to a client connected to the server.
 * Maybe we could set a name at client start...
 */
type GamePlayer struct {
	name string          // The name of the player
	ws   *websocket.Conn // The connection informations
	hp   int             // Player health points
}

/* An action is starting by the client.
 * It contains all player and spell informations, and the current stage of this spell launch.
 */
type Action struct {
	player *GamePlayer // The player that initiated this action
	spell  Spell       // Spell chosen when the action started

	//spellStep      int // The index of the current ingredient used
	//ingredientStep int // The index of the next expected key

	spellCheck      []bool // Saved results of spell checks
	ingredientCheck []bool // Saved resultas of ingredients checks

	initialTime time.Time // Time at the start of this spell
}

type Spell struct {
	name           string       // Spell name
	spellType      string       // Spell type (DirectDamage, DamageOverTime, Mesmerize...)
	value          int          // the effectiveness of the spell
	duration       int          // duration of the effect in ms
	ingredientList []Ingredient // List of spell ingredients
}

type Ingredient struct {
	name           string // Ingredient name
	keyCombination []rune // The keys the client have to fire during the ingredient set
}

// Initializing lists that will be use to store spells and actions
var actions = []Action{}
var spellList = []Spell{}

func CheckActionList() {
	//go func() {
	for {

		if len(actions) == 0 {
			continue
		}
		fmt.Println("Il y a " + (string)(len(actions)) + " actions.")
		for a := range actions {
			fmt.Println(a.player.name + " lance le sort" + a.spell.name + " à la date : " + a.initialTime.Format(time.UnixDate))
		}

	}
	//}()
}

func main() {

	p1 := GamePlayer{"Tyriada", nil, 100}

	i_batWing := Ingredient{"bat wing", []rune{'h', 'j', 'k'}}
	i_bearClaw := Ingredient{"bear claw", []rune{'g', 'h', 'j'}}

	s_magicMissile := Spell{"Magic missile",
		"DirectDamage",
		12,
		0,
		[]Ingredient{i_batWing, i_batWing, i_bearClaw}}

	a1 := Action{&p1, s_magicMissile, []bool{}, []bool{}, time.Now()}
	fmt.Println(a1.spell.name)
	CheckActionList()
}
