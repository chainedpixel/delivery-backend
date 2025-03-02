package zones

type Coverage struct {
	ZoneID              string  `gorm:"column:zone_id;type:char(36);primaryKey"`
	CoverageArea        []byte  `gorm:"column:coverage_area;type:polygon;not null"`
	OperatingHours      string  `gorm:"column:operating_hours;type:json;not null"`
	MaxConcurrentOrders int     `gorm:"column:max_concurrent_orders;type:int;not null;default:10"`
	SurgeMultiplier     float64 `gorm:"column:surge_multiplier;type:decimal(3,2);default:1.00"`
	CoverageRules       string  `gorm:"column:coverage_rules;type:json"`

	// Inverse relationships
	Zone *Zone `gorm:"foreignKey:ZoneID;references:ID"`
}

func (Coverage) TableName() string {
	return "zone_coverage"
}
