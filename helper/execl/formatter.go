package execl

import "github.com/lazyfury/pulse/framework/types"

func TimeFormatter(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		return v
	case types.LocalTime:
		return val.Format("2006-01-02 15:04:05")
	default:
		return v
	}
}
