package value_objects

import (
	"fmt"
	"math"
)

type Percentage struct {
	value float64 // Valor de 0 a 100
}

func NewPercentage(value float64) *Percentage {
	return &Percentage{value: value}
}

func (p *Percentage) IsValid() bool {
	return p.value >= 0 && p.value <= 100 && !math.IsNaN(p.value)
}

func (p *Percentage) ToString() string {
	return fmt.Sprintf("%.2f%%", p.value)
}

func (p *Percentage) Equals(value ValidaterObject[float64]) bool {
	// Comparar con una pequeña tolerancia para valores de punto flotante
	const epsilon = 0.01
	return math.Abs(p.value-value.GetValue()) < epsilon
}

func (p *Percentage) GetValue() float64 {
	return p.value
}

// AsDecimal devuelve el valor como decimal (0-1)
func (p *Percentage) AsDecimal() float64 {
	return p.value / 100.0
}

// FromDecimal crea un porcentaje a partir de un valor decimal
func NewPercentageFromDecimal(decimal float64) *Percentage {
	return NewPercentage(decimal * 100.0)
}
