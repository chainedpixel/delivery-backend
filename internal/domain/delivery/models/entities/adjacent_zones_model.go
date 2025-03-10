package entities

type AdjacentZone struct {
	ZoneID          string  `gorm:"column:zone_id;type:char(36);primaryKey"`
	AdjacentZoneID  string  `gorm:"column:adjacent_zone_id;type:char(36);primaryKey"`
	Distance        float64 `gorm:"column:distance;type:decimal(10,2);not null"`
	TravelTime      int     `gorm:"column:travel_time;type:int;not null"`
	CoverageOverlap float64 `gorm:"column:coverage_overlap;type:decimal(5,2)"`
	IsActive        bool    `gorm:"column:is_active;type:boolean;default:true"`

	// Relationships
	Zone         *Zone `gorm:"foreignKey:ZoneID;references:ID"`
	AdjacentZone *Zone `gorm:"foreignKey:AdjacentZoneID;references:ID"`
}

func (AdjacentZone) TableName() string {
	return "adjacent_zones"
}
