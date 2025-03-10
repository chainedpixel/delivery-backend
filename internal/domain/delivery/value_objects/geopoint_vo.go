package value_objects

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type GeoPoint struct {
	latitude  float64
	longitude float64
}

// NewGeoPoint crea un GeoPoint desde valores de latitud y longitud
func NewGeoPoint(latitude, longitude float64) *GeoPoint {
	return &GeoPoint{
		latitude:  latitude,
		longitude: longitude,
	}
}

// NewGeoPointFromString crea un GeoPoint desde un string en formato "latitud,longitud"
func NewGeoPointFromString(coordStr string) (*GeoPoint, error) {
	parts := strings.Split(strings.TrimSpace(coordStr), ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid coordinate format, expected 'latitude,longitude'")
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude: %w", err)
	}

	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude: %w", err)
	}

	return NewGeoPoint(lat, lon), nil
}

// FromWKT crea un GeoPoint desde un string WKT (Well-Known Text)
// Ejemplo: "POINT(73.123456 40.123456)"
func NewGeoPointFromWKT(wkt string) (*GeoPoint, error) {
	wkt = strings.TrimSpace(wkt)
	if !strings.HasPrefix(wkt, "POINT(") || !strings.HasSuffix(wkt, ")") {
		return nil, fmt.Errorf("invalid WKT POINT format")
	}

	// Extraer las coordenadas dentro del paréntesis
	coords := wkt[6 : len(wkt)-1]
	parts := strings.Split(coords, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid WKT POINT coordinates")
	}

	// NOTA: WKT usa el formato "longitude latitude"
	lon, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude in WKT: %w", err)
	}

	lat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude in WKT: %w", err)
	}

	return NewGeoPoint(lat, lon), nil
}

func (p *GeoPoint) IsValid() bool {
	return p.latitude >= -90 && p.latitude <= 90 &&
		p.longitude >= -180 && p.longitude <= 180
}

// ToString devuelve el GeoPoint en formato "latitud,longitud"
func (p *GeoPoint) ToString() string {
	return fmt.Sprintf("%.6f,%.6f", p.latitude, p.longitude)
}

// ToWKT devuelve el GeoPoint en formato WKT
func (p *GeoPoint) ToWKT() string {
	return fmt.Sprintf("POINT(%.6f %.6f)", p.longitude, p.latitude)
}

// ToGeoJSON devuelve el GeoPoint en formato GeoJSON
func (p *GeoPoint) ToGeoJSON() string {
	return fmt.Sprintf(`{"type":"Point","coordinates":[%.6f,%.6f]}`, p.longitude, p.latitude)
}

func (p *GeoPoint) Equals(value ValidaterObject[GeoPoint]) bool {
	other := value.GetValue()
	// Comparar con una pequeña tolerancia para valores de punto flotante
	const epsilon = 0.000001
	return math.Abs(p.latitude-other.latitude) < epsilon &&
		math.Abs(p.longitude-other.longitude) < epsilon
}

func (p *GeoPoint) GetValue() GeoPoint {
	return *p
}

func (p *GeoPoint) Latitude() float64 {
	return p.latitude
}

func (p *GeoPoint) Longitude() float64 {
	return p.longitude
}

// DistanceTo calcula la distancia en kilómetros a otro punto usando la fórmula Haversine
func (p *GeoPoint) DistanceTo(other *GeoPoint) float64 {
	const earthRadius = 6371.0

	lat1 := p.latitude * math.Pi / 180
	lon1 := p.longitude * math.Pi / 180
	lat2 := other.latitude * math.Pi / 180
	lon2 := other.longitude * math.Pi / 180

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
