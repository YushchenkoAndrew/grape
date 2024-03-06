package filter

import (
	"fmt"
	"grape/src/common/client"
	"grape/src/common/helper"
	e "grape/src/filter/entities"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type filterService struct {
	db     *gorm.DB
	client *redis.Client
}

func NewFilterService(client *client.Clients) *filterService {
	return &filterService{db: client.DB}
}

func (c *filterService) TraceIP(ip string) ([]e.IpLocationEntity, error) {
	var model []e.IpLocationEntity
	err, _ := helper.Getcache(c.db.Where("geoname_id IN (?)", c.db.Select("geoname_id").Where("network >>= ?::inet", ip).Model(&e.IpBlockEntity{})), c.client, "INDEX", fmt.Sprintf("BLOCK:%s", ip), &model)
	return model, err
}
