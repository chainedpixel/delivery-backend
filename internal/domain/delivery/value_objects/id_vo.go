package value_objects

import (
	"github.com/google/uuid"
	"strings"
)

type ID struct {
	value string
}

func NewID(value string) *ID {
	return &ID{value: value}
}

func GenerateID() *ID {
	return &ID{value: uuid.NewString()}
}

func (id *ID) IsValid() bool {
	_, err := uuid.Parse(id.value)
	return err == nil
}

func (id *ID) ToString() string {
	return id.value
}

func (id *ID) Equals(value ValidaterObject[string]) bool {
	return strings.EqualFold(id.value, value.GetValue())
}

func (id *ID) GetValue() string {
	return id.value
}
