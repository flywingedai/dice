package results

import (
	"math/big"
	"sort"
)

type Chances struct {
	Values           []int64
	Weights          map[int64]*big.Int
	CumulativeWeight map[int64]*big.Int
	Total            *big.Int
}

/*
Creates a generator function that returns all the different combinations of
chances in a given aggregation. Pass in a function which determines how to
create a new T based on the existing chances and T that we are looking at.
*/
func CombineChances_Hash[T comparable](
	chancesList []*Chances,
	createHash func(oldHash *T, newValue int64) T,
	convertToValue func(hashValue T) int64,
) *Chances {

	total := big.NewInt(1)
	iterator := map[T]*big.Int{}

	for _, chances := range chancesList {
		loopIterator := map[T]*big.Int{}
		total.Mul(total, chances.Total)

		for value, valueWeight := range chances.Weights {

			if len(iterator) > 0 {
				for key, keyWeight := range iterator {
					newKey := createHash(&key, value)

					if weight, ok := loopIterator[newKey]; ok {
						weight.Add(weight, new(big.Int).Mul(valueWeight, keyWeight))
					} else {
						loopIterator[newKey] = new(big.Int).Mul(valueWeight, keyWeight)
					}
				}
			} else {
				newKey := createHash(nil, value)
				loopIterator[newKey] = valueWeight
			}

		}

		iterator = loopIterator
	}

	weights := map[int64]*big.Int{}
	values := []int64{}
	for key, weight := range iterator {
		value := convertToValue(key)
		weights[value] = weight
		values = append(values, value)
	}

	sort.SliceStable(values, func(i, j int) bool { return values[i] < values[j] })

	return &Chances{
		Values:  values,
		Weights: weights,
		Total:   total,
	}
}
