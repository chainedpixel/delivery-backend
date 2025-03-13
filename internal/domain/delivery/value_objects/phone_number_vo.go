package value_objects

import (
	"regexp"
)

type PhoneNumber struct {
	value string
}

func NewPhoneNumber(value string) *PhoneNumber {
	// Eliminar caracteres no numéricos
	re := regexp.MustCompile(`\D`)
	cleaned := re.ReplaceAllString(value, "")
	return &PhoneNumber{value: cleaned}
}

func (p *PhoneNumber) IsValid() bool {
	// Permitimos números internacionales de hasta 15 dígitos
	return len(p.value) >= 8 && len(p.value) <= 15
}

func (p *PhoneNumber) ToString() string {
	return p.value
}

func (p *PhoneNumber) Equals(value ValidaterObject[string]) bool {
	return p.value == value.GetValue()
}

func (p *PhoneNumber) GetValue() string {
	return p.value
}
