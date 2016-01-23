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

func (a Action) IngredientSucceeded(index int) bool {
	if len(a.ingredientCheck) < len(a.spell.ingredientList[index].keyCombination) {
		return false
	}

	for _, v := range a.ingredientCheck {
		if !v {
			return false
		}
	}
	return true

}

func (a Action) SpellSucceeded() bool {
	if len(a.spellCheck) < len(a.spell.ingredientList) {
		return false
	}

	for _, v := range a.spellCheck {
		if !v {
			return false
		}
	}

	return true
}

func (a Action) startMessage() {
	fmt.Println(a.player.name + " spelling " + a.spell.name + " to date : " + a.initialTime.Format(time.UnixDate))
}

func (a Action) spellSucceededMessage() {
	fmt.Println(a.player.name + " lance correctement le sort" + a.spell.name)
}

func (a Action) endOfCastTime() time.Time {
	return a.initialTime.Add(time.Duration(a.spell.casttime) * time.Millisecond)
}

func (a Action) remainingToEndDuration() time.Duration {
	return a.endOfCastTime().Sub(time.Now())
}

func (a Action) castTimeFinished() bool {
	return a.endOfCastTime().Before(time.Now())
}

func (a *Action) StartSpell() {
	a.initialTime = time.Now()
	a.startMessage()

	for {
		if a.castTimeFinished() {

			fmt.Println("End of cast !")
			break
		} else {

		}
	}

}

type Spell struct {
	name      string // Spell name
	spellType string // Spell type (DirectDamage, DamageOverTime, Mesmerize...)
	value     int    // the effectiveness of the spell
	casttime  int    // cast time of the spell in ms

	duration       int          // duration of the effect in ms
	ingredientList []Ingredient // List of spell ingredients
}

type Ingredient struct {
	name           string // Ingredient name
	keyCombination []rune // The keys the client have to fire during the ingredient set
}

// Initializing lists that will be use to store spells and actions

var spellList = []Spell{}

func WaitAndSee() {
	for {
		// We should look to key use here
	}
}

func main() {

	p1 := GamePlayer{"Tyriada", nil, 100}

	i_batWing := Ingredient{"bat wing", []rune{'h', 'j', 'k'}}
	i_bearClaw := Ingredient{"bear claw", []rune{'g', 'h', 'j'}}

	s_magicMissile := Spell{"Magic missile",
		"DirectDamage",
		12,
		5000,
		0,
		[]Ingredient{i_batWing, i_batWing, i_bearClaw}}

	a1 := Action{&p1, s_magicMissile, []bool{}, []bool{}, time.Now()}

	go a1.StartSpell()

	WaitAndSee()
}
