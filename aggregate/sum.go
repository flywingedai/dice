package aggregate

import (
	"sort"

	"github.com/flywingedai/dice/core"
)

/////////////
// SUM ALL //
/////////////

type aggregate_Sum struct{}

func (a *aggregate_Sum) Load(params map[string]interface{}) {
	core.SetParams(a, params)
}

func (a *aggregate_Sum) Aggregate(result *core.Result) {
	for _, roll := range result.Results {
		result.Total += roll.Total
		result.Values = append(result.Values, roll.Total)
	}
}

//////////////////////////
// SUM SPECIFIC INDICES //
//////////////////////////

type aggregate_SumIndex struct {
	Indices []int `json:"indices"`
}

func (a *aggregate_SumIndex) Load(params map[string]interface{}) {
	core.SetParams(a, params)
}

func (a *aggregate_SumIndex) Aggregate(result *core.Result) {

	// Sort the roll values in ascending order
	sort.SliceStable(result.Results, func(i, j int) bool {
		return result.Results[i].Total < result.Results[j].Total
	})

	// Select and sum the indices from the selected rolls
	for _, index := range a.Indices {
		if index < 0 {
			index += len(result.Results)
		}
		result.Values = append(result.Values, result.Results[index].Total)
		result.Total += result.Results[index].Total
	}

}
