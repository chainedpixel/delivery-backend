package value_objects

import (
	"encoding/json"
	"fmt"
	"time"
)

type TimeRange struct {
	start time.Time
	end   time.Time
}

func NewTimeRange(start, end time.Time) *TimeRange {
	return &TimeRange{
		start: start,
		end:   end,
	}
}

// FromJSON crea TimeRange desde un string JSON
func NewTimeRangeFromJSON(jsonStr string) (*TimeRange, error) {
	var data struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}

	start, err := time.Parse(time.RFC3339, data.Start)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, data.End)
	if err != nil {
		return nil, err
	}

	return NewTimeRange(start, end), nil
}

func (tr *TimeRange) IsValid() bool {
	return !tr.start.IsZero() && !tr.end.IsZero() && tr.start.Before(tr.end)
}

func (tr *TimeRange) ToString() string {
	return fmt.Sprintf("%s - %s", tr.start.Format(time.RFC3339), tr.end.Format(time.RFC3339))
}

func (tr *TimeRange) Equals(value ValidaterObject[TimeRange]) bool {
	other := value.GetValue()
	return tr.start.Equal(other.start) && tr.end.Equal(other.end)
}

func (tr *TimeRange) GetValue() TimeRange {
	return *tr
}

func (tr *TimeRange) Start() time.Time {
	return tr.start
}

func (tr *TimeRange) End() time.Time {
	return tr.end
}

// Duration calcula la duración
func (tr *TimeRange) Duration() time.Duration {
	return tr.end.Sub(tr.start)
}

// Contains verifica si un momento específico está dentro del rango
func (tr *TimeRange) Contains(t time.Time) bool {
	return (t.Equal(tr.start) || t.After(tr.start)) && (t.Equal(tr.end) || t.Before(tr.end))
}

// Overlaps verifica si este rango se solapa con otro
func (tr *TimeRange) Overlaps(other *TimeRange) bool {
	return tr.Contains(other.start) || tr.Contains(other.end) ||
		other.Contains(tr.start) || other.Contains(tr.end)
}
