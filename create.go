package dice

import (
	"github.com/flywingedai/dice/aggregate"
	"github.com/flywingedai/dice/core"
	"github.com/flywingedai/dice/roll"
)

// Initialization function for the package
func init() {
	roll.Initialize()
	aggregate.Initialize()
}

/*
Create a new *Definition with a specified number of sides. A "sides" value of 6
would be equivalent to rolling a standard 6 sided dice, where ever value 1-6 has
the same chance of being rolled.
*/
func New(sides int) *core.Definition {
	d := &core.Definition{
		Children:   []*core.Definition{},
		RollType:   core.ROLL_SIDES,
		RollParams: map[string]interface{}{"sides": sides},
	}
	return d.SetDefaultSource()
}

/*
Create a new *Definition with a specified weighting for each specified value.
A weights map of:

	{
		1: 1,
		2: 1,
		3: 1,
		4: 1,
		5: 1,
		6: 2,
	}

describes a 6 normal sided die, but a roll of 6 being twice as likely as all
the other values.
*/
func NewWeighted(weights map[int]int) *core.Definition {
	d := &core.Definition{
		Children:   []*core.Definition{},
		RollType:   core.ROLL_WEIGHTED,
		RollParams: map[string]interface{}{"weights": weights},
	}
	return d.SetDefaultSource()
}
