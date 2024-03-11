package entities

import (
	e "grape/src/common/entities"
)

type LocationEntity struct {
	*e.BasicEntity

	LocaleCode     string `gorm:"not null;size:2"`
	ContinentCode  string `gorm:"not null;size:2"`
	ContinentName  string `gorm:"not null"`
	CountryIsoCode string `gorm:"not null;size:2"`
	CountryName    string `gorm:"not null"`
	GeonameId      int64  `gorm:"not null"`

	IpBlock NetworkEntity `gorm:"-"`
}

func (*LocationEntity) TableName() string {
	return "locations"
}
