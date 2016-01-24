package models

import (
	"fmt"
	"math/rand"
	"time"
)

type SpellEffect struct {
	action    Action    // the spell that fired this effect
	startTime time.Time // when the effect has started
	ID        int
	forceStop bool
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func (se *SpellEffect) Start() {

	se.ID = randInt(0, 99999)

	if se.action.target != nil {
		se.action.target.AddEffect(se)
	}

	go func() {

		if se.action.caster == nil || se.forceStop {
			se.End()
			return
		}

		se.OnDirectEffect()

		if se.action.spell.duration > 0 {
			for {
				if !se.OnPulse() || se.forceStop {
					se.End()
					return
				}
				time.Sleep(time.Duration(se.action.spell.timer) * time.Millisecond)
			}
		}

		se.End()
	}()

}

func (se SpellEffect) End() {

	if se.action.target != nil {
		se.action.target.RemoveEffect(se.ID)
	}

	switch se.action.spell.spellType {

	case "DamageOverTime":
		fmt.Println(se.action.target.Name + " is not poisoned yet.")
	case "Sleep":
		fmt.Println(se.action.target.Name + " is not asleep yet.")
	}

	SendOnEffectChannels(SpellEffectEvent{SPELL_EFFECT_END, se})
}

func (se SpellEffect) OnDirectEffect() bool {
	if se.action.target.Hp == 0 {
		return false
	}

	switch se.action.spell.spellType {

	case "Heal":
		se.action.target.Hp += se.action.spell.value
		fmt.Println(se.action.caster.Name+" heals "+se.action.target.Name+", and gives ", se.action.spell.value, " health points ! (", se.action.target.Hp, " / ", se.action.target.MaxHp, " )")

	case "DirectDamage":
		se.action.target.Hp -= se.action.spell.value
		fmt.Println(se.action.caster.Name+" hits "+se.action.target.Name+", and does ", se.action.spell.value, " damage points ! (", se.action.target.Hp, " / ", se.action.target.MaxHp, " )")

	case "DamageOverTime":
		se.action.target.Hp -= se.action.spell.value
		fmt.Println(se.action.caster.Name+" hits "+se.action.target.Name+", and does ", se.action.spell.value, " damage points ! (", se.action.target.Hp, " / ", se.action.target.MaxHp, " )")
	case "CurePoison":
		if se.action.target.RemovePoison() {
			fmt.Println(se.action.target.Name + " healed from poison.")
		} else {
			fmt.Println(se.action.target.Name + " is not poisoned.")
		}
	case "Sleep":
		fmt.Println(se.action.target.Name + " falls asleep.")

	}

	SendOnEffectChannels(SpellEffectEvent{SPELL_EFFECT_DIRECT, se})
	return true
}

func (se SpellEffect) OnPulse() bool {

	if se.action.target.Hp == 0 {
		return false
	}

	endDate := se.startTime.Add(time.Duration(se.action.spell.duration) * time.Millisecond)

	for endDate.Before(time.Now()) {
		return false
	}

	switch se.action.spell.spellType {
	case "DamageOverTime":
		se.action.target.Hp -= se.action.spell.value
		fmt.Println(se.action.target.Name+" suffers and lose ", se.action.spell.value, " health points ! (", se.action.target.Hp, " / ", se.action.target.MaxHp, " )")

	}

	SendOnEffectChannels(SpellEffectEvent{SPELL_EFFECT_PULSE, se})
	return true
}

type SpellEffectEventType int

const (
	SPELL_EFFECT_DIRECT SpellEffectEventType = iota + 1
	SPELL_EFFECT_PULSE
	SPELL_EFFECT_END
)

type SpellEffectEvent struct {
	eventType   SpellEffectEventType
	spellEffect SpellEffect
}

var EffectChannel_Test = make(chan SpellEffectEvent, 2)
var EffectChannel_Server = make(chan SpellEffectEvent, 2)

func SendOnEffectChannels(spellEffectEvent SpellEffectEvent) {
	go func() {
		EffectChannel_Test <- spellEffectEvent
	}()
	go func() {
		EffectChannel_Server <- spellEffectEvent
	}()
}
