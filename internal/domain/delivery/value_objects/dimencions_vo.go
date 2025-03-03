package value_objects

import (
	"encoding/json"
	"fmt"
	"math"
)

type Dimensions struct {
	length float64
	width  float64
	height float64
	unit   string // "cm" por defecto
}

func NewDimensions(length, width, height float64, unit string) *Dimensions {
	if unit == "" {
		unit = "cm"
	}
	return &Dimensions{
		length: length,
		width:  width,
		height: height,
		unit:   unit,
	}
}

// FromJSON crea Dimensions desde un string JSON
func NewDimensionsFromJSON(jsonStr string) (*Dimensions, error) {
	var data struct {
		Length float64 `json:"length"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
		Unit   string  `json:"unit"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}

	if data.Unit == "" {
		data.Unit = "cm"
	}

	return NewDimensions(data.Length, data.Width, data.Height, data.Unit), nil
}

func (d *Dimensions) IsValid() bool {
	return d.length > 0 && d.width > 0 && d.height > 0
}

func (d *Dimensions) ToString() string {
	return fmt.Sprintf("%.2f x %.2f x %.2f %s", d.length, d.width, d.height, d.unit)
}

func (d *Dimensions) Equals(value ValidaterObject[Dimensions]) bool {
	other := value.GetValue()
	// Comparar con una pequeña tolerancia para valores de punto flotante
	const epsilon = 0.01
	return math.Abs(d.length-other.length) < epsilon &&
		math.Abs(d.width-other.width) < epsilon &&
		math.Abs(d.height-other.height) < epsilon &&
		d.unit == other.unit
}

func (d *Dimensions) GetValue() Dimensions {
	return *d
}

func (d *Dimensions) Length() float64 {
	return d.length
}

func (d *Dimensions) Width() float64 {
	return d.width
}

func (d *Dimensions) Height() float64 {
	return d.height
}

func (d *Dimensions) Unit() string {
	return d.unit
}

// Volume calcula el volumen
func (d *Dimensions) Volume() float64 {
	return d.length * d.width * d.height
}

// ToJSON convierte a formato JSON
func (d *Dimensions) ToJSON() (string, error) {
	data := struct {
		Length float64 `json:"length"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
		Unit   string  `json:"unit"`
	}{
		Length: d.length,
		Width:  d.width,
		Height: d.height,
		Unit:   d.unit,
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
