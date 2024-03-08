package filter

import (
	"grape/src/common/service"
	e "grape/src/filter/entities"

	"gorm.io/gorm"
)

type filterService struct {
	db *gorm.DB
	// client *redis.Client
}

func NewFilterService(client *service.CommonService) *filterService {
	return &filterService{db: client.DB}
}

func (c *filterService) TraceIP(ip string) ([]e.IpLocationEntity, error) {
	// var model []e.IpLocationEntity
	// err, _ := helper.Getcache(c.db.Where("geoname_id IN (?)", c.db.Select("geoname_id").Where("network >>= ?::inet", ip).Model(&e.IpBlockEntity{})), c.client, "INDEX", fmt.Sprintf("BLOCK:%s", ip), &model)
	// return model, err
	return []e.IpLocationEntity{}, nil
}
