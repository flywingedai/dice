package roll

import (
	"math/rand"
	"sort"

	"github.com/flywingedai/dice/core"
)

type roll_Weighted struct {
	Weights map[int]int `json:"weights"`

	values []int `json:"-"`
	total  int   `json:"-"`
}

func (r *roll_Weighted) Load(params map[string]interface{}) {
	core.SetParams(r, params)

	r.total = 0
	r.values = []int{}
	for value, weight := range r.Weights {
		r.total += weight
		r.values = append(r.values, value)
	}
	sort.SliceStable(r.values, func(i, j int) bool { return r.values[i] < r.values[j] })

}

func (r *roll_Weighted) Roll(source *rand.Rand, _ []*core.Definition) *core.Result {

	randomChoice := source.Intn(r.total)

	selectedValue := 0
	cumulativeWeight := 0

	for _, value := range r.values {
		weight := r.Weights[value]
		cumulativeWeight += weight
		if cumulativeWeight > randomChoice {
			selectedValue = value
			break
		}
	}

	return &core.Result{
		Base:   true,
		Values: []int{selectedValue},
		Total:  selectedValue,
	}
}
