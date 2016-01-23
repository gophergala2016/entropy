# Spell

* name
* description
* effects
  * success (eg. do damage)
  * failure (eg. do damage to yourself)
  * mistype (eg. makes you lose 1 sec on the next casting)
* casting
  * components
    * type of component (vocal, gesture, material component, focus)
    * key/letter combination
  * casting time

# Spellbook

* list of spells

# Casting tracker

Helps the server knows where in the casting process the player is

* Spell: the spell that is currently being casted
* current keytype (eg. already type a, b, d, z. Still needs f, r, a to cast the spell)


# User

* name
* spellbook
* current state (disconnected, free to fight, fighting)
