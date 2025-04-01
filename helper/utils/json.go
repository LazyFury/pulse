package utils

import (
	"encoding/json"
)

func JSONToMap[T any](obj any) map[string]T {
	str, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	jsonStr := string(str)
	var res map[string]T
	_ = json.Unmarshal([]byte(jsonStr), &res)
	return res
}
