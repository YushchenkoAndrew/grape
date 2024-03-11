package entities

import (
	e "grape/src/common/entities"
)

type NetworkEntity struct {
	*e.BasicEntity

	Network   string `gorm:"not null;type:cidr"`
	GeonameId int64  `gorm:"not null"`
}

func (*NetworkEntity) TableName() string {
	return "networks"
}
