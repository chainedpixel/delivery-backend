package value_objects

import (
	"net/url"
	"strings"
)

type URL struct {
	value string
}

func NewURL(value string) *URL {
	return &URL{value: strings.TrimSpace(value)}
}

func (u *URL) IsValid() bool {
	_, err := url.ParseRequestURI(u.value)
	if err != nil {
		return false
	}

	parsed, err := url.Parse(u.value)
	if err != nil {
		return false
	}

	return parsed.Scheme != "" && parsed.Host != ""
}

func (u *URL) ToString() string {
	return u.value
}

func (u *URL) Equals(value ValidaterObject[string]) bool {
	return u.value == value.GetValue()
}

func (u *URL) GetValue() string {
	return u.value
}

// GetDomain obtiene el dominio de la URL
func (u *URL) GetDomain() string {
	parsed, err := url.Parse(u.value)
	if err != nil {
		return ""
	}
	return parsed.Host
}

// GetPath obtiene la ruta de la URL
func (u *URL) GetPath() string {
	parsed, err := url.Parse(u.value)
	if err != nil {
		return ""
	}
	return parsed.Path
}
