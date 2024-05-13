package entities

import (
	e "grape/src/common/entities"
)

type LinkEntity struct {
	*e.UuidEntity
	*e.DroppableEntity

	Name string `gorm:"not null"`
	Link string `gorm:"not null"`

	LinkableID   int64  `gorm:"not null" copier:"-"`
	LinkableType string `gorm:"not null" copier:"-"`
}

func (*LinkEntity) TableName() string {
	return "links"
}

func NewLinkEntity() *LinkEntity {
	return &LinkEntity{
		UuidEntity:      e.NewUuidEntity(),
		DroppableEntity: e.NewDroppableEntity(),
	}
}
