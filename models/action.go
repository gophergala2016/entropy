package models

import (
	"fmt"
	"time"
)

/* An action is starting by the client.
 * It contains all player and spell informations, and the current stage of this spell launch.
 */
type Action struct {
	caster *GamePlayer // The player that initiated this action
	target *GamePlayer // The target of this spell
	spell  Spell       // Spell chosen when the action started

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

	a.keyChan = make(chan rune, 2)

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

			a.LaunchEffect(success)

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

func (a *Action) LaunchEffect(success bool) {

	if success {
		fmt.Println(a.caster.Name + " succeeds launching " + a.spell.name)
	} else {
		fmt.Println(a.caster.Name + " missed launching " + a.spell.name)
	}

	var effect SpellEffect

	if success {
		effect = SpellEffect{*a, time.Now()}

	} else {
		effect = SpellEffect{Action{}, time.Now()}
	}
	effect.Start()

}
