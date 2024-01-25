package roll

import (
	"github.com/flywingedai/dice/results"
	"github.com/flywingedai/dice/utils"
)

var RollTypes = map[RollType]func() Roll{
	Roll_Multiple: func() Roll { return &roll_Multiple{} },
}

type RollType string

const Roll_Base = RollType("")

type Roll interface {
	Load(params map[string]interface{})
	Chances([]*results.Chances) []*results.Chances
	Results([]*Definition) []*results.Result
}

//////////////
// MULTIPLE //
//////////////

// Generic type for handling multiple rolls
type roll_Multiple struct {
	Count int `json:"count"`
}

func (r *roll_Multiple) Load(params map[string]interface{}) {
	utils.SetParams(r, params)
}

func (r *roll_Multiple) Chances(chances []*results.Chances) []*results.Chances {
	newChances := []*results.Chances{}
	for i := 0; i < r.Count; i++ {
		newChances = append(newChances, &results.Chances{
			Weights: chances[0].Weights,
			Total:   chances[0].Total,
		})
	}

	return newChances
}

func (r *roll_Multiple) Results(definitions []*Definition) []*results.Result {
	results := []*results.Result{}
	for i := 0; i < r.Count; i++ {
		results = append(results, definitions[0].RollResult())
	}
	return results
}

const Roll_Multiple = RollType("multiple")
