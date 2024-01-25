package dice

import (
	"github.com/flywingedai/dice/aggregate"
	"github.com/flywingedai/dice/roll"
)

func (d *Definition) Multiple(n int) *roll.Definition {
	newDefinition := &roll.Definition{
		Children: []*roll.Definition{d.Definition},
		Rand:     d.Rand,
	}

	newDefinition.RollType = roll.Roll_Multiple
	newDefinition.RollParams = map[string]interface{}{
		"count": n,
	}
	newDefinition.AggName = aggregate.Aggregate_Sum
	newDefinition.AggParams = map[string]interface{}{}

	newDefinition.Load()

	return newDefinition
}
