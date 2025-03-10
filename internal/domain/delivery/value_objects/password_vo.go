package value_objects

import "regexp"

type Password struct {
	value string
}

func NewPassword(value string) *Password {
	return &Password{value: value}
}

func (p *Password) GetValue() string {
	return p.value
}

func (p *Password) Equals(other ValidaterObject[string]) bool {
	return p.value == other.GetValue()
}

func (p *Password) ToString() string {
	return p.value
}

func (p *Password) IsValid() bool {
	regex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
	return regexp.MustCompile(regex).MatchString(p.value)
}
