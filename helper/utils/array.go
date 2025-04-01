package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func InterfaceArrToStringArr[T string | int | float64 | interface{}](arr []T) []string {
	var res []string
	for _, v := range arr {
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
		res = append(res, str)
	}
	return res
}

func TryInterfaceToStringToArray(v interface{}) []string {
	switch val := v.(type) {
	case string:
		return strings.Split(val, ",")
	case []string:
		return val
	case []int:
		return InterfaceArrToStringArr(val)
	case []int64:
		return InterfaceArrToStringArr(val)
	case []float64:
		return InterfaceArrToStringArr(val)
	case []interface{}:
		return InterfaceArrToStringArr(val)
	}
	return []string{}
}

func InArray[T string | int | float64 | interface{}](arr []T, obj T) bool {
	for _, a := range arr {
		if fmt.Sprint(a) == fmt.Sprint(obj) {
			return true
		}
	}
	return false
}
