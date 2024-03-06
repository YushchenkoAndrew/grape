package entities

type IpBlockEntity struct {
	Network                     string `csv:"network" gorm:"type:cidr"`
	GeonameId                   int64  `csv:"geoname_id"`
	RegisteredCountryGeonameId  int64  `csv:"registered_country_geoname_id"`
	RepresentedCountryGeonameId int64  `csv:"represented_country_geoname_id"`
	IsAnonymousProxy            bool   `csv:"is_anonymous_proxy"`
	IsSatelliteProvider         bool   `csv:"is_satellite_provider"`
}

func (*IpBlockEntity) TableName() string {
	return "ip_blocks"
}

// func (c *GeoIpBlocks) Migrate(db *gorm.DB, forced bool) error {
// 	if forced {
// 		db.Migrator().DropTable(c)
// 	}

// 	if err := db.AutoMigrate(c); err != nil {
// 		return err
// 	}

// 	var nSize int64
// 	if db.Model(c).Count(&nSize); nSize == 0 {

// 		// The most quick and easiest way !!!
// 		db.Exec(fmt.Sprintf("copy geo_ip_blocks from '%s/GeoLite2-Country-Blocks.csv' delimiter ',' csv header;", config.ENV.MigrationPath))
// 		db.Exec("CREATE INDEX geo_ip_blocks_network_idx ON geo_ip_blocks USING gist (network inet_ops);")
// 	}

// 	return nil
// }
