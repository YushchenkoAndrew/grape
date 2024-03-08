package entities

import (
	e "grape/src/common/entities"
)

type UserEntity struct {
	*e.UuidEntity

	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
	Status   int    `gorm:"not null;default:0"`

	// Organization OrganizationEntity `gorm:"foreignKey:organization_id;not null"`
}

func (*UserEntity) TableName() string {
	return "users"
}
