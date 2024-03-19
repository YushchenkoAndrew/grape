package entities

import (
	e "grape/src/common/entities"
)

type LocationEntity struct {
	*e.UuidEntity

	LocaleCode     string `gorm:"not null;size:2"`
	ContinentCode  string `gorm:"not null;size:2"`
	ContinentName  string `gorm:"not null"`
	CountryIsoCode string `gorm:"not null;size:2"`
	CountryName    string `gorm:"not null"`

	Network NetworkEntity `gorm:"foreignKey:LocationID"`
}

func (*LocationEntity) TableName() string {
	return "locations"
}

func NewLocationEntity() *LocationEntity {
	return &LocationEntity{UuidEntity: e.NewUuidEntity()}
}
