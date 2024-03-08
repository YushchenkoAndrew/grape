package entities

import "time"

type IdEntity struct {
	ID        int64     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
}

type UuidEntity struct {
	IdEntity
	UUID string `gorm:"unique;not null"`
}
