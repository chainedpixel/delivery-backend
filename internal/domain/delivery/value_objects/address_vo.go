package value_objects

import "strings"

type Address struct {
	line1 string
	line2 string
	city  string
	state string
}

func NewAddress(line1, line2, city, state, postalCode string) *Address {
	return &Address{
		line1: strings.TrimSpace(line1),
		line2: strings.TrimSpace(line2),
		city:  strings.TrimSpace(city),
		state: strings.TrimSpace(state),
	}
}

func (a *Address) IsValid() bool {
	return a.line1 != "" && a.city != "" && a.state != ""
}

func (a *Address) ToString() string {
	address := a.line1
	if a.line2 != "" {
		address += ", " + a.line2
	}
	address += ", " + a.city + ", " + a.state
	return address
}

func (a *Address) Equals(value ValidaterObject[Address]) bool {
	other := value.GetValue()
	return a.line1 == other.line1 &&
		a.line2 == other.line2 &&
		a.city == other.city &&
		a.state == other.state
}

func (a *Address) GetValue() Address {
	return *a
}

func (a *Address) Line1() string {
	return a.line1
}

func (a *Address) Line2() string {
	return a.line2
}

func (a *Address) City() string {
	return a.city
}

func (a *Address) State() string {
	return a.state
}
