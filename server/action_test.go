package main

import (
	"testing"
	"time"

	. "github.com/gophergala2016/entropy/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStartSpell(t *testing.T) {
	Convey("Given player Tyriada with 100hp and magic missile spell placed as an action", t, func() {
		p1 := GamePlayer{"Tyriada", nil, 100, StateConnected}
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

			a1.StartSpell()
		})
	})

}
