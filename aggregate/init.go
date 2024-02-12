package aggregate

import "github.com/flywingedai/dice/core"

func Initialize() {
	core.AddAggregationType(core.AGGREGATE_SUM, func() core.Aggregation { return &aggregate_Sum{} })
	core.AddAggregationType(core.AGGREGATE_SUM_INDEX, func() core.Aggregation { return &aggregate_SumIndex{} })
}
