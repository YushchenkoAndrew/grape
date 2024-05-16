package entities

import (
	e "grape/src/common/entities"
)

type TagEntity struct {
	*e.UuidEntity

	Name string `gorm:"not null"`

	TaggableID   int64  `gorm:"not null" copier:"-"`
	TaggableType string `gorm:"not null" copier:"-"`
}

func (*TagEntity) TableName() string {
	return "tags"
}

func NewTagEntity() *TagEntity {
	return &TagEntity{UuidEntity: e.NewUuidEntity()}
}
