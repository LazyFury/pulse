package validate

import (
	"fmt"
	"regexp"
	"time"
)

type Validate interface {
	GetField() string
	IsValid(val interface{}) (valid bool, errMsg string)
	GetMessage() string
}

type DefValidate struct {
	Field    string
	Required bool
	Msg      string
}

var _ Validate = (*DefValidate)(nil)

func (s *DefValidate) IsValid(val interface{}) (valid bool, errMsg string) {
	return true, ""
}

func (s *DefValidate) GetField() string {
	return s.Field
}

func (s *DefValidate) GetMessage() string {
	return s.Msg
}

type StringValidate struct {
	*DefValidate
	Min  int
	Max  int
	Regx *regexp.Regexp
}

var _ Validate = (*StringValidate)(nil)

func (s *StringValidate) IsValid(val interface{}) (valid bool, errMsg string) {
	str, ok := val.(string)
	if !ok {
		return false, fmt.Sprintf("%s not string", s.Field)
	}
	if s.Required && len(str) == 0 {
		return false, s.GetMessage()
	}

	if s.Min > 0 && len(str) < s.Min {
		return false, fmt.Sprintf("%s too short", s.Field)
	}
	if s.Max > 0 && len(str) > s.Max {
		return false, fmt.Sprintf("%s too long", s.Field)
	}
	if s.Regx != nil && !s.Regx.MatchString(str) {
		return false, fmt.Sprintf("%s not match", s.Field)
	}
	return true, ""

}

type NumberValidate struct {
	*DefValidate
	Min int
	Max int
}

var _ Validate = (*NumberValidate)(nil)

func (s *NumberValidate) IsValid(val interface{}) (valid bool, errMsg string) {
	switch val.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		if s.Min > 0 && val.(int) < s.Min {
			return false, s.Field + " too small"
		}
		if s.Max > 0 && val.(int) > s.Max {
			return false, s.Field + " too big"
		}
		return true, ""
	default:
		return false, s.Field + " not number"
	}
}

type EmailValidate struct {
	*DefValidate
}

var _ Validate = (*EmailValidate)(nil)

func (s *EmailValidate) IsValid(val interface{}) (valid bool, errMsg string) {
	str, ok := val.(string)
	if !ok {
		return false, s.Field + " not string"
	}
	if s.Required && len(str) == 0 {
		return false, s.GetMessage()
	}
	regx := "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

	valid = regexp.MustCompile(regx).MatchString(str)
	if !valid {
		return false, "not email"
	}
	return true, ""
}

type DateValidate struct {
	*DefValidate
}

var _ Validate = (*DateValidate)(nil)

func (s *DateValidate) IsValid(val interface{}) (valid bool, errMsg string) {
	t, ok := val.(time.Time)
	if !ok {
		return false, "not time.Time"
	}
	if s.Required && t.IsZero() {
		return false, s.GetMessage()
	}
	return true, ""
}

type Validator struct {
	Validates    []Validate
	AllowKeys    []string
	NotAllowKeys []string
}

func (v *Validator) AddValidate(validate Validate) {
	v.Validates = append(v.Validates, validate)
}

func (v *Validator) Validate(data map[string]interface{}) (valid bool, message string) {
	for _, validate := range v.Validates {
		val := data[validate.GetField()]
		valid, message = validate.IsValid(val)
		if !valid {
			return valid, message
		}
	}
	return valid, ""
}

func NewValidator() *Validator {
	return &Validator{}
}

func NewStringValidate(field string, required bool, msg string, min int, max int, regx *regexp.Regexp) *StringValidate {
	return &StringValidate{
		DefValidate: &DefValidate{
			Field:    field,
			Required: required,
			Msg:      msg,
		},
		Min:  min,
		Max:  max,
		Regx: regx,
	}
}

func NewNumberValidate(field string, required bool, msg string, min int, max int) *NumberValidate {
	return &NumberValidate{
		DefValidate: &DefValidate{
			Field:    field,
			Required: required,
			Msg:      msg,
		},
		Min: min,
		Max: max,
	}
}

func NewEmailValidate(field string, required bool, msg string) *EmailValidate {
	return &EmailValidate{
		DefValidate: &DefValidate{
			Field:    field,
			Required: required,
			Msg:      msg,
		},
	}
}
