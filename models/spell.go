package models

type Spell struct {
	name      string // Spell name
	spellType string // Spell type (DirectDamage, DamageOverTime, Mesmerize...)
	value     int    // the effectiveness of the spell
	casttime  int    // cast time of the spell in ms

	duration       int          // duration of the effect in ms
	timer          int          // duration between each pulse, let 0 if no pulse
	ingredientList []Ingredient // List of spell ingredients
}

type Ingredient struct {
	name           string // Ingredient name
	keyCombination []rune // The keys the client have to fire during the ingredient set
}

// Initializing lists that will be use to store spells

var i_batWing = Ingredient{"bat wing", []rune{'h', 'j', 'k'}}
var i_bearClaw = Ingredient{"bear claw", []rune{'g', 'h', 'j'}}

var s_magicMissile = Spell{name: "Magic missile",
	spellType:      "DirectDamage",
	value:          12,
	casttime:       2000,
	duration:       0,
	timer:          0,
	ingredientList: []Ingredient{i_batWing, i_batWing, i_bearClaw}}

var s_cyanide = Spell{name: "Cyanide",
	spellType:      "DamageOverTime",
	value:          3,
	casttime:       2500,
	duration:       10000,
	timer:          1200,
	ingredientList: []Ingredient{i_bearClaw, i_batWing}}

var spellList = []Spell{s_magicMissile, s_cyanide}
