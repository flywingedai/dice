package roll

import (
	"math/rand"

	"github.com/flywingedai/dice/core"
)

type roll_Sides struct {
	Sides int `json:"sides"`
}

func (r *roll_Sides) Load(params map[string]interface{}) {
	core.SetParams(r, params)
}

func (r *roll_Sides) Roll(source *rand.Rand, _ []*core.Definition) *core.Result {
	value := source.Intn(r.Sides) + 1
	return &core.Result{
		Base:   true,
		Values: []int{value},
		Total:  value,
	}
}
