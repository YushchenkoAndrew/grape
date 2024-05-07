package entities

import (
	e "grape/src/common/entities"
)

type ContextEntity struct {
	*e.UuidEntity

	Name  string `gorm:"not null"`
	Order int    `gorm:"not null;default:1" copier:"-"`

	ContextableID   int64  `gorm:"not null"`
	ContextableType string `gorm:"not null"`

	ContextFields []ContextFieldEntity `gorm:"foreignKey:ContextID;references:ID" copier:"-"`
}

func (*ContextEntity) TableName() string {
	return "contexts"
}

func (c *ContextEntity) SetOrder(order int) {
	c.Order = order
}

func (c *ContextEntity) GetOrder() int {
	return c.Order
}

func NewContextEntity() *ContextEntity {
	return &ContextEntity{UuidEntity: e.NewUuidEntity()}
}
