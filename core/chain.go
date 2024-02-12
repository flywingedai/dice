package core

/*
TODO
*/
func (d *Definition) Multiple(n int) *Definition {
	newDefinition := new(d, ROLL_MULTIPLE, AGGREGATE_SUM)
	newDefinition.RollParams["count"] = n
	return newDefinition
}

/*
TODO
*/
func (d *Definition) Advantage() *Definition {
	newDefinition := new(d, ROLL_MULTIPLE, AGGREGATE_SUM_INDEX)
	newDefinition.RollParams["count"] = 2
	newDefinition.AggregationParams["indices"] = []int{-1}
	return newDefinition
}

/*
TODO
*/
func (d *Definition) Disadvantage() *Definition {
	newDefinition := new(d, ROLL_MULTIPLE, AGGREGATE_SUM_INDEX)
	newDefinition.RollParams["count"] = 2
	newDefinition.AggregationParams["indices"] = []int{0}
	return newDefinition
}

///////////////////
// CHAIN HELPERS //
///////////////////

/*
Creates a new definition with standard fields filled out
*/
func new(d *Definition, rollType, aggregationType string) *Definition {
	return &Definition{
		RollType:        rollType,
		AggregationType: aggregationType,

		RollParams:        map[string]interface{}{},
		AggregationParams: map[string]interface{}{},

		Children: []*Definition{d},
	}
}
