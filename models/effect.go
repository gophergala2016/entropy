package models

import (
	"fmt"
	"time"
)

type SpellEffect struct {
	action    Action    // the spell that fired this effect
	startTime time.Time // when the effect has started
}

func (se SpellEffect) Start() {

	go func() {

		if se.action.caster == nil {
			se.End()
			return
		}

		se.OnDirectEffect()

		if se.action.spell.duration > 0 {
			for {
				if !se.OnPulse() {
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
	SendOnEffectChannels(SpellEffectEvent{SPELL_EFFECT_END, se})
}

func (se SpellEffect) OnDirectEffect() bool {
	if se.action.target.Hp == 0 {
		return false
	}

	se.action.target.Hp -= se.action.spell.value
	fmt.Println(se.action.caster.Name+" hits "+se.action.target.Name+", and does ", se.action.spell.value, " damage points ! (", se.action.target.Hp, " / ", se.action.target.MaxHp, " )")

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
