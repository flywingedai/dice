package roll

import (
	"math/rand"

	"github.com/flywingedai/dice/core"
)

type roll_Skip struct{}

func (r *roll_Skip) Load(params map[string]interface{}) {
	core.SetParams(r, params)
}

func (r *roll_Skip) Roll(source *rand.Rand, definitions []*core.Definition) *core.Result {
	result := &core.Result{
		Base:    false,
		Results: []*core.Result{},
		Values:  []int{},
		Total:   0,
	}

	// Just perform each of the definitions individually
	for _, definition := range definitions {
		definitionResult := definition.Roll()
		result.Results = append(result.Results, definitionResult)
	}

	return result
}
