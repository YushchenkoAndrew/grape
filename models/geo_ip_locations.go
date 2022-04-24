package models

import (
	"api/config"
	"api/interfaces"
	"fmt"

	"gorm.io/gorm"
)

type GeoIpLocations struct {
	GeonameId         int64  `csv:"geoname_id" json:"geoname_id" xml:"geoname_id" example:"690791"`
	LocaleCode        string `csv:"locale_code" gorm:"size:2" json:"locale_code" xml:"locale_code" example:"en"`
	ContinentCode     string `csv:"continent_code" gorm:"size:2" json:"continent_code" xml:"continent_code" example:"EU"`
	ContinentName     string `csv:"continent_name" gorm:"size:255" json:"continent_name" xml:"continent_name" example:"Europe"`
	CountryIsoCode    string `csv:"country_iso_code" gorm:"size:2" json:"country_iso_code" xml:"country_iso_code" example:"UA"`
	CountryName       string `csv:"country_name" gorm:"size:255" json:"country_name" xml:"country_name" example:"Ukraine"`
	IsInEuropeanUnion bool   `csv:"is_in_european_union" json:"is_in_european_union" xml:"is_in_european_union" example:"false"`
}

func NewGeoIpLocations() interfaces.Table {
	return &GeoIpLocations{}
}

func (*GeoIpLocations) TableName() string {
	return "geo_ip_locations"
}

func (c *GeoIpLocations) Migrate(db *gorm.DB, forced bool) error {
	if forced {
		db.Migrator().DropTable(c)
	}

	if err := db.AutoMigrate(c); err != nil {
		return err
	}

	var nSize int64
	if db.Model(c).Count(&nSize); nSize == 0 {

		// The most quick and easiest way !!!
		db.Exec(fmt.Sprintf("copy geo_ip_locations from '%s/GeoLite2-Country-Locations-en.csv' delimiter ',' csv header;", config.ENV.MigrationPath))
	}

	return nil
}
