package value_objects

import (
	"regexp"
	"strings"
)

type TrackingNumber struct {
	value string
}

func NewTrackingNumber(value string) *TrackingNumber {
	return &TrackingNumber{value: strings.ToUpper(strings.TrimSpace(value))}
}

func (t *TrackingNumber) IsValid() bool {
	regex := regexp.MustCompile(`^[A-Z0-9]{8,30}$`)
	return regex.MatchString(t.value)
}

func (t *TrackingNumber) ToString() string {
	return t.value
}

func (t *TrackingNumber) Equals(value ValidaterObject[string]) bool {
	return t.value == value.GetValue()
}

func (t *TrackingNumber) GetValue() string {
	return t.value
}
