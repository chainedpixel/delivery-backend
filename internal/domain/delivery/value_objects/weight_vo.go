package value_objects

import (
	"fmt"
	"math"
)

type Weight struct {
	value float64 // en kg
	unit  string  // "kg" por defecto
}

func NewWeight(value float64, unit string) *Weight {
	if unit == "" {
		unit = "kg"
	}
	return &Weight{
		value: value,
		unit:  unit,
	}
}

func (w *Weight) IsValid() bool {
	return w.value >= 0 && !math.IsNaN(w.value) && !math.IsInf(w.value, 0)
}

func (w *Weight) ToString() string {
	return fmt.Sprintf("%.2f %s", w.value, w.unit)
}

func (w *Weight) Equals(value ValidaterObject[Weight]) bool {
	other := value.GetValue()
	const epsilon = 0.01
	return math.Abs(w.value-other.value) < epsilon &&
		w.unit == other.unit
}

func (w *Weight) GetValue() Weight {
	return *w
}

func (w *Weight) Value() float64 {
	return w.value
}

func (w *Weight) Unit() string {
	return w.unit
}
