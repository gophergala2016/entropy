package models

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStartSpell(t *testing.T) {
	Convey("Given player Tyriada with 100hp and magic missile spell placed as an action", t, func() {
		p1 := GamePlayer{"Tyriada", nil, 100, 100, StateConnected, []*SpellEffect{}}
		p2 := GamePlayer{"Maeltor", nil, 100, 100, StateConnected, []*SpellEffect{}}

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
						fmt.Println("End of effect detected")
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

		Convey("Test about damageOverTime and effectlist", func() {

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

					if ev.eventType == SPELL_EFFECT_PULSE {
						So(len(p2.EffectList), ShouldEqual, 1)
					}
					if ev.eventType == SPELL_EFFECT_END {
						So(len(p2.EffectList), ShouldEqual, 0)
						endTest = true
					}
				}

				if endTest {
					break
				}
			}

		})

		Convey("Test about damageOverTime curePoison", func() {

			a1 := Action{&p1, &p2, s_cyanide, []bool{}, []bool{}, time.Now(), make(chan rune)}
			a2 := Action{&p2, &p2, s_curePoison, []bool{}, []bool{}, time.Now(), make(chan rune)}
			go func() {

				go a1.StartSpell()

				time.Sleep(1 * time.Second)

				a1.keyChan <- 'g'
				a1.keyChan <- 'h'
				a1.keyChan <- 'j'
				a1.keyChan <- 'h'
				a1.keyChan <- 'j'
				a1.keyChan <- 'k'

				time.Sleep(1 * time.Second)

				go a2.StartSpell()

				time.Sleep(1 * time.Second)

				a2.keyChan <- 'g'
				a2.keyChan <- 'h'
				a2.keyChan <- 'j'
				a2.keyChan <- 'k'
				a2.keyChan <- 'k'
				a2.keyChan <- 'j'
				a2.keyChan <- 'h'

				time.Sleep(1 * time.Second)

			}()

			var endTest = false

			for {
				time.Sleep(50 * time.Millisecond)
				for len(EffectChannel_Test) > 0 {
					var ev = <-EffectChannel_Test

					if ev.eventType == SPELL_EFFECT_PULSE {
						fmt.Println("Effect list :", ev.spellEffect.action.target.EffectList)
					}
					if ev.eventType == SPELL_EFFECT_END {
						fmt.Println("End of effect", ev.spellEffect.action.spell.name, "detected")

						if ev.spellEffect.action.spell.name == "Cyanide" {
							endTest = true
						}
					}
				}

				if endTest {
					break
				}
			}
			fmt.Println(p2.EffectList)
			So(len(p2.EffectList), ShouldEqual, 0)
		})

		Convey("Test about sleep", func() {

			a1 := Action{&p1, &p2, s_sleep, []bool{}, []bool{}, time.Now(), make(chan rune)}
			go func() {

				go a1.StartSpell()

				time.Sleep(1 * time.Second)

				a1.keyChan <- 'k'
				a1.keyChan <- 'k'
				a1.keyChan <- 'j'
				a1.keyChan <- 'h'

				a1.keyChan <- 'k'
				a1.keyChan <- 'k'
				a1.keyChan <- 'j'
				a1.keyChan <- 'h'

				a1.keyChan <- 'h'
				a1.keyChan <- 'j'
				a1.keyChan <- 'k'

			}()

			var endTest = false

			for {
				time.Sleep(50 * time.Millisecond)
				for len(EffectChannel_Test) > 0 {
					var ev = <-EffectChannel_Test

					if ev.eventType == SPELL_EFFECT_PULSE {
						So(p2.IsAsleep(), ShouldEqual, true)
					}
					if ev.eventType == SPELL_EFFECT_END {
						So(p2.IsAsleep(), ShouldEqual, false)
						endTest = true
					}
				}

				if endTest {
					break
				}
			}

		})

		Convey("Test about heal", func() {

			a1 := Action{&p1, &p1, s_layingonofhands, []bool{}, []bool{}, time.Now(), make(chan rune)}
			go func() {

				p1.Hp = 40

				go a1.StartSpell()

				time.Sleep(1 * time.Second)

				a1.keyChan <- 'k'
				a1.keyChan <- 'k'
				a1.keyChan <- 'j'
				a1.keyChan <- 'h'

				a1.keyChan <- 'g'
				a1.keyChan <- 'h'
				a1.keyChan <- 'j'

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
						So(p1.Hp, ShouldEqual, 65)
						endTest = true
					}
				}

				if endTest {
					break
				}
			}

		})

	})

}
