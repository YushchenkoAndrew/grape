package entities

import (
	e "grape/src/common/entities"
)

type LinkEntity struct {
	e.UuidEntity

	Name      string `gorm:"not null"`
	Link      string `gorm:"not null"`
	ProjectID int64  `gorm:"foreignKey:ProjectID;not null"`
}

func (*LinkEntity) TableName() string {
	return "links"
}
