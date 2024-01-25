package dice

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/flywingedai/dice/results"
	"github.com/flywingedai/dice/roll"
)

/*
Inherited instance of Definition to allow easily accesible public methods
for roll creation/handling.
*/
type Definition struct {
	*roll.Definition
}

/*
Create a new roll definition based on a dice with "n" sides. Will set the
chances of each side being rolled to 1/n
*/
func New(n int64) *Definition {
	weights := map[int64]*big.Int{}
	values := []int64{}
	for i := int64(1); i <= n; i++ {
		weights[i] = big.NewInt(1)
		values = append(values, i)
	}

	d := &roll.Definition{
		RollType: roll.Roll_Base,
		Rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		Chances: &results.Chances{
			Values:  values,
			Weights: weights,
			Total:   big.NewInt(n),
		},
	}
	d.SetChances()

	return &Definition{Definition: d}
}
