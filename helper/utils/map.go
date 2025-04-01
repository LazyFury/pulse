package utils

import (
	"fmt"
	"strconv"
)

func InterfaceMapToStringMap[T string | int | float64 | interface{}](m map[string]T) map[string]string {
	res := map[string]string{}
	for k, v := range m {
		str := ""
		switch val := any(v).(type) {
		case string:
			str = val
		case int:
			str = strconv.Itoa(val)
		case float64:
			str = strconv.FormatFloat(val, 'f', -1, 64)
		case interface{}:
			str = fmt.Sprint(val)
		default:
			str = fmt.Sprint(v)
		}
		res[k] = str
	}
	return res
}

func StringMapToInterfaceMap(m map[string]string) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		res[k] = v
	}
	return res
}

func GetFromMapWithDefault[T string | int | float64 | interface{}](m map[string]T, key string, defaultVal T) T {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultVal
}
