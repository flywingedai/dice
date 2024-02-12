package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
)

/*
Lock mutex to make sure nothing causes a race condition within the package
(Unless a user passes in a function that would cause such a race condition.)
*/
var lock = sync.Mutex{}

///////////
// ROLLS //
//////////

type Roll interface {
	Load(params map[string]interface{})
	Roll(*rand.Rand, []*Definition) *Result
}

var rollTypes = map[string]func() Roll{}

func AddRollType(rollType string, newFunction func() Roll) {
	lock.Lock()
	defer lock.Unlock()
	rollTypes[rollType] = newFunction
}

func GetRollType(rollType string) Roll {
	lock.Lock()
	defer lock.Unlock()
	rollFunction, ok := rollTypes[rollType]
	if !ok {
		panic(fmt.Sprintf("Invalid rollType %s provided to GetRollType", rollType))
	}
	return rollFunction()
}

//////////////////
// AGGREGATIONS //
//////////////////

type Aggregation interface {
	Load(params map[string]interface{})
	Aggregate(*Result)
}

var aggregationTypes = map[string]func() Aggregation{}

func AddAggregationType(aggregationType string, newFunction func() Aggregation) {
	lock.Lock()
	defer lock.Unlock()
	aggregationTypes[aggregationType] = newFunction
}

func GetAggregationType(aggregationType string) Aggregation {
	lock.Lock()
	defer lock.Unlock()
	aggregationFunction, ok := aggregationTypes[aggregationType]
	if !ok {
		panic(fmt.Sprintf("Invalid aggregationType %s provided to GetAggregationType", aggregationType))
	}
	return aggregationFunction()
}

/*
Helper function to be used in Roll and Aggregation interface structs.
Spcifically it makes the Load(map[string]interface{}) method pretty trivial
and allows all load functions to look like this:

	func (s *interfaceStruct) Load(params map[string]interface{}) {
		core.SetParams(s, params)
	}
*/
func SetParams[T any](c *T, params map[string]interface{}) {
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(paramsBytes, c)
	if err != nil {
		panic(err)
	}
}
