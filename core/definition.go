package core

import "math/rand"

/*
The Definition struct is the base object for the entirity of the "dice" package.
The Roll Definition contains many methods which are expanded upon in other
files.
*/
type Definition struct {

	// The source used for all random events that occur for this definition
	source *rand.Rand `json:"-"`

	// If the definition is a parent, it will have an array of child Definitions
	// that are necessary to facilitate the roll
	Children []*Definition `json:"children,omitempty"`

	// Interfaces for the Roll and Aggregation
	roll        Roll        `json:"-"`
	aggregation Aggregation `json:"-"`

	// Params necessary for the roll and aggregation to be saved and loaded
	RollType   string                 `json:"rollType"`
	RollParams map[string]interface{} `json:"rollParams"`

	AggregationType   string                 `json:"aggregationType"`
	AggregationParams map[string]interface{} `json:"aggregationParams"`

	// Tag which indicates whether or not the definition has been loaded or not
	// This will automatically happen if the roll and aggregation are nil
	loaded bool `json:"-"`
}

// Copy function to make a new duplicate of a roll
func (d *Definition) Copy() *Definition {
	newDefinition := &Definition{
		Children:          []*Definition{},
		RollType:          d.RollType,
		RollParams:        d.RollParams,
		AggregationType:   d.AggregationType,
		AggregationParams: d.AggregationParams,
	}

	for _, child := range d.Children {
		newDefinition.Children = append(newDefinition.Children, child.Copy())
	}

	return newDefinition
}

// Internal function for managing definition loads
func (d *Definition) tryLoad() {

	// No need to load if the definition is already loaded
	if d.loaded {
		return
	}

	// Otherwise, load all children first
	for _, child := range d.Children {
		child.tryLoad()
	}

	// Load the roll
	d.roll = GetRollType(d.RollType)
	d.roll.Load(d.RollParams)

	// Load the aggregation.
	// Skip empty aggregations as base rolls don't use them
	if d.AggregationType != "" {
		d.aggregation = GetAggregationType(d.AggregationType)
		d.aggregation.Load(d.AggregationParams)
	}

	// Set the loaded flag to skip this process in the future
	d.loaded = true

}

// Perform the roll described by the definition, and return a result object
func (d *Definition) Roll() *Result {

	// Try to load the roll
	d.tryLoad()

	result := d.roll.Roll(d.source, d.Children)
	if result.Base {
		return result
	}
	d.aggregation.Aggregate(result)
	return result

}
