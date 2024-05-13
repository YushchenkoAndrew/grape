package entities

import (
	e "grape/src/common/entities"
)

type UserEntity struct {
	*e.UuidEntity
	*e.DeleteableEntity

	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`

	OrganizationID int64              `gorm:"not null"`
	Organization   OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID"`
}

func (*UserEntity) TableName() string {
	return "users"
}

func NewUserEntity() *UserEntity {
	return &UserEntity{
		UuidEntity:       e.NewUuidEntity(),
		DeleteableEntity: e.NewDeleteableEntity(),
	}
}
