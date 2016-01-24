package models

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStartSpell(t *testing.T) {
	Convey("Given player Tyriada with 100hp and magic missile spell placed as an action", t, func() {
		p1 := GamePlayer{"Tyriada", nil, 100, 100, StateConnected, []SpellEffect{}}
		p2 := GamePlayer{"Maeltor", nil, 100, 100, StateConnected, []SpellEffect{}}

		Convey("Test about right choices", func() {

			a1 := Action{&p1, &p2, s_magicMissile, []bool{}, []bool{}, time.Now(), make(chan rune)}
			go func() {

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
			}()

			var endTest = false

			for {
				time.Sleep(50 * time.Millisecond)
				for len(EffectChannel_Test) > 0 {
					var ev = <-EffectChannel_Test
					if ev.eventType == SPELL_EFFECT_END {
						fmt.Println("End of cast detected")
						endTest = true
					}
				}

				if endTest {
					break
				}
			}

			So(a1.HasGoodIngredients(), ShouldEqual, true)
		})

		Convey("Test about wrong choices", func() {

			a1 := Action{&p1, &p2, s_magicMissile, []bool{}, []bool{}, time.Now(), make(chan rune)}
			go func() {

				go a1.StartSpell()

				time.Sleep(1 * time.Second)

				a1.keyChan <- 'a'
				a1.keyChan <- 'j'
				a1.keyChan <- 'k'
				a1.keyChan <- 'h'
				a1.keyChan <- 'j'
				a1.keyChan <- 'k'
				a1.keyChan <- 'g'
				a1.keyChan <- 'h'
				a1.keyChan <- 'j'
			}()

			var endTest = false

			for {
				time.Sleep(50 * time.Millisecond)
				for len(EffectChannel_Test) > 0 {
					var ev = <-EffectChannel_Test
					if ev.eventType == SPELL_EFFECT_END {
						fmt.Println("End of cast detected")
						endTest = true
					}
				}

				if endTest {
					break
				}
			}

			So(a1.HasGoodIngredients(), ShouldEqual, false)
		})

		Convey("Test about damageOverTime", func() {

			a1 := Action{&p1, &p2, s_cyanide, []bool{}, []bool{}, time.Now(), make(chan rune)}
			go func() {

				go a1.StartSpell()

				time.Sleep(1 * time.Second)

				a1.keyChan <- 'g'
				a1.keyChan <- 'h'
				a1.keyChan <- 'j'
				a1.keyChan <- 'h'
				a1.keyChan <- 'j'
				a1.keyChan <- 'k'
			}()

			var endTest = false

			for {
				time.Sleep(50 * time.Millisecond)
				for len(EffectChannel_Test) > 0 {
					var ev = <-EffectChannel_Test
					if ev.eventType == SPELL_EFFECT_END {
						fmt.Println("End of cast detected")
						endTest = true
					}
				}

				if endTest {
					break
				}
			}

			So(a1.HasGoodIngredients(), ShouldEqual, true)
		})

	})

}
