package main

import (
	"testing"
	"time"

	. "github.com/gophergala2016/entropy/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStartSpell(t *testing.T) {
	Convey("Given player Tyriada with 100hp and magic missile spell placed as an action", t, func() {
		p1 := GamePlayer{"Tyriada", nil, 100, 100, StateConnected}
		p2 := GamePlayer{"Maeltor", nil, 100, 100, StateConnected}

		a1 := Action{&p1, &p2, s_magicMissile, []bool{}, []bool{}, time.Now(), make(chan rune)}

		Convey("Test about right choices", func() {

			go a1.StartSpell()

			time.Sleep(1 * time.Second)

			a1.keyChan <- 'h'
			a1.keyChan <- 'j'
			a1.keyChan <- 'k'
			a1.keyChan <- 'h'
			a1.keyChan <- 'j'
			a1.keyChan <- 'k'
			a1.keyChan <- 'g'
			a1.keyChan <- 'h'
			a1.keyChan <- 'j'

			time.Sleep(5 * time.Second)

			So(a1.HasGoodIngredients(), ShouldEqual, true)
		})

		Convey("Test about wrong choices", func() {

			go a1.StartSpell()

			time.Sleep(1 * time.Second)

			a1.keyChan <- 'h'
			a1.keyChan <- 'j'
			a1.keyChan <- 'k'
			a1.keyChan <- 'h'
			a1.keyChan <- 'j'
			a1.keyChan <- 'k'
			a1.keyChan <- 'g'
			a1.keyChan <- 'h'
			a1.keyChan <- 's'

			time.Sleep(5 * time.Second)

			So(a1.HasGoodIngredients(), ShouldEqual, false)
		})
	})

}
