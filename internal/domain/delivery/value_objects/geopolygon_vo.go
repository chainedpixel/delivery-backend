package value_objects

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// GeoPolygon representa un polígono geográfico como una serie de puntos
type GeoPolygon struct {
	vertices []*GeoPoint // Lista de vértices del polígono
}

// NewGeoPolygon crea un nuevo polígono a partir de una lista de puntos
func NewGeoPolygon(points []*GeoPoint) *GeoPolygon {
	return &GeoPolygon{
		vertices: points,
	}
}

// NewGeoPolygonFromString crea un GeoPolygon desde un string con formato "lat1,lng1;lat2,lng2;..."
func NewGeoPolygonFromString(polyStr string) (*GeoPolygon, error) {
	pointStrings := strings.Split(polyStr, ";")
	if len(pointStrings) < 3 {
		return nil, fmt.Errorf("polygon must have at least 3 vertices")
	}

	points := make([]*GeoPoint, len(pointStrings))
	for i, pointStr := range pointStrings {
		point, err := NewGeoPointFromString(pointStr)
		if err != nil {
			return nil, fmt.Errorf("invalid point at index %d: %w", i, err)
		}
		points[i] = point
	}

	return NewGeoPolygon(points), nil
}

// NewGeoPolygonFromWKT crea un GeoPolygon desde un string WKT (Well-Known Text)
// Ejemplo: "POLYGON((73.123 40.123, 73.234 40.234, 73.345 40.345, 73.123 40.123))"
func NewGeoPolygonFromWKT(wkt string) (*GeoPolygon, error) {
	wkt = strings.TrimSpace(wkt)
	if !strings.HasPrefix(wkt, "POLYGON((") || !strings.HasSuffix(wkt, "))") {
		return nil, fmt.Errorf("invalid WKT POLYGON format")
	}

	// Extraer las coordenadas
	coords := wkt[9 : len(wkt)-2]
	pointStrings := strings.Split(coords, ",")
	if len(pointStrings) < 3 {
		return nil, fmt.Errorf("polygon must have at least 3 vertices")
	}

	points := make([]*GeoPoint, len(pointStrings))
	for i, pointStr := range pointStrings {
		pointStr = strings.TrimSpace(pointStr)
		parts := strings.Split(pointStr, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid WKT coordinate format at index %d", i)
		}

		// WKT usa el formato "longitude latitude"
		lng, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude at index %d: %w", i, err)
		}

		lat, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude at index %d: %w", i, err)
		}

		points[i] = NewGeoPoint(lat, lng)
	}

	return NewGeoPolygon(points), nil
}

// NewGeoPolygonFromGeoJSON crea un GeoPolygon desde un string GeoJSON
func NewGeoPolygonFromGeoJSON(geojson string) (*GeoPolygon, error) {
	var data struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	}

	if err := json.Unmarshal([]byte(geojson), &data); err != nil {
		return nil, fmt.Errorf("invalid GeoJSON: %w", err)
	}

	if data.Type != "Polygon" {
		return nil, fmt.Errorf("expected GeoJSON type 'Polygon', got '%s'", data.Type)
	}

	if len(data.Coordinates) == 0 || len(data.Coordinates[0]) < 3 {
		return nil, fmt.Errorf("polygon must have at least 3 vertices")
	}

	// Usamos el primer anillo de coordenadas (exterior)
	ring := data.Coordinates[0]
	points := make([]*GeoPoint, len(ring))

	for i, coord := range ring {
		if len(coord) < 2 {
			return nil, fmt.Errorf("invalid coordinate at index %d", i)
		}
		// GeoJSON usa [longitude, latitude]
		points[i] = NewGeoPoint(coord[1], coord[0])
	}

	return NewGeoPolygon(points), nil
}

func (p *GeoPolygon) IsValid() bool {
	// Un polígono necesita al menos 3 vértices
	if len(p.vertices) < 3 {
		return false
	}

	// El primer y último punto deben ser iguales (polígono cerrado)
	lastIdx := len(p.vertices) - 1
	first := p.vertices[0]
	last := p.vertices[lastIdx]

	if !first.Equals(last) && !p.vertices[0].Equals(p.vertices[lastIdx]) {
		// Si no están cerrados, los validamos individualmente
		for _, vertex := range p.vertices {
			if !vertex.IsValid() {
				return false
			}
		}
		// El polígono es válido pero no está cerrado
	}

	return true
}

// ToString devuelve el polígono en formato simple lat,lng;lat,lng;...
func (p *GeoPolygon) ToString() string {
	points := make([]string, len(p.vertices))
	for i, vertex := range p.vertices {
		points[i] = vertex.ToString()
	}
	return strings.Join(points, ";")
}

// ToWKT devuelve el polígono en formato WKT
func (p *GeoPolygon) ToWKT() string {
	points := make([]string, len(p.vertices))
	for i, vertex := range p.vertices {
		points[i] = fmt.Sprintf("%.6f %.6f", vertex.Longitude(), vertex.Latitude())
	}
	return fmt.Sprintf("POLYGON((%s))", strings.Join(points, ", "))
}

// ToGeoJSON devuelve el polígono en formato GeoJSON
func (p *GeoPolygon) ToGeoJSON() string {
	coords := make([]string, len(p.vertices))
	for i, vertex := range p.vertices {
		coords[i] = fmt.Sprintf("[%.6f,%.6f]", vertex.Longitude(), vertex.Latitude())
	}
	return fmt.Sprintf(`{"type":"Polygon","coordinates":[[%s]]}`, strings.Join(coords, ","))
}

func (p *GeoPolygon) Equals(value ValidaterObject[GeoPolygon]) bool {
	other := value.GetValue()
	if len(p.vertices) != len(other.vertices) {
		return false
	}

	for i := 0; i < len(p.vertices); i++ {
		if !p.vertices[i].Equals(other.vertices[i]) {
			return false
		}
	}
	return true
}

func (p *GeoPolygon) GetValue() GeoPolygon {
	return *p
}

func (p *GeoPolygon) Vertices() []*GeoPoint {
	return p.vertices
}

// ContainsPoint verifica si un punto está dentro del polígono usando el algoritmo "ray casting"
func (p *GeoPolygon) ContainsPoint(point *GeoPoint) bool {
	inside := false
	j := len(p.vertices) - 1

	for i := 0; i < len(p.vertices); i++ {
		if ((p.vertices[i].Latitude() > point.Latitude()) != (p.vertices[j].Latitude() > point.Latitude())) &&
			(point.Longitude() < p.vertices[i].Longitude()+(p.vertices[j].Longitude()-p.vertices[i].Longitude())*
				(point.Latitude()-p.vertices[i].Latitude())/(p.vertices[j].Latitude()-p.vertices[i].Latitude())) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// Centroid calcula el punto central (centroide) del polígono
func (p *GeoPolygon) Centroid() *GeoPoint {
	var area float64
	var centerX, centerY float64

	for i := 0; i < len(p.vertices); i++ {
		j := (i + 1) % len(p.vertices)
		f := p.vertices[i].Latitude()*p.vertices[j].Longitude() -
			p.vertices[j].Latitude()*p.vertices[i].Longitude()

		centerX += (p.vertices[i].Latitude() + p.vertices[j].Latitude()) * f
		centerY += (p.vertices[i].Longitude() + p.vertices[j].Longitude()) * f
		area += f
	}

	area /= 2
	area = math.Abs(area)

	centerX /= (6 * area)
	centerY /= (6 * area)

	return NewGeoPoint(centerX, centerY)
}

// Area calcula el área aproximada del polígono en metros cuadrados
// Usa la fórmula de Gauss (shoelace formula)
func (p *GeoPolygon) Area() float64 {
	const R = 6378137

	if len(p.vertices) < 3 {
		return 0
	}

	area := 0.0
	for i := 0; i < len(p.vertices); i++ {
		j := (i + 1) % len(p.vertices)

		// Convertir a radianes
		lat1 := p.vertices[i].Latitude() * math.Pi / 180
		lng1 := p.vertices[i].Longitude() * math.Pi / 180
		lat2 := p.vertices[j].Latitude() * math.Pi / 180
		lng2 := p.vertices[j].Longitude() * math.Pi / 180

		// Fórmula de área geodésica
		area += (lng2 - lng1) * (2 + math.Sin(lat1) + math.Sin(lat2))
	}

	area = area * R * R / 2.0
	return math.Abs(area)
}

// Perimeter calcula el perímetro aproximado del polígono en metros
func (p *GeoPolygon) Perimeter() float64 {
	perimeter := 0.0
	for i := 0; i < len(p.vertices); i++ {
		j := (i + 1) % len(p.vertices)
		perimeter += p.vertices[i].DistanceTo(p.vertices[j]) * 1000
	}
	return perimeter
}
