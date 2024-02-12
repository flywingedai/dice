package roll

import (
	"math/rand"

	"github.com/flywingedai/dice/core"
)

type roll_Multiple struct {
	Count int `json:"count"`
}

func (r *roll_Multiple) Load(params map[string]interface{}) {
	core.SetParams(r, params)
}

func (r *roll_Multiple) Roll(source *rand.Rand, definitions []*core.Definition) *core.Result {
	result := &core.Result{
		Base:    false,
		Results: []*core.Result{},
		Values:  []int{},
		Total:   0,
	}

	for i := 0; i < r.Count; i++ {
		result.Results = append(result.Results, definitions[0].Roll())
	}

	return result
}
