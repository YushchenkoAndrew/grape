package models

import (
	"api/interfaces"
	"time"

	"gorm.io/gorm"
)

type World struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
	Country   string    `gorm:"size:2;unique" json:"country" xml:"country" example:"UK"`
	Visitors  uint16    `gorm:"default:0" json:"visitors" xml:"visitors" example:"5"`
}

func NewWorld() interfaces.Table {
	return &World{}
}

func (*World) TableName() string {
	return "world"
}

func (c *World) Migrate(db *gorm.DB, forced bool) error {
	if forced {
		db.Migrator().DropTable(c)
	}

	return db.AutoMigrate(c)
}

type WorldDto struct {
	// ID        uint32
	// UpdatedAt time.Time
	Country  string  `json:"country" xml:"country"`
	Visitors *uint16 `json:"visitors" xml:"visitors"`
}
