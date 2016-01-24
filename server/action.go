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
	caster *models.GamePlayer // The player that initiated this action
	target *models.GamePlayer // The target of this spell
	spell  Spell              // Spell chosen when the action started

	spellCheck      []bool // Saved results of spell checks
	ingredientCheck []bool // Saved resultas of ingredients checks

	initialTime time.Time // Time at the start of this spell

	keyChan chan rune // The key channel that manage the action

}

func (a Action) HasGoodKeyCombination(index int) bool {
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

func (a Action) HasGoodIngredients() bool {
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
	fmt.Println(a.caster.Name + " spelling " + a.spell.name + " to date : " + a.initialTime.Format(time.UnixDate))
}

func (a Action) goodKey(ingredientPos int, keyPos int) {
	fmt.Println("  + " + a.caster.Name + " used good key for ingredient " + a.spell.ingredientList[ingredientPos].name)
}

func (a Action) badKey(ingredientPos int, keyPos int) {
	fmt.Println("  + " + a.caster.Name + " used bad key for ingredient " + a.spell.ingredientList[ingredientPos].name)
}

func (a Action) goodIngredient(ingredientPos int) {
	fmt.Println("    + " + a.caster.Name + " succeeded adding " + a.spell.ingredientList[ingredientPos].name)
}

func (a Action) badIngredient(ingredientPos int) {
	fmt.Println("    + " + a.caster.Name + " missed adding " + a.spell.ingredientList[ingredientPos].name)
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

	a.keyChan = make(chan rune, 100)

	// if stepFinished is true,  success have been set
	var stepFinished = false
	// the result of the key combination
	var success = false

	for {
		time.Sleep(50 * time.Millisecond)

		if a.castTimeFinished() {

			fmt.Println("End of cast !")

			if !stepFinished {
				// the player didn't send enough key, he missed !
				success = false
			}

			a.spell.StartEffect(a.caster, a.target, success)

			break
		} else {

			if !stepFinished {
				for len(a.keyChan) > 0 {
					r := <-a.keyChan

					var ingredientPos = len(a.spellCheck)
					var keyPos = len(a.ingredientCheck)
					var keySuccess = (a.spell.ingredientList[ingredientPos].keyCombination[keyPos] == r)

					if keySuccess {
						a.goodKey(ingredientPos, keyPos)
					} else {
						a.badKey(ingredientPos, keyPos)
					}

					a.ingredientCheck = append(a.ingredientCheck, keySuccess)

					if len(a.ingredientCheck) == len(a.spell.ingredientList[ingredientPos].keyCombination) {
						var ingredientSuccess = a.HasGoodKeyCombination(ingredientPos)

						if ingredientSuccess {
							a.goodIngredient(ingredientPos)
						} else {
							a.badIngredient(ingredientPos)
						}

						a.spellCheck = append(a.spellCheck, ingredientSuccess)
						a.ingredientCheck = []bool{}

					}

					if len(a.spellCheck) == len(a.spell.ingredientList) {
						stepFinished = true
						success = a.HasGoodIngredients()
					}

				}
			}
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

func (s Spell) StartEffect(caster *models.GamePlayer, target *models.GamePlayer, success bool) {

	if success {
		fmt.Println(caster.Name + " succeeds launching " + s.name)
	} else {
		fmt.Println(caster.Name + " missed launching " + s.name)
	}
	if success {
		switch s.spellType {
		case "DirectDamage":
			target.Hp -= s.value
			fmt.Println(caster.Name+" hits "+target.Name+", and does ", s.value, " damage points ! (", target.Hp, " / ", target.MaxHp, " )")
		default:
		}
	} else {
		// TODO : random spell startEffect !
	}

}

type Ingredient struct {
	name           string // Ingredient name
	keyCombination []rune // The keys the client have to fire during the ingredient set
}

// Initializing lists that will be use to store spells

var i_batWing = Ingredient{"bat wing", []rune{'h', 'j', 'k'}}
var i_bearClaw = Ingredient{"bear claw", []rune{'g', 'h', 'j'}}

var s_magicMissile = Spell{"Magic missile",
	"DirectDamage",
	12,
	5000,
	0,
	[]Ingredient{i_batWing, i_batWing, i_bearClaw}}

var spellList = []Spell{s_magicMissile}
