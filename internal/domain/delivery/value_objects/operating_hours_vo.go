package value_objects

import (
	"encoding/json"
	"fmt"
	"time"
)

// OperatingHoursDay representa las horas de operación para un día
type OperatingHoursDay struct {
	Start string `json:"start"` // Formato HH:MM
	End   string `json:"end"`   // Formato HH:MM
}

// OperatingHours representa las horas de operación para una sucursal
type OperatingHours struct {
	Weekdays OperatingHoursDay `json:"weekdays"` // Lunes a Viernes
	Weekends OperatingHoursDay `json:"weekends"` // Sábados y Domingos
}

func NewOperatingHours(weekdaysStart, weekdaysEnd, weekendsStart, weekendsEnd string) *OperatingHours {
	return &OperatingHours{
		Weekdays: OperatingHoursDay{
			Start: weekdaysStart,
			End:   weekdaysEnd,
		},
		Weekends: OperatingHoursDay{
			Start: weekendsStart,
			End:   weekendsEnd,
		},
	}
}

// NewOperatingHoursFromJSON crea un objeto OperatingHours desde un string JSON
func NewOperatingHoursFromJSON(jsonStr string) (*OperatingHours, error) {
	var operatingHours OperatingHours
	if err := json.Unmarshal([]byte(jsonStr), &operatingHours); err != nil {
		return nil, err
	}

	return &operatingHours, nil
}

// IsValid verifica que las horas de operación sean válidas
func (oh *OperatingHours) IsValid() bool {
	return isValidTimeFormat(oh.Weekdays.Start) &&
		isValidTimeFormat(oh.Weekdays.End) &&
		isValidTimeFormat(oh.Weekends.Start) &&
		isValidTimeFormat(oh.Weekends.End) &&
		isStartBeforeEnd(oh.Weekdays.Start, oh.Weekdays.End) &&
		isStartBeforeEnd(oh.Weekends.Start, oh.Weekends.End)
}

// ToString convierte las horas de operación a un string legible
func (oh *OperatingHours) ToString() string {
	return fmt.Sprintf("Weekdays: %s - %s, Weekends: %s - %s",
		oh.Weekdays.Start, oh.Weekdays.End,
		oh.Weekends.Start, oh.Weekends.End)
}

// Equals compara si dos objetos OperatingHours son iguales
func (oh *OperatingHours) Equals(value ValidaterObject[OperatingHours]) bool {
	other := value.GetValue()
	return oh.Weekdays.Start == other.Weekdays.Start &&
		oh.Weekdays.End == other.Weekdays.End &&
		oh.Weekends.Start == other.Weekends.Start &&
		oh.Weekends.End == other.Weekends.End
}

// GetValue retorna el valor del objeto
func (oh *OperatingHours) GetValue() OperatingHours {
	return *oh
}

// ToJSON convierte el objeto a un string JSON
func (oh *OperatingHours) ToJSON() (string, error) {
	data, err := json.Marshal(oh)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsOpen verifica si está abierto en un momento específico
func (oh *OperatingHours) IsOpen(t time.Time) bool {
	dayOfWeek := t.Weekday()
	timeStr := t.Format("15:04")

	// Verificar si es día de semana (1-5) o fin de semana (0, 6)
	if dayOfWeek >= time.Monday && dayOfWeek <= time.Friday {
		return isTimeBetween(timeStr, oh.Weekdays.Start, oh.Weekdays.End)
	} else {
		return isTimeBetween(timeStr, oh.Weekends.Start, oh.Weekends.End)
	}
}

// Helper para verificar el formato de la hora (HH:MM)
func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// Helper para verificar que la hora de inicio sea anterior a la de fin
func isStartBeforeEnd(start, end string) bool {
	startTime, err1 := time.Parse("15:04", start)
	endTime, err2 := time.Parse("15:04", end)

	if err1 != nil || err2 != nil {
		return false
	}

	return startTime.Before(endTime)
}

// Helper para verificar si una hora está entre un inicio y un fin
func isTimeBetween(check, start, end string) bool {
	checkTime, err1 := time.Parse("15:04", check)
	startTime, err2 := time.Parse("15:04", start)
	endTime, err3 := time.Parse("15:04", end)

	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	// Manejo de caso especial donde el horario cruza la medianoche
	if endTime.Before(startTime) {
		return checkTime.After(startTime) || checkTime.Before(endTime)
	}

	return (checkTime.After(startTime) || checkTime.Equal(startTime)) &&
		(checkTime.Before(endTime) || checkTime.Equal(endTime))
}
