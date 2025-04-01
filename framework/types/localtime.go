package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

var _ JSON = (*LocalTime)(nil)
var _ VALUE = (*LocalTime)(nil)

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	if tTime.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}
func (l *LocalTime) UnmarshalJSON(b []byte) error {
	cleanB := b[1 : len(b)-1]
	tTime, err := time.Parse("2006-01-02 15:04:05", string(cleanB))
	if err != nil {
		return err
	}
	*l = LocalTime(tTime)
	return nil
}

func (l *LocalTime) Scan(v interface{}) error {
	switch v := v.(type) {
	case time.Time:
		*l = LocalTime(v)
		return nil
	default:
		return nil
	}
}

func (l LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tTime := time.Time(l)
	if tTime.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(l), nil
}

// format
func (l LocalTime) Format(format string) string {
	return time.Time(l).Format(format)
}
