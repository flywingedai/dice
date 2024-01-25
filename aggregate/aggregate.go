package aggregate

import (
	"github.com/flywingedai/dice/results"
	"github.com/flywingedai/dice/utils"
)

var AggregationTypes = map[AggregationName]func() Aggregation{
	Aggregate_Sum: func() Aggregation { return &aggregation_Sum{} },
}

type AggregationName string
type Aggregation interface {
	Load(params map[string]interface{})
	AggregateChances([]*results.Chances) *results.Chances
	AggregateResult([]*results.Result) int64
}

///////////////
// SUMMATION //
///////////////

type aggregation_Sum struct{}

func (a *aggregation_Sum) Load(params map[string]interface{}) {
	utils.SetParams(a, params)
}

func (a *aggregation_Sum) AggregateChances(chancesList []*results.Chances) *results.Chances {
	return results.CombineChances_Hash(
		chancesList,
		func(oldHash *int64, newValue int64) int64 {
			if oldHash == nil {
				return newValue
			}
			return *oldHash + newValue
		}, func(hashValue int64) int64 {
			return hashValue
		},
	)
}

func (a *aggregation_Sum) AggregateResult(results []*results.Result) int64 {
	total := int64(0)
	for _, r := range results {
		total += r.Result
	}
	return total
}

const Aggregate_Sum = AggregationName("sum")
