package entities

import (
	e "grape/src/common/entities"
)

type CustomerEntity struct {
	*e.UuidEntity

	NetworkID int64          `gorm:"not null" copier:"-"`
	Network   *NetworkEntity `gorm:"foreignKey:NetworkID;references:ID" copier:"-"`

	LocationID int64           `gorm:"not null" copier:"-"`
	Location   *LocationEntity `gorm:"foreignKey:LocationID;references:ID" copier:"-"`
}

func (*CustomerEntity) TableName() string {
	return "customers"
}

// func (c *World) Migrate(db *gorm.DB, forced bool) error {
// 	if forced {
// 		db.Migrator().DropTable(c)
// 	}

// 	return db.AutoMigrate(c)
// }

// type WorldDto struct {
// 	// ID        uint32
// 	// UpdatedAt time.Time
// 	Country  string  `json:"country" xml:"country"`
// 	Visitors *uint16 `json:"visitors" xml:"visitors"`
// }
