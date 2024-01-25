package utils

import "encoding/json"

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
