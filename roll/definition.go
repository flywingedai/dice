package roll

import (
	"math/big"
	"math/rand"
	"sort"

	"github.com/flywingedai/dice/aggregate"
	"github.com/flywingedai/dice/results"
)

/*
A Definition describes how a roll takes place.
*/
type Definition struct {

	// Random seed to use for rolls
	Rand *rand.Rand

	// Roll interface for managing the action of actually rolling dice
	roll Roll

	// Aggregation interface for managing the action of aggregating whatever was rolled
	aggregation aggregate.Aggregation

	// Whether or not the chances have been set for this definition
	Set bool

	// The chances object describes the chance of each roll value happening
	Chances *results.Chances

	// Whether or not the definition has properly loaded in the .roll and
	// .aggregation attributes.
	Loaded bool

	// Fields that are necessary for loading/customization
	Children []*Definition `json:"children"`

	RollType RollType                  `json:"rollType"`
	AggName  aggregate.AggregationName `json:"aggName"`

	RollParams map[string]interface{} `json:"rollParams"`
	AggParams  map[string]interface{} `json:"aggParams"`
}

/*
Load in all the json fields to a proper roll and aggregation interfaces. Must
be called in order to set the dice rolls for the system.
*/
func (d *Definition) Load() *Definition {

	if d.Loaded || d.RollType == Roll_Base {
		return d
	}
	d.Loaded = true

	// Load in the aggregation and the roll
	d.roll = RollTypes[d.RollType]()
	d.roll.Load(d.RollParams)

	d.aggregation = aggregate.AggregationTypes[d.AggName]()
	d.aggregation.Load(d.AggParams)

	return d
}

/* Loads the definition based on the parameters given */
func (d *Definition) SetChances() *Definition {
	d.Load()

	if d.Set {
		return d
	}
	d.Set = true

	// Make sure all the children params are loaded in as well.
	if d.RollType != Roll_Base {
		chances := []*results.Chances{}
		for _, child := range d.Children {
			if !child.Set {
				child.SetChances()
			}
			chances = append(chances, child.Chances)
		}

		// Compute the chances for this new roll
		chancesList := d.roll.Chances(chances)
		d.Chances = d.aggregation.AggregateChances(chancesList)
	}

	// Created the cumulative weight
	d.Chances.CumulativeWeight = map[int64]*big.Int{}
	currentWeight := big.NewInt(0)
	for _, value := range d.Chances.Values {
		currentWeight = new(big.Int).Add(currentWeight, d.Chances.Weights[value])
		d.Chances.CumulativeWeight[value] = currentWeight
	}

	return d
}

/*
Roll a Definition without generating results. If you want to accumulate results
during the run and return the full details instead, call RollResult() instead.
*/
func (d *Definition) Roll() int64 {

	/*
		Get a random big.Int using the built-in big.Int.Rand function. This will
		be used to select a value from the weights map in the Definition.
	*/
	r := new(big.Int).Rand(d.Rand, d.Chances.Total)

	/*
		Loop through the weights map and select the first value that is higher
		than the randomly generated number, r
	*/
	i := sort.Search(len(d.Chances.Values), func(i int) bool {
		return r.Cmp(d.Chances.CumulativeWeight[d.Chances.Values[i]]) < 0
	})
	return d.Chances.Values[i]

}

/*
Roll a Definition without generating results. If you want to accumulate results
furing the run and return the full details instead, call RollResult() instead.
*/
func (d *Definition) RollResult() *results.Result {

	// Create the result object
	result := &results.Result{
		Chances: d.Chances,
	}

	// Make sure all the children are rolled properly
	if d.RollType != Roll_Base {
		result.Children = d.roll.Results(d.Children)
		result.Result = d.aggregation.AggregateResult(result.Children)
	} else {
		result.Result = d.Roll()
	}

	return result

}

/*
Get the probabilities of every roll happening
*/
func (d *Definition) TrueAnalysis() map[int64]*big.Float {
	d.SetChances()

	probabilities := map[int64]*big.Float{}
	total := new(big.Float).SetInt(d.Chances.Total)

	for value, weight := range d.Chances.Weights {
		probabilities[value] = new(big.Float).Quo(new(big.Float).SetInt(weight), total)
	}

	return probabilities
}

/*
Get the probabilities of every roll happening
*/
func (d *Definition) SimAnalysis(n int64) map[int64]*big.Float {
	probabilities := map[int64]*big.Float{}
	total := big.NewFloat(float64(n))

	for i := int64(0); i < n; i++ {
		r := d.RollResult().Result
		if probabilities[r] == nil {
			probabilities[r] = big.NewFloat(0)
		}
		probabilities[r].Add(probabilities[r], big.NewFloat(1))
	}

	for _, value := range d.Chances.Values {
		if probabilities[value] == nil {
			probabilities[value] = big.NewFloat(0)
		} else {
			probabilities[value] = new(big.Float).Quo(probabilities[value], total)
		}
	}

	return probabilities
}
