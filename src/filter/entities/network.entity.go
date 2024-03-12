package entities

import (
	e "grape/src/common/entities"
)

type NetworkEntity struct {
	*e.BasicEntity

	Network    string `gorm:"not null;type:cidr"`
	LocationID int64  `gorm:"foreignKey:LocationID;not null"`
}

func (*NetworkEntity) TableName() string {
	return "networks"
}
