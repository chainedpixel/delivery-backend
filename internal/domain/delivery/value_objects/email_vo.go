package value_objects

import (
	"regexp"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) *Email {
	return &Email{value: strings.TrimSpace(strings.ToLower(value))}
}

func (e *Email) IsValid() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(e.value)
}

func (e *Email) ToString() string {
	return e.value
}

func (e *Email) Equals(value ValidaterObject[string]) bool {
	return strings.EqualFold(e.value, value.GetValue())
}

func (e *Email) GetValue() string {
	return e.value
}
