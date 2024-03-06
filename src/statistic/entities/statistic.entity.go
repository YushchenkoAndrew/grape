package entities

import (
	e "grape/src/common/entities"
)

type StatisticEntity struct {
	e.UuidEntity

	// Countries string    `json:"countries" xml:"contries" example:"UK,US"`
	Views    uint16 `gorm:"default:0" json:"views" xml:"views" example:"1"`
	Clicks   uint16 `gorm:"default:0" json:"clicks" xml:"clicks" example:"2"`
	Media    uint16 `gorm:"default:0" json:"media" xml:"media" example:"3"`
	Visitors uint16 `gorm:"default:0" json:"visitors" xml:"visitors" example:"4"`
}

func (*StatisticEntity) TableName() string {
	return "statistics"
}
