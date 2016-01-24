package main

import (
	"fmt"
	"time"

	"github.com/gophergala2016/entropy/models"
)

/* An action is starting by the client.
 * It contains all player and spell informations, and the current stage of this spell launch.
 */
type Action struct {
	player *models.GamePlayer // The player that initiated this action
	spell  Spell              // Spell chosen when the action started

	//spellStep      int // The index of the current ingredient used
	//ingredientStep int // The index of the next expected key

	spellCheck      []bool // Saved results of spell checks
	ingredientCheck []bool // Saved resultas of ingredients checks

	initialTime time.Time // Time at the start of this spell

	keyEvenChan chan keyboardEvent // The key event channel that manage the action

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
		time.Sleep(50 * time.Millisecond)

		fmt.Println("test")
		if !v {
			return false
		}
	}

	return true
}

func (a Action) startMessage() {
	fmt.Println(a.player.Name + " spelling " + a.spell.name + " to date : " + a.initialTime.Format(time.UnixDate))
}

func (a Action) spellSucceededMessage() {
	fmt.Println(a.player.Name + " lance correctement le sort" + a.spell.name)
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

// Initializing lists that will be use to store spells

var spellList = []Spell{}
