package types

import "database/sql/driver"

type JSON interface {
	UnmarshalJSON([]byte) error
	MarshalJSON() ([]byte, error)
}

type MarshalJSON interface {
	MarshalJSON() ([]byte, error)
}

type VALUE interface {
	Scan(interface{}) error
	Value() (driver.Value, error)
}
