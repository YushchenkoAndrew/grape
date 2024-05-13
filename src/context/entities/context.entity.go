package entities

import (
	e "grape/src/common/entities"
)

type ContextEntity struct {
	*e.UuidEntity
	*e.DroppableEntity

	Name string `gorm:"not null"`

	ContextableID   int64  `gorm:"not null"`
	ContextableType string `gorm:"not null"`

	ContextFields []ContextFieldEntity `gorm:"foreignKey:ContextID;references:ID" copier:"-"`
}

func (*ContextEntity) TableName() string {
	return "contexts"
}

func NewContextEntity() *ContextEntity {
	return &ContextEntity{
		UuidEntity:      e.NewUuidEntity(),
		DroppableEntity: e.NewDroppableEntity(),
	}
}
