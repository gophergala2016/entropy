package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestStartSpell(t *testing.T) {
	Convey("Given player Tyriada with 100hp and magic missile spell placed as an action", t, func() {
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
		Convey("The spell should function", func() {

			go a1.StartSpell()
			WaitAndSee()
		})
	})

}
