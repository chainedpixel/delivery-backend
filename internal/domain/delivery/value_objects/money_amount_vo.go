package value_objects

import (
	"fmt"
	"math"
)

type MoneyAmount struct {
	value    float64
	currency string
}

func NewMoneyAmount(value float64, currency string) *MoneyAmount {
	// Redondear a 2 decimales
	rounded := math.Round(value*100) / 100
	return &MoneyAmount{
		value:    rounded,
		currency: currency,
	}
}

func (m *MoneyAmount) IsValid() bool {
	return !math.IsNaN(m.value) && !math.IsInf(m.value, 0)
}

func (m *MoneyAmount) ToString() string {
	return fmt.Sprintf("%.2f %s", m.value, m.currency)
}

func (m *MoneyAmount) Equals(value ValidaterObject[MoneyAmount]) bool {
	other := value.GetValue()
	// Comparar con una pequeña tolerancia para valores de punto flotante
	const epsilon = 0.005 // Tolerancia de medio centavo
	return math.Abs(m.value-other.value) < epsilon &&
		m.currency == other.currency
}

func (m *MoneyAmount) GetValue() MoneyAmount {
	return *m
}

func (m *MoneyAmount) Amount() float64 {
	return m.value
}

func (m *MoneyAmount) Currency() string {
	return m.currency
}

// Add suma otro MoneyAmount (de la misma moneda)
func (m *MoneyAmount) Add(other *MoneyAmount) (*MoneyAmount, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("cannot add different currencies: %s and %s", m.currency, other.currency)
	}
	return NewMoneyAmount(m.value+other.value, m.currency), nil
}

// Subtract resta otro MoneyAmount (de la misma moneda)
func (m *MoneyAmount) Subtract(other *MoneyAmount) (*MoneyAmount, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("cannot subtract different currencies: %s and %s", m.currency, other.currency)
	}
	return NewMoneyAmount(m.value-other.value, m.currency), nil
}
