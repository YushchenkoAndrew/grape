package entities

import (
	e "grape/src/common/entities"
)

type OrganizationEntity struct {
	*e.UuidEntity

	Name string `gorm:"not null"`
}

func (*OrganizationEntity) TableName() string {
	return "organizations"
}
