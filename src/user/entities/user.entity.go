package entities

import (
	e "grape/src/common/entities"
	t "grape/src/user/types"
)

type UserEntity struct {
	e.UuidEntity

	Name     string           `gorm:"not null"`
	Password string           `gorm:"not null"`
	Status   t.UserStatusEnum `gorm:"not null;default:1"`

	OrganizationID int64              `gorm:"not null"`
	Organization   OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID"`
}

func (*UserEntity) TableName() string {
	return "users"
}
